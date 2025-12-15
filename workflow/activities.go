package workflow

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func SendEmailActivity(ctx context.Context, recipient, message string) (string, error) {
	// Test retry policy
	if rand.Intn(100) < 20 { // 20% err
		return "", fmt.Errorf("SMTP server unavailable for %s", recipient)
	}

	// dispatch time
	time.Sleep(time.Second * time.Duration(rand.Intn(2)+1))

	fmt.Printf("_____Email sent to: %s_____\nMessage: %s\n", recipient, message)
	return fmt.Sprintf("Email sent to %s", recipient), nil
}

func WaitForConfirmationActivity(ctx context.Context, recipient, message string) (bool, error) {
	timer := time.NewTicker(time.Second * 5)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	case <-timer.C: // Demo 5 second
		// 70% confirmation
		if rand.Intn(100) < 70 {
			fmt.Printf("_____Confirmation received from: %s_____\n", recipient)
			return true, nil
		}
		return false, nil
	}
}

func EscalateToManagerActivity(ctx context.Context, recipient, requesterID, message string) error {
	fmt.Printf("_____ESCALATION: Recipient %s didn't confirm. Requester: %s_____\nMessage: %s\n",
		recipient, requesterID, message)

	time.Sleep(time.Second * 1) // work imitation
	return nil
}
