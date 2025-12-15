package main

import (
	"Hellowrld/workflow"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Connection error:", err)
	}
	defer c.Close()

	w := worker.New(c, "notification-queue", worker.Options{})

	// Workflow
	w.RegisterWorkflow(workflow.NotificationWorkflow)

	// Activity
	w.RegisterActivity(workflow.SendEmailActivity)
	w.RegisterActivity(workflow.WaitForConfirmationActivity)
	w.RegisterActivity(workflow.EscalateToManagerActivity)

	log.Println("Starting notification worker...")
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Worker error:", err)
	}
}
