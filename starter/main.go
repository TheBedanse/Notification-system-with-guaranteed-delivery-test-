package main

import (
	"Hellowrld/workflow"
	"context"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Connection error:", err)
	}
	defer c.Close()

	// Test data
	request := workflow.NotificationRequest{
		Message:     "Order №123 is ready",
		Recipients:  []string{"user1@example.com", "user2@example.com", "user3@example.com"},
		RequesterID: "order_system_123",
	}

	options := client.StartWorkflowOptions{
		ID:                 fmt.Sprintf("notification-%d", time.Now().Unix()),
		TaskQueue:          "notification-queue",
		WorkflowRunTimeout: time.Hour * 25, // Time reserve
	}

	// Run workflow
	we, err := c.ExecuteWorkflow(context.Background(), options,
		workflow.NotificationWorkflow, request)
	if err != nil {
		log.Fatalln("Failed to start workflow:", err)
	}

	// Receipt result
	var reports []workflow.DeliveryReport
	err = we.Get(context.Background(), &reports)
	if err != nil {
		log.Fatalln("Failed to get result:", err)
	}

	fmt.Println("\n=== Delivery Report ===")
	for _, report := range reports {
		fmt.Printf("%s: %s", report.Recipient, report.Status)
		if report.Error != "" {
			fmt.Printf(" (Error: %s)", report.Error)
		}
		fmt.Println()
	}
}
