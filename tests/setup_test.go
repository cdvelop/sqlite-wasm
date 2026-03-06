// Package driver_test contains integration tests for the driver/ package.
package driver_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup: ensure tests run from repo root for any relative path dependencies
	os.Exit(m.Run())
}
