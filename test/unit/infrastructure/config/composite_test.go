package config_test

import (
	"testing"
	"time"

	domainconfig "github.com/fintechain/skeleton/internal/domain/config"
	infraconfig "github.com/fintechain/skeleton/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

type mockSource struct {
	values  map[string]interface{}
	loadErr error
}

func (m *mockSource) LoadConfig() error { return m.loadErr }
func (m *mockSource) GetValue(key string) (interface{}, bool) {
	v, ok := m.values[key]
	return v, ok
}
func (m *mockSource) GetAllValues() map[string]interface{} {
	return m.values
}

func TestCompositeConfig_InterfaceCompliance(t *testing.T) {
	var _ domainconfig.Configuration = (*infraconfig.CompositeConfig)(nil)
}

func TestCompositeConfig_LoadConfig_Precendence(t *testing.T) {
	src1 := &mockSource{values: map[string]interface{}{
		"database.host": "host1",
		"database.port": 5432,
		"feature.flag":  false,
	}}
	src2 := &mockSource{values: map[string]interface{}{
		"database.host": "host2", // should override src1
		"feature.flag":  true,    // should override src1
		"logging.level": "debug",
	}}
	composite := infraconfig.NewCompositeConfig(src1, src2)
	assert.NoError(t, composite.LoadConfig())

	// src2 should take precedence
	assert.Equal(t, "host2", composite.GetString("database.host"))
	assert.Equal(t, 5432, composite.GetIntDefault("database.port", 0))
	assert.Equal(t, true, composite.GetBoolDefault("feature.flag", false))
	assert.Equal(t, "debug", composite.GetString("logging.level"))
}

func TestCompositeConfig_Exists(t *testing.T) {
	src := &mockSource{values: map[string]interface{}{"foo": "bar"}}
	composite := infraconfig.NewCompositeConfig(src)
	_ = composite.LoadConfig()
	assert.True(t, composite.Exists("foo"))
	assert.False(t, composite.Exists("baz"))
}

func TestCompositeConfig_GetStringDefault(t *testing.T) {
	src := &mockSource{values: map[string]interface{}{"foo": "bar"}}
	composite := infraconfig.NewCompositeConfig(src)
	_ = composite.LoadConfig()
	assert.Equal(t, "bar", composite.GetStringDefault("foo", "default"))
	assert.Equal(t, "default", composite.GetStringDefault("baz", "default"))
}

func TestCompositeConfig_GetInt(t *testing.T) {
	src := &mockSource{values: map[string]interface{}{"int": 42, "str": "123", "float": 7.0}}
	composite := infraconfig.NewCompositeConfig(src)
	_ = composite.LoadConfig()
	v, err := composite.GetInt("int")
	assert.NoError(t, err)
	assert.Equal(t, 42, v)
	v, err = composite.GetInt("str")
	assert.NoError(t, err)
	assert.Equal(t, 123, v)
	v, err = composite.GetInt("float")
	assert.NoError(t, err)
	assert.Equal(t, 7, v)
	_, err = composite.GetInt("missing")
	assert.Error(t, err)
}

func TestCompositeConfig_GetBool(t *testing.T) {
	src := &mockSource{values: map[string]interface{}{"b1": true, "b2": "true", "b3": "1", "b4": "false", "b5": "0"}}
	composite := infraconfig.NewCompositeConfig(src)
	_ = composite.LoadConfig()
	v, err := composite.GetBool("b1")
	assert.NoError(t, err)
	assert.True(t, v)
	v, err = composite.GetBool("b2")
	assert.NoError(t, err)
	assert.True(t, v)
	v, err = composite.GetBool("b3")
	assert.NoError(t, err)
	assert.True(t, v)
	v, err = composite.GetBool("b4")
	assert.NoError(t, err)
	assert.False(t, v)
	v, err = composite.GetBool("b5")
	assert.NoError(t, err)
	assert.False(t, v)
	_, err = composite.GetBool("missing")
	assert.Error(t, err)
}

func TestCompositeConfig_GetDuration(t *testing.T) {
	src := &mockSource{values: map[string]interface{}{"d1": time.Second, "d2": "2s"}}
	composite := infraconfig.NewCompositeConfig(src)
	_ = composite.LoadConfig()
	d, err := composite.GetDuration("d1")
	assert.NoError(t, err)
	assert.Equal(t, time.Second, d)
	d, err = composite.GetDuration("d2")
	assert.NoError(t, err)
	assert.Equal(t, 2*time.Second, d)
	_, err = composite.GetDuration("missing")
	assert.Error(t, err)
}

func TestCompositeConfig_GetObject(t *testing.T) {
	type obj struct {
		A string
		B int
	}
	val := map[string]interface{}{"A": "foo", "B": 42}
	src := &mockSource{values: map[string]interface{}{"obj": val}}
	composite := infraconfig.NewCompositeConfig(src)
	_ = composite.LoadConfig()
	var o obj
	assert.NoError(t, composite.GetObject("obj", &o))
	assert.Equal(t, "foo", o.A)
	assert.Equal(t, 42, o.B)
	var missing obj
	assert.Error(t, composite.GetObject("missing", &missing))
}

func TestCompositeConfig_ThreadSafety(t *testing.T) {
	src := &mockSource{values: map[string]interface{}{"foo": "bar"}}
	composite := infraconfig.NewCompositeConfig(src)
	_ = composite.LoadConfig()
	done := make(chan bool)
	for i := 0; i < 20; i++ {
		go func() {
			_ = composite.GetString("foo")
			done <- true
		}()
	}
	for i := 0; i < 20; i++ {
		<-done
	}
}
