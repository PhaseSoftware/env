// Package env provides convenience functions to retrieve
// environment variables according to the deployment target.
//
// The target is read from the ENV environment variable, which
// must be set to one of: PRODUCTION|STAGING|DEBUG.
package env

import (
	"encoding/base64"
	"os"
	"runtime/debug"
	"strconv"
)

type environment string

const (
	envProduction environment = "PRODUCTION"
	envStaging    environment = "STAGING"
	envDebug      environment = "DEBUG"

	envKey string = "ENV"
)

//
//	Detect Environment
//

func isEnvironment(env environment) bool {
	return Get(envKey, string(envDebug)) == string(env)
}

// True if production environment.
func IsProduction() bool {
	return isEnvironment(envProduction)
}

// True if staging environment.
func IsStaging() bool {
	return isEnvironment(envStaging)
}

// True if environment is either production or staging.
func IsLive() bool {
	return IsProduction() || IsStaging()
}

// Returns the first value if either production or staging, otherwise the second.
func IfLiveElse[T any](live T, debug T) T {
	if IsLive() {
		return live
	}

	return debug
}

// Returns the short git hash.
// Only available for go build.
func Version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value[:7]
			}
		}
	}

	return "unknown"
}

//
//	Retrieve Environment Variables
//

// Returns true if an environment variable with the given name exists and is not empty.
func Has(name string) bool {
	value, ok := os.LookupEnv(name)
	return ok && len(value) > 0
}

// Returns an environment variable or the provided fallback,
// if no variable with that name exists.
func Get(name string, fallback string) string {
	value, ok := os.LookupEnv(name)

	if ok && value != "" {
		return value
	}

	return fallback
}

// Returns an environment variable as integer or the
// provided fallback, if no variable with that name exists.
func GetInt(name string, fallback int) int {
	value, ok := os.LookupEnv(name)

	if !ok {
		return fallback
	}

	intval, err := strconv.Atoi(value)

	if err == nil {
		return intval
	}

	return fallback
}

// Returns an environment variable as base64 decoded bytes or nil
func GetBase64(name string) []byte {
	value, err := base64.StdEncoding.DecodeString(Get(name, ""))
	if err != nil || len(value) == 0 {
		return nil
	}

	return value
}

// Returns an environment variable or panics.
func MustGet(name string) string {
	value, ok := os.LookupEnv(name)

	if ok {
		return value
	}

	panic("environment variable " + name + " is not available")
}

// Panic if production and variable not present.
func MustGetProd(name string, fallback string) string {
	if IsLive() {
		return MustGet(name)
	}

	return Get(name, fallback)
}

// Panics if environment variable is not present.
func MustGetBase64(name string) []byte {
	value := GetBase64(name)
	if len(value) == 0 {
		panic("environment variable " + name + " is not available")
	}

	return value
}
