package queue

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAWSConfig_UsesInjectedLoader(t *testing.T) {
	expectedRegion := "sa-east-1"
	expectedCfg := aws.Config{Region: expectedRegion}
	called := false

	original := loadAWSConfig
	loadAWSConfig = func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
		called = true
		if len(optFns) > 0 {
			var opts config.LoadOptions
			require.NoError(t, optFns[0](&opts))
			assert.Equal(t, expectedRegion, opts.Region)
		}
		return expectedCfg, nil
	}
	t.Cleanup(func() {
		loadAWSConfig = original
	})

	cfg, err := loadAWSConfig(context.Background(), config.WithRegion(expectedRegion))

	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, expectedCfg, cfg)
}

func TestQueueLogFatalf_UsesInjectedLogger(t *testing.T) {
	called := false
	var capturedFormat string
	var capturedArgs []interface{}

	original := queueLogFatalf
	queueLogFatalf = func(format string, v ...interface{}) {
		called = true
		capturedFormat = format
		capturedArgs = v
	}
	t.Cleanup(func() {
		queueLogFatalf = original
	})

	queueLogFatalf("error: %s", "boom")

	assert.True(t, called)
	assert.Equal(t, "error: %s", capturedFormat)
	assert.Len(t, capturedArgs, 1)
	assert.Equal(t, "boom", capturedArgs[0])
}
