package healthcheck

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	service := NewService().Validate()

	t.Run("HealthCheck returns expected response", func(t *testing.T) {
		ctx := context.Background()

		expectedMessage := "Invoice System - Server up and running"
		expectedVersion := "v1.0.0"

		res, err := service.HealthCheck(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedMessage, res.Message)
		assert.Equal(t, expectedVersion, res.Version)
		assert.WithinDuration(t, time.Now(), parseTimeRFC1123(res.ServerTime), time.Second)
	})
}

func parseTimeRFC1123(value string) time.Time {
	parsedTime, _ := time.Parse(time.RFC1123, value)
	return parsedTime
}
