package workflow

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type DeliveryStatus string

const (
	StatusPending   DeliveryStatus = "pending"
	StatusSent      DeliveryStatus = "sent"
	StatusDelivered DeliveryStatus = "delivered"
	StatusFailed    DeliveryStatus = "failed"
	StatusEscalated DeliveryStatus = "escalated"
)

type NotificationRequest struct {
	Message     string   `json:"message"`
	Recipients  []string `json:"recipients"`
	RequesterID string   `json:"requester_id"`
}

type DeliveryReport struct {
	Recipient string         `json:"recipient"`
	Status    DeliveryStatus `json:"status"`
	Error     string         `json:"error,omitempty"`
}

func NotificationWorkflow(ctx workflow.Context, req NotificationRequest) ([]DeliveryReport, error) {
	logger := workflow.GetLogger(ctx)
	reports := make([]DeliveryReport, len(req.Recipients))

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	for i, recipient := range req.Recipients {
		//Send email
		var sendResult string
		err := workflow.ExecuteActivity(ctx, SendEmailActivity, recipient, req.Message).Get(ctx, &sendResult)

		if err != nil {
			//Retry
			logger.Error("Failed to send email", "recipient", recipient, "error", err)
			reports[i] = DeliveryReport{
				Recipient: recipient,
				Status:    StatusFailed,
				Error:     err.Error(),
			}
			continue
		}

		//Waiting 24h
		var confirmation bool
		confirmationCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout: time.Hour * 24,
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 1,
			},
		})

		err = workflow.ExecuteActivity(confirmationCtx, WaitForConfirmationActivity, recipient, req.Message).Get(ctx, &confirmation)

		if err != nil || !confirmation {
			//Sending to manager
			logger.Warn("No confirmation received, escalating", "recipient", recipient)

			escalateCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				StartToCloseTimeout: time.Second * 30,
			})

			workflow.ExecuteActivity(escalateCtx, EscalateToManagerActivity, recipient, req.RequesterID, req.Message).Get(ctx, nil)

			reports[i] = DeliveryReport{
				Recipient: recipient,
				Status:    StatusEscalated,
			}
		} else {
			reports[i] = DeliveryReport{
				Recipient: recipient,
				Status:    StatusDelivered,
			}
		}
	}

	return reports, nil
}
