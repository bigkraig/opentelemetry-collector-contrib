// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package signalfxexporter

import (
	"testing"
	"time"

	"github.com/open-telemetry/opentelemetry-collector/config/configcheck"
	"github.com/open-telemetry/opentelemetry-collector/config/configerror"
	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := Factory{}
	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, configcheck.ValidateConfig(cfg))
}

func TestCreateMetricsExporter(t *testing.T) {
	factory := Factory{}
	cfg := factory.CreateDefaultConfig()

	assert.Equal(t, typeStr, factory.Type())
	_, err := factory.CreateMetricsExporter(zap.NewNop(), cfg)
	assert.NoError(t, err)
}

func TestCreateTraceExporter(t *testing.T) {
	factory := Factory{}
	cfg := factory.CreateDefaultConfig()
	_, err := factory.CreateTraceExporter(zap.NewNop(), cfg)
	assert.Equal(t, configerror.ErrDataTypeIsNotSupported, err)
}

func TestCreateInstanceViaFactory(t *testing.T) {
	factory := Factory{}

	cfg := factory.CreateDefaultConfig()
	exp, err := factory.CreateMetricsExporter(
		zap.NewNop(),
		cfg)
	assert.NoError(t, err)
	assert.NotNil(t, exp)

	// Set values that don't have a valid default.
	expCfg := cfg.(*Config)
	expCfg.AccessToken = "testToken"
	expCfg.Realm = "us1"
	exp, err = factory.CreateMetricsExporter(
		zap.NewNop(),
		cfg)
	assert.NoError(t, err)
	require.NotNil(t, exp)

	assert.NoError(t, exp.Shutdown())
}

func TestFactory_CreateMetricsExporter(t *testing.T) {
	f := &Factory{}
	config := &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		AccessToken: "testToken",
		Realm:       "us1",
		Headers: map[string]string{
			"added-entry": "added value",
			"dot.test":    "test",
		},
		Timeout: 2 * time.Second,
	}

	te, err := f.CreateMetricsExporter(zap.NewNop(), config)
	assert.NoError(t, err)
	assert.NotNil(t, te)
}

func TestFactory_CreateMetricsExporterFails(t *testing.T) {
	tests := []struct {
		name         string
		config       *Config
		errorMessage string
	}{
		{
			name: "negative_duration",
			config: &Config{
				ExporterSettings: configmodels.ExporterSettings{
					TypeVal: typeStr,
					NameVal: typeStr,
				},
				AccessToken: "testToken",
				Realm:       "lab",
				Timeout:     -2 * time.Second,
			},
			errorMessage: "\"signalfx\" config cannot have a negative \"timeout\"",
		},
		{
			name: "empty_realm_and_url",
			config: &Config{
				ExporterSettings: configmodels.ExporterSettings{
					TypeVal: typeStr,
					NameVal: typeStr,
				},
			},
			errorMessage: "\"signalfx\" config requires a non-empty \"realm\" or \"url\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Factory{}
			te, err := f.CreateMetricsExporter(zap.NewNop(), tt.config)
			assert.EqualError(t, err, tt.errorMessage)
			assert.Nil(t, te)
		})
	}
}
