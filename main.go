package corser

import (
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"net/http"
)

const DefaultConfigKey = "cors_options"

// defaultOptions is a set of default cors options
var defaultOptions = cors.Options{
	AllowedOrigins: []string{"*"},
	AllowedMethods: []string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	},
	AllowedHeaders:   []string{"*"},
	AllowCredentials: false,
}

// Corser is an interface for cors configuration
type Corser interface {
	CorsOptions() cors.Options
	CorsHandler() func(next http.Handler) http.Handler
}

// CorsSettings defines a set of options to allow cors requests
type CorsOptions struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string `fig:"allowed_origins"`

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string `fig:"allowed_methods"`

	// AllowedHeaders is list of non-simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string `fig:"allowed_headers"`

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string `fig:"exposed_headers"`

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool `fig:"allow_credentials"`

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge int `fig:"max_age"`

	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool `fig:"options_passthrough"`

	// Debugging flag adds additional output to debug server side CORS issues
	Debug bool `fig:"debug"`
}

// corser is an implementation of Corser interface
type corser struct {
	getter kv.Getter
	once   comfig.Once
	key    string
}

// NewCorser creates a new instance of Corser
func NewCorser(getter kv.Getter, configKey string) Corser {
	return &corser{
		getter: getter,
		key:    configKey,
	}
}

// CorsOptions returns cors options from config file.
// If config file is empty, default options will be returned.
func (c *corser) CorsOptions() cors.Options {
	return c.once.Do(func() interface{} {
		var opts = defaultOptions

		raw, err := c.getter.GetStringMap(c.key)
		if err != nil {
			panic(errors.Wrap(err, "failed to get cors options"))
		}

		if len(raw) == 0 {
			return opts
		}

		newOpts := CorsOptions{}

		if err := figure.Out(&newOpts).From(raw).Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out cors options"))
		}

		return cors.Options{
			AllowedOrigins:     newOpts.AllowedOrigins,
			AllowedMethods:     newOpts.AllowedMethods,
			AllowedHeaders:     newOpts.AllowedHeaders,
			ExposedHeaders:     newOpts.ExposedHeaders,
			AllowCredentials:   newOpts.AllowCredentials,
			MaxAge:             newOpts.MaxAge,
			OptionsPassthrough: newOpts.OptionsPassthrough,
			Debug:              newOpts.Debug,
		}
	}).(cors.Options)
}

// CorsHandler returns cors handler based on cors configuration
func (c *corser) CorsHandler() func(next http.Handler) http.Handler {
	return cors.Handler(c.CorsOptions())
}
