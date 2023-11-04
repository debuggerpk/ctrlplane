package cloudrun

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"cloud.google.com/go/iam/apiv1/iampb"
	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/gocql/gocql"
	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/shared"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

var registerOnce sync.Once

type (
	Constructor struct{}
	workflows   struct{}

	Resource struct {
		ID                         gocql.UUID
		Cpu                        string
		Memory                     string
		Generation                 uint8
		Port                       int32
		Envs                       []*runpb.EnvVar
		OutputEnvs                 map[string]string
		Region                     string // from blueprint
		Image                      string // from workload
		Config                     string
		Name                       string
		Revision                   string
		LastRevision               string
		MinInstances               int32
		MaxInstances               int32
		AllowUnauthenticatedAccess bool
		CpuIdle                    bool
		Project                    string
		ServiceName                string
	}

	Workload struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	GCPConfig struct {
		Project string
	}
)

var (
	activities *Activities
)

// Create creates cloud run resource
func (c *Constructor) Create(name string, region string, config string, providerConfig string) (core.CloudResource, error) {
	cr := &Resource{Name: name, Region: region, Config: config}
	cr.AllowUnauthenticatedAccess = true
	cr.Cpu = "2000m"
	cr.Memory = "1024Mi"
	cr.MinInstances = 0
	cr.MaxInstances = 5
	cr.Generation = 2

	cr.Port = 8000
	cr.CpuIdle = true

	// TODO: Get env values from config
	cr.Envs = append(cr.Envs, &runpb.EnvVar{Name: "CARGOFLO_DEBUG", Values: &runpb.EnvVar_Value{Value: "false"}})
	cr.Envs = append(cr.Envs, &runpb.EnvVar{Name: "CARGOFLO_TEMPORAL_HOST", Values: &runpb.EnvVar_Value{Value: "10.10.0.3"}})
	cr.Envs = append(cr.Envs, &runpb.EnvVar{Name: "CARGOFLO_DB_HOST", Values: &runpb.EnvVar_Value{Value: "10.69.49.8"}})
	cr.Envs = append(cr.Envs, &runpb.EnvVar{Name: "CARGOFLO_DB_NAME", Values: &runpb.EnvVar_Value{Value: "cargoflo"}})
	cr.Envs = append(cr.Envs, &runpb.EnvVar{Name: "CARGOFLO_DB_USER", Values: &runpb.EnvVar_Value{Value: "cargoflo"}})
	cr.Envs = append(cr.Envs, &runpb.EnvVar{Name: "CARGOFLO_DB_PASS", Values: &runpb.EnvVar_Value{Value: "cargoflo"}})

	// get gcp project from configuration
	pconfig := new(GCPConfig)
	err := json.Unmarshal([]byte(providerConfig), pconfig)
	if err != nil {
		shared.Logger().Error("Unable to parse provider config for cloudrun")
		return nil, err
	}

	cr.Project = pconfig.Project

	shared.Logger().Info("cloud run", "object", providerConfig, "umarshaled", pconfig, "project", cr.Project)
	w := &workflows{}
	registerOnce.Do(func() {
		coreWrkr := shared.Temporal().Worker(shared.CoreQueue)
		coreWrkr.RegisterWorkflow(w.DeployWorkflow)
		coreWrkr.RegisterWorkflow(w.UpdateTrafficWorkflow)
	})
	return cr, nil
}

// CreateFromJson creates a Resource object from JSON
func (c *Constructor) CreateFromJson(data []byte) core.CloudResource {
	cr := &Resource{}
	json.Unmarshal(data, cr)
	return cr
}

// Marshal marshals the Resource object
func (r *Resource) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Provision provisions the cloud resource
func (r *Resource) Provision(ctx workflow.Context) (workflow.Future, error) {

	// do nothing, the infra will be provisioned with deployment
	return nil, nil
}

// DeProvision deprovisions the cloudrun resource
func (r *Resource) DeProvision() error {

	return nil
}

// UpdateTraffic updates the traffic distribution on latest and previous revision as per the input
// parameter trafficpcnt is the percentage traffic to be deployed on latest revision
func (r *Resource) UpdateTraffic(ctx workflow.Context, trafficpcnt int32) error {

	// UpdateTraffic will execute a workflow to update the resource. This workflow is not directly called
	// from provisioninfra workflow to avoid passing resource interface as argument
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		ChildWorkflowOptions(
			shared.WithWorkflowParent(ctx),
			shared.WithWorkflowBlock("Resource"),
			shared.WithWorkflowBlockID(r.Name),
			shared.WithWorkflowElement("UpdateTraffic"),
		)

	cctx := workflow.WithChildOptions(ctx, opts)

	shared.Logger().Info("Executing Update traffic workflow")

	w := &workflows{}
	err := workflow.
		ExecuteChildWorkflow(cctx, w.UpdateTrafficWorkflow, r, trafficpcnt).Get(cctx, nil)

	if err != nil {
		shared.Logger().Error("Could not execute UpdateTraffic workflow", "error", err)
		return err
	}
	return nil
}

func (r *Resource) Deploy(ctx workflow.Context, wl []core.Workload, changesetID gocql.UUID) error {
	shared.Logger().Info("deploying", "cloudrun", r, "workload", wl)

	if len(wl) != 1 {
		shared.Logger().Error("Cannot deploy more than one workloads on cloud run", "number of workloads", len(wl))
		return errors.New("multiple workloads defined for cloud run")
	}

	// provision with execute a workflow to provision the resources. This workflow is not directly called
	// from provisioninfra workflow to avoid passing resource interface as argument

	crworkload := &Workload{}
	json.Unmarshal([]byte(wl[0].Container), crworkload)
	crworkload.Image = crworkload.Image + ":" + changesetID.String()
	crworkload.Name = wl[0].Name

	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		ChildWorkflowOptions(
			shared.WithWorkflowParent(ctx),
			shared.WithWorkflowBlock("Resource"),
			shared.WithWorkflowBlockID(r.Name),
			shared.WithWorkflowElement("Deploy"),
		)

	cctx := workflow.WithChildOptions(ctx, opts)

	shared.Logger().Info("starting DeployCloudRun workflow")

	w := &workflows{}
	err := workflow.
		ExecuteChildWorkflow(cctx, w.DeployWorkflow, r, crworkload).Get(cctx, r)

	if err != nil {
		shared.Logger().Error("Could not start DeployCloudRun workflow", "error", err)
		return err
	}
	return nil
}

func (w *workflows) DeployWorkflow(ctx workflow.Context, r *Resource, wl *Workload) (*Resource, error) {

	r.ServiceName = wl.Name
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)
	err := workflow.ExecuteActivity(actx, activities.GetNextRevision, r).Get(actx, r)
	if err != nil {
		shared.Logger().Error("Error in Executing activity: GetNextRevision", "error", err)
		return r, err
	}

	err = workflow.ExecuteActivity(actx, activities.DeployRevision, r, wl).Get(actx, nil)
	if err != nil {
		shared.Logger().Error("Error in Executing activity: DeployDummy", "error", err)
		return r, err
	}
	return r, nil
}

// UpdateTraffic workflow executes UpdateTrafficActivity
func (w *workflows) UpdateTrafficWorkflow(ctx workflow.Context, r *Resource, trafficpcnt int32) error {
	shared.Logger().Info("Distributing traffic between revisions", r.Revision, r.LastRevision)
	activityOpts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	actx := workflow.WithActivityOptions(ctx, activityOpts)
	err := workflow.ExecuteActivity(actx, activities.UpdateTrafficActivity, r, trafficpcnt).Get(ctx, r)
	if err != nil {
		shared.Logger().Error("Error in Executing activity: UpdateTrafficActivity", "error", err)
		return err
	}
	return nil
}

func (r *Resource) GetServiceClient() (*run.ServicesClient, error) {
	client, err := run.NewServicesRESTClient(context.Background())
	if err != nil {
		shared.Logger().Error("New service rest client", "error", err)
		return nil, err
	}

	return client, err
}

// GetService gets a cloud run service from GCP
func (r *Resource) GetService(ctx context.Context) *runpb.Service {

	logger := activity.GetLogger(ctx)
	serviceClient, err := run.NewServicesRESTClient(ctx)
	if err != nil {
		shared.Logger().Error("New service rest client", "Error", err)
		return nil
	}
	defer serviceClient.Close()

	svcpath := r.GetParent() + "/services/" + r.ServiceName
	req := &runpb.GetServiceRequest{Name: svcpath}

	svc, err := serviceClient.GetService(ctx, req)

	if err != nil {
		logger.Error("Get Service", "Error, returning nil", err)
		return nil
	}

	logger.Debug("Get service", "service", svc, "error", err)
	return svc
}

// AllowAccessToAll Sets IAM policy to allow access to all users
func (r *Resource) AllowAccessToAll(ctx context.Context) error {
	logger := activity.GetLogger(ctx)
	client, err := run.NewServicesRESTClient(context.Background())
	if err != nil {
		logger.Error("New service rest client", "Error", err)
		return nil
	}

	defer func() { _ = client.Close() }()

	rsc := r.GetParent() + "/services/" + r.ServiceName

	binding := new(iampb.Binding)
	binding.Members = []string{"allUsers"}
	binding.Role = "roles/run.invoker"

	_, err = client.SetIamPolicy(
		context.Background(),
		&iampb.SetIamPolicyRequest{Resource: rsc, Policy: &iampb.Policy{Bindings: []*iampb.Binding{binding}}},
	)
	if err != nil {
		logger.Error("Set policy", "Error", err)
		return err
	}

	return nil
}

// GetServiceTemplate creates and returns the revision template for cloud run from the workload to be deployed
// revision template specifies the resource requirements, image to be deployed and traffic distribution etc.
// this template will be used for first deployment only, from next deployments the already deployed template will be
// fetched from cloudrun and the same will be used for next revision
// TODO: the above design will not work if resource definition is changed
func (r *Resource) GetServiceTemplate(ctx context.Context, wl *Workload) *runpb.Service {
	activity.GetLogger(ctx).Info("setting service template for", "revision", r.Revision)
	resources := &runpb.ResourceRequirements{Limits: map[string]string{"cpu": r.Cpu, "memory": r.Memory}, CpuIdle: r.CpuIdle}

	// unmarshaling the container here assuming that container definition will be specific to a resource
	// this can be done at a common location if the container definition turns out to be same for all resources

	containerPort := &runpb.ContainerPort{ContainerPort: r.Port}
	container := &runpb.Container{Name: wl.Name, Image: wl.Image, Resources: resources, Ports: []*runpb.ContainerPort{containerPort}, Env: r.Envs}

	scaling := &runpb.RevisionScaling{MinInstanceCount: r.MinInstances, MaxInstanceCount: r.MaxInstances}

	rt := &runpb.RevisionTemplate{Containers: []*runpb.Container{container}, Scaling: scaling,
		ExecutionEnvironment: runpb.ExecutionEnvironment(r.Generation), Revision: r.Revision}

	service := &runpb.Service{Template: rt}

	tt := &runpb.TrafficTarget{Type: runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST, Percent: 100}
	service.Traffic = []*runpb.TrafficTarget{tt}

	return service
}

func (r *Resource) GetParent() string {
	return "projects/" + r.Project + "/locations/" + r.Region
}

func (r *Resource) GetFirstRevision() string {
	return r.ServiceName + "-0"
}
