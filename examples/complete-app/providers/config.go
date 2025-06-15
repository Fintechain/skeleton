package providers

import (
	"fmt"
	"time"

	"github.com/fintechain/skeleton/internal/domain/config"
)

// CustomConfiguration demonstrates a custom configuration implementation
type CustomConfiguration struct {
	data map[string]interface{}
}

// NewCustomConfiguration creates a new custom configuration instance
func NewCustomConfiguration() config.Configuration {
	return &CustomConfiguration{
		data: map[string]interface{}{
			"app.name":       "Custom Framework App",
			"app.version":    "2.0.0",
			"database.host":  "custom-db-host",
			"database.port":  5432,
			"server.port":    8080,
			"logging.level":  "debug",
			"features.cache": true,
		},
	}
}

func (c *CustomConfiguration) GetString(key string) string {
	if val, ok := c.data[key].(string); ok {
		return val
	}
	return ""
}

func (c *CustomConfiguration) GetStringDefault(key, defaultValue string) string {
	if val := c.GetString(key); val != "" {
		return val
	}
	return defaultValue
}

func (c *CustomConfiguration) GetInt(key string) (int, error) {
	if val, ok := c.data[key].(int); ok {
		return val, nil
	}
	return 0, fmt.Errorf("key %s not found or not an int", key)
}

func (c *CustomConfiguration) GetIntDefault(key string, defaultValue int) int {
	if val, err := c.GetInt(key); err == nil {
		return val
	}
	return defaultValue
}

func (c *CustomConfiguration) GetBool(key string) (bool, error) {
	if val, ok := c.data[key].(bool); ok {
		return val, nil
	}
	return false, fmt.Errorf("key %s not found or not a bool", key)
}

func (c *CustomConfiguration) GetBoolDefault(key string, defaultValue bool) bool {
	if val, err := c.GetBool(key); err == nil {
		return val
	}
	return defaultValue
}

func (c *CustomConfiguration) GetDuration(key string) (time.Duration, error) {
	if val, ok := c.data[key].(time.Duration); ok {
		return val, nil
	}
	return 0, fmt.Errorf("key %s not found or not a duration", key)
}

func (c *CustomConfiguration) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	if val, err := c.GetDuration(key); err == nil {
		return val
	}
	return defaultValue
}

func (c *CustomConfiguration) GetObject(key string, target interface{}) error {
	if val, ok := c.data[key]; ok {
		// Simple assignment - in real implementation you'd use reflection or JSON marshaling
		if ptr, ok := target.(*interface{}); ok {
			*ptr = val
			return nil
		}
		return fmt.Errorf("unsupported target type for key %s", key)
	}
	return fmt.Errorf("key %s not found", key)
}

func (c *CustomConfiguration) Exists(key string) bool {
	_, exists := c.data[key]
	return exists
}
