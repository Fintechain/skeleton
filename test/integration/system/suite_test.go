package system

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite provides a test suite for all integration tests
type IntegrationTestSuite struct {
	suite.Suite
}

// SetupSuite runs once before all tests in the suite
func (s *IntegrationTestSuite) SetupSuite() {
	// Ensure test data directory exists
	err := os.MkdirAll("./test-data", 0755)
	s.Require().NoError(err)

	// Set environment variables for testing
	os.Setenv("TEST_MODE", "integration")
}

// TearDownSuite runs once after all tests in the suite
func (s *IntegrationTestSuite) TearDownSuite() {
	// Clean up test data
	os.RemoveAll("./test-data")

	// Clean up environment variables
	os.Unsetenv("TEST_MODE")
}

// TestIntegrationSuite runs the integration test suite
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// Helper function to check if we're running in CI environment
func IsCI() bool {
	return os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != ""
}

// Helper function to check if we should run performance tests
func ShouldRunPerformanceTests() bool {
	return !testing.Short() && os.Getenv("SKIP_PERFORMANCE_TESTS") == ""
}
