package ymirtestsuites

import (
	"os"
	"testing"
)

func SkipIfIntegrationNotEnabled(t *testing.T) {
	if val, found := os.LookupEnv("YMIR_CI_INTEGRATION_TESTS_ENABLED"); !found || val != "true" {
		t.Skip("Skipping integration test, set YMIR_CI_INTEGRATION_TESTS_ENABLED=true to enable it.")
	}
}