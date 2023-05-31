package corser

import (
	"github.com/go-chi/cors"
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
	"gitlab.com/distributed_lab/kit/kv"
	"os"
	"testing"
)

func TestCorser(t *testing.T) {
	os.Setenv("KV_VIPER_FILE", "test-config.yaml")
	defer os.Unsetenv("KV_VIPER_FILE")

	getter := kv.MustFromEnv()

	for name, test := range getTestCases() {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.Expected, NewCorser(getter, test.CfgKey).CorsOptions())
		})
	}
}

type testCase struct {
	CfgKey   string
	Expected cors.Options
}

func getTestCases() map[string]testCase {
	return map[string]testCase{
		"default_options": {"non-existent-key", defaultOptions},
		"allowed_origins": {"origins", cors.Options{AllowedOrigins: []string{"http://localhost:3000"}}},
		"allowed_methods": {"methods", cors.Options{AllowedMethods: []string{"GET", "POST"}}},
		"allowed_headers": {"headers", cors.Options{AllowedHeaders: []string{"X-Header"}}},
		"exposed_headers": {"exposed", cors.Options{ExposedHeaders: []string{"X-Header"}}},
		"full_config": {DefaultConfigKey, cors.Options{
			AllowedOrigins: []string{
				"http://localhost:3000",
				"http://localhost:3001",
				"http://localpost:3000",
			},
			AllowOriginFunc:    nil,
			AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders:     []string{"X-Header", "X-Header2", "X-Header3"},
			ExposedHeaders:     []string{"X-Header"},
			AllowCredentials:   true,
			MaxAge:             3600,
			OptionsPassthrough: true,
			Debug:              true,
		}},
	}
}
