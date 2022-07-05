package main

import (
	"log"

	tw "go.temporal.io/sdk/worker"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/temporal/activities"
	"go.breu.io/ctrlplane/internal/temporal/workflows"
)

func init() {
	conf.ReadSvcConfig("worker::webhooks")
	conf.ReadDBConfig()
	conf.InitDBSession()
	conf.ReadGithubConfig()
	conf.ReadTemporalConfig()
	conf.InitTemporalClient()
}

func main() {
	defer conf.Temporal.Client.Close()

	queue := conf.Temporal.Queues.Webhooks
	options := tw.Options{}
	worker := tw.New(conf.Temporal.Client, queue, options)

	worker.RegisterWorkflow(workflows.OnGithubInstall)
	worker.RegisterActivity(activities.SaveGithubInstallation)

	err := worker.Run(tw.InterruptCh())

	if err != nil {
		log.Fatal(err)
	}
}
