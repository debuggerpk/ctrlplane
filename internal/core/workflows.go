// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/core/mutex"
	"go.breu.io/ctrlplane/internal/shared"
)

const (
	unLockTimeOutStackMutex             time.Duration = time.Minute * 30 // TODO: adjust this
	OnPullRequestWorkflowPRSignalsLimit               = 1000             // TODO: adjust this
)

type (
	Workflows        struct{}
	GetAssetsPayload struct {
		StackID     string
		RepoID      gocql.UUID
		ChangeSetID gocql.UUID
	}
)

// ChangesetController controls the rollout lifecycle for one changeset.
func (w *Workflows) ChangesetController(id string) error {
	return nil
}

// StackController runs indefinitely and controls and synchronizes all actions on stack.
// This workflow will start when createStack call is received. it will be the master workflow for all child stack workflows
// for tasks like creating infrastructure, doing deployment, apperture controller etc.
//
// The workflow waits for the signals from the git provider. It consumes events for PR created, updated, merged etc.
func (w *Workflows) StackController(ctx workflow.Context, stackID string) error {
	// deployment map is designed to be used in OnPullRequestWorkflow only
	logger := workflow.GetLogger(ctx)
	lockID := "stack." + stackID // stack.<stack id>
	deployments := make(Deployments)
	activeInfra := make(Infra)

	// create and initialize mutex, initializing mutex will start a mutex workflow
	logger.Info("creating mutex for stack", "stack", stackID)
	lock := mutex.New(
		mutex.WithCallerContext(ctx),
		mutex.WithID(lockID),
	)

	if err := lock.Start(ctx); err != nil {
		logger.Debug("unable to start mutex workflow", "error", err)
	}

	triggerChannel := workflow.GetSignalChannel(ctx, shared.WorkflowSignalDeploymentStarted.String())
	assetsChannel := workflow.GetSignalChannel(ctx, WorkflowSignalAssetsRetrieved.String())
	infrachannel := workflow.GetSignalChannel(ctx, WorkflowSignalInfraProvisioned.String())
	deploymentchannel := workflow.GetSignalChannel(ctx, WorkflowSignalDeploymentCompleted.String())
	manualOverrideChannel := workflow.GetSignalChannel(ctx, WorkflowSignalManaulOverride.String())

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(triggerChannel, onDeploymentStartedSignal(ctx, stackID, deployments))
	selector.AddReceive(assetsChannel, onAssetsRetreivedSignal(ctx, stackID, deployments))
	selector.AddReceive(infrachannel, onInfraProvisionedSignal(ctx, stackID, lock, deployments, activeInfra))
	selector.AddReceive(deploymentchannel, onDeploymentCompletedSignal(ctx, stackID, deployments))
	selector.AddReceive(manualOverrideChannel, onManualOverrideSignal(ctx, stackID, deployments))

	// var prSignalsCounter int = 0
	// return continue as new if this workflow has processed signals upto a limit
	// if prSignalsCounter >= OnPullRequestWorkflowPRSignalsLimit {
	// 	return workflow.NewContinueAsNewError(ctx, w.OnPullRequestWorkflow, stackID)
	// }
	for {
		logger.Info("waiting for signals ....")
		selector.Select(ctx)
	}
}

// DeProvisionInfra de-provisions the infrastructure created for stack deployment.
func (w *Workflows) DeProvisionInfra(ctx workflow.Context, stackID string, resourceData *ResourceConfig) error {
	return nil
}

// GetAssets gets assests for stack including resources, workloads and blueprint.
func (w *Workflows) GetAssets(ctx workflow.Context, payload *GetAssetsPayload) error {
	var (
		future workflow.Future
		err    error = nil
	)

	shared.Logger().Info("Get assets workflow")

	logger := workflow.GetLogger(ctx)
	assets := NewAssets()
	workloads := SlicedResult[Workload]{}
	resources := SlicedResult[Resource]{}
	repos := SlicedResult[Repo]{}
	blueprint := Blueprint{}

	selector := workflow.NewSelector(ctx)
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)
	providerActivityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           shared.Temporal().Queue(shared.ProvidersQueue).Name(),
	}
	pctx := workflow.WithActivityOptions(ctx, providerActivityOpts)

	// get resources for stack
	future = workflow.ExecuteActivity(actx, activities.GetResources, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &resources); err != nil {
			logger.Error("GetResources activity failed", "error", err)
			return
		}
	})

	// get workloads for stack
	future = workflow.ExecuteActivity(actx, activities.GetWorkloads, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &workloads); err != nil {
			logger.Error("GetWorkloads activity failed", "error", err)
			return
		}
	})

	// get repos for stack
	future = workflow.ExecuteActivity(actx, activities.GetRepos, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &repos); err != nil {
			logger.Error("GetRepos activity failed", "error", err)
			return
		}
	})

	// get blueprint for stack
	future = workflow.ExecuteActivity(actx, activities.GetBluePrint, payload.StackID)
	selector.AddFuture(future, func(f workflow.Future) {
		if err = f.Get(ctx, &blueprint); err != nil {
			logger.Error("GetBluePrint activity failed", "error", err)
			return
		}
	})

	// TODO: come up with a better logic for this
	for i := 0; i < 4; i++ {
		selector.Select(ctx)
		// return if activity failed. TODO: handle race conditions as the 'err' variable is shared among all activities
		if err != nil {
			logger.Error("Exiting due to activity failure")
			return err
		}
	}

	// get commits against the repos
	repoMarker := make([]ChangeSetRepoMarker, len(repos.Data))

	for idx, repo := range repos.Data {
		marker := &repoMarker[idx]
		// p := Instance().Provider(repo.Provider) // get the specific provider
		p := Instance().RepoProvider(repo.Provider) // get the specific provider
		commitID := ""

		if err := workflow.
			ExecuteActivity(pctx, p.GetLatestCommit, repo.ProviderID, repo.DefaultBranch).
			Get(ctx, &commitID); err != nil {
			logger.Error("Error in getting latest commit ID", "repo", repo.Name, "provider", repo.Provider)
			return fmt.Errorf("Error in getting latest commit ID repo:%s, provider:%s", repo.Name, repo.Provider.String())
		}

		marker.CommitID = commitID
		marker.HasChanged = repo.ID == payload.RepoID
		marker.Provider = repo.Provider.String()
		marker.RepoID = repo.ID.String()
		logger.Debug("Repo", "Name", repo.Name, "Repo marker", marker)
	}

	// save changeset
	stackID, _ := gocql.ParseUUID(payload.StackID)
	changeset := &ChangeSet{
		RepoMarkers: repoMarker,
		ID:          payload.ChangeSetID,
		StackID:     stackID,
	}

	err = workflow.ExecuteActivity(actx, activities.CreateChangeset, changeset, payload.ChangeSetID).Get(ctx, nil)
	if err != nil {
		logger.Error("Error in creating changeset")
	}

	assets.ChangesetID = payload.ChangeSetID
	assets.Blueprint = blueprint
	assets.Repos = append(assets.Repos, repos.Data...)
	assets.Resources = append(assets.Resources, resources.Data...)
	assets.Workloads = append(assets.Workloads, workloads.Data...)
	logger.Debug("Assets retreived", "Assets", assets)

	// signal parent workflow
	parent := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	_ = workflow.
		SignalExternalWorkflow(ctx, parent, "", WorkflowSignalAssetsRetrieved.String(), assets).
		Get(ctx, nil)

	return nil
}

// ProvisionInfra provisions the infrastructure required for stack deployment.
func (w *Workflows) ProvisionInfra(ctx workflow.Context, assets *Assets) error {

	logger := workflow.GetLogger(ctx)
	var children []workflow.Future

	shared.Logger().Debug("provision infra", "assets", assets)
	for _, rsc := range assets.Resources {
		logger.Info("Creating resource", "Name", rsc.Name)

		// get the resource contructor specific to the driver e.g gke, cloudrun for GCP, sns, fargate for AWS
		resconstr := Instance().CloudResources(rsc.Provider, rsc.Driver)
		if rsc.IsImmutable {
			// assuming a single region for now
			region := getRegion(rsc.Provider, &assets.Blueprint)
			providerConfig := getProviderConfig(rsc.Provider, &assets.Blueprint)
			r, err := resconstr.Create(rsc.Name, region, rsc.Config, providerConfig)

			if err != nil {
				logger.Error("could not create resource object", "ID", rsc.ID, "name", rsc.Name, "Error", err)
				return err
			}

			// resource is an interface and cannot be sent as a parameter in workflow because workflow cannot unmarshal an interface. So we need to
			// send the marshalled value in workflow and then unmarshal and resconstruct the resource again in the workflow
			ser, err := r.Marshal()
			if err != nil {
				logger.Error("could not marshal resource", "ID", rsc.ID, "name", rsc.Name, "Error", err)
				return err
			}
			assets.Infra[rsc.ID] = ser

			// TODO: initiate all resource provisions in parallel in child workflows and wait for all child workflows
			// completion before sending infra provisioned signal
			if f, err := r.Provision(ctx); err != nil {
				logger.Error("could not start resource provisioning", "ID", rsc.ID, "name", rsc.Name, "Error", err)
				return err

			} else if f != nil {
				children = append(children, f)
			}
		}
	}

	for _, child := range children {
		child.Get(ctx, nil)
	}

	shared.Logger().Info("Signaling infra provisioned")
	prWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	_ = workflow.SignalExternalWorkflow(ctx, prWorkflowID, "", WorkflowSignalInfraProvisioned.String(), assets).Get(ctx, nil)

	shared.Logger().Info("INFRA PROVISIONED")
	return nil
}

// Deploy deploys the stack.
func (w *Workflows) Deploy(ctx workflow.Context, stackID string, lock *mutex.Lock, assets *Assets) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Deployment initiated", "changeset", assets.ChangesetID, "infra", assets.Infra)
	Infra := make(Infra)

	// Acquire lock
	err := lock.Acquire(ctx)
	if err != nil {
		logger.Error("Error in acquiring lock", "Error", err)
		return err
	}

	// create deployable, map of one or more workloads against each resource
	deployables := make(map[gocql.UUID][]Workload) // map of resource id and workloads
	for _, w := range assets.Workloads {
		_, ok := deployables[w.ResourceID]
		if ok == false {
			deployables[w.ResourceID] = make([]Workload, 0)
		}
		deployables[w.ResourceID] = append(deployables[w.ResourceID], w)
	}

	// create the resource object again from marshaled data
	for _, rsc := range assets.Resources {
		resconstr := Instance().CloudResources(rsc.Provider, rsc.Driver)
		inf := assets.Infra[rsc.ID] // get marshaled resource from ID
		r := resconstr.CreateFromJson(inf)
		Infra[rsc.ID] = r
		r.Deploy(ctx, deployables[rsc.ID])
	}

	// update traffic on resource from 50 to 100
	var i int32
	for i = 50; i <= 100; i += 25 {
		for id, r := range Infra {
			shared.Logger().Info("updating traffic", id, r)
			r.UpdateTraffic(ctx, i)
			// workflow.Sleep(ctx, 10*time.Second)
		}
	}

	_ = lock.Release(ctx)

	prWorkflowID := workflow.GetInfo(ctx).ParentWorkflowExecution.ID
	workflow.SignalExternalWorkflow(ctx, prWorkflowID, "", WorkflowSignalDeploymentCompleted.String(), assets)

	return nil
}

func onManualOverrideSignal(ctx workflow.Context, stackID string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	triggerID := int64(0)

	return func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, &triggerID)
		logger.Info("manual override for", "Trigger ID", triggerID)
	}
}

// onDeploymentStartedSignal is the channel handler for trigger channel
// It will execute GetAssets and update PR deployment state to "GettingAssets".
func onDeploymentStartedSignal(ctx workflow.Context, stackID string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		// Receive signal data
		payload := &shared.PullRequestSignal{}
		channel.Receive(ctx, payload)
		logger.Info("received deployment request", "Trigger ID", payload.TriggerID)

		// We want to filter workflows with changeset ID, so create changeset ID here and use it for creating workflow ID
		changesetID, _ := gocql.RandomUUID()

		// Set childworkflow options
		opts := shared.Temporal().
			Queue(shared.CoreQueue).
			ChildWorkflowOptions(
				shared.WithWorkflowParent(ctx),
				shared.WithWorkflowElement("get_assets"),
				shared.WithWorkflowMod("trigger"),
				shared.WithWorkflowModID(strconv.FormatInt(payload.TriggerID, 10)),
			)

		getAssetsPayload := &GetAssetsPayload{
			StackID:     stackID,
			RepoID:      payload.RepoID,
			ChangeSetID: changesetID,
		}

		// execute GetAssets and wait until spawned
		var execution workflow.Execution

		cctx := workflow.WithChildOptions(ctx, opts)
		err := workflow.
			ExecuteChildWorkflow(cctx, w.GetAssets, getAssetsPayload).
			GetChildWorkflowExecution().
			Get(cctx, &execution)

		if err != nil {
			logger.Error("TODO: Error in executing GetAssets", "Error", err)
		}

		// create and save deployment data against a changeset
		deployment := NewDeployment()
		deployments[changesetID] = deployment
		deployment.state = GettingAssets
		deployment.workflows.GetAssets = execution.ID
	}
}

// onAssetsRetreivedSignal will receive assets sent by GetAssets, update deployment state and execute ProvisionInfra.
func onAssetsRetreivedSignal(ctx workflow.Context, stackID string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := NewAssets()
		channel.Receive(ctx, assets)
		logger.Info("received Assets", "changeset", assets.ChangesetID)

		// update deployment state
		deployment := deployments[assets.ChangesetID]
		deployment.state = GotAssets

		// execute provision infra workflow
		logger.Info("Executing provision Infra workflow")

		var execution workflow.Execution

		opts := shared.Temporal().
			Queue(shared.CoreQueue).
			ChildWorkflowOptions(
				shared.WithWorkflowParent(ctx),
				shared.WithWorkflowBlock("changeset"), // TODO: shouldn't this be part of the changeset controller?
				shared.WithWorkflowBlockID(assets.ChangesetID.String()),
				shared.WithWorkflowElement("provision_infra"),
			)

		cctx := workflow.WithChildOptions(ctx, opts)

		err := workflow.
			ExecuteChildWorkflow(cctx, w.ProvisionInfra, assets).
			GetChildWorkflowExecution().Get(cctx, &execution)

		if err != nil {
			logger.Error("TODO: Error in executing ProvisionInfra", "Error", err)
		}

		logger.Info("Executed provision Infra workflow")

		deployment.state = ProvisioningInfra
		deployment.workflows.ProvisionInfra = execution.ID
	}
}

// onInfraProvisionedSignal will receive assets by ProvisionInfra, update deployment state and execute Deploy.
func onInfraProvisionedSignal(ctx workflow.Context, stackID string, lock mutex.Mutex, deployments Deployments, activeinfra Infra) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)
	w := &Workflows{}

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := NewAssets()
		channel.Receive(ctx, assets)
		logger.Info("Infra provisioned", "changeset", assets.ChangesetID)

		deployment := deployments[assets.ChangesetID]
		deployment.state = InfraProvisioned

		// deployment.OldInfra = activeinfra // All traffic is currently being routed to this infra
		deployment.NewInfra = assets.Infra // handling zero traffic, no workload is deployed

		var execution workflow.Execution
		opts := shared.Temporal().
			Queue(shared.CoreQueue).
			ChildWorkflowOptions(
				shared.WithWorkflowParent(ctx),
				shared.WithWorkflowBlock("changeset"), // TODO: shouldn't this be part of the changeset controller?
				shared.WithWorkflowBlockID(assets.ChangesetID.String()),
				shared.WithWorkflowElement("deploy"),
			)
		cctx := workflow.WithChildOptions(ctx, opts)

		err := workflow.
			ExecuteChildWorkflow(cctx, w.Deploy, stackID, lock.(*mutex.Lock), assets).
			GetChildWorkflowExecution().Get(cctx, &execution)
		if err != nil {
			logger.Error("Error in Executing deployment workflow", "Error", err)
		}

		deployment.state = CreatingDeployment
		deployment.workflows.ProvisionInfra = execution.ID
	}
}

// onDeploymentCompletedSignal will conclude the deployment.
func onDeploymentCompletedSignal(ctx workflow.Context, stackID string, deployments Deployments) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		assets := NewAssets()
		channel.Receive(ctx, assets)
		logger.Info("Deployment complete", "changeset", assets.ChangesetID)
		delete(deployments, assets.ChangesetID)

		logger.Info("Deleted deployment data", "changeset", assets.ChangesetID)
	}
}
