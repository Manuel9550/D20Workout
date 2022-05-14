package environment

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// If other environment variable are added in the future, they can be added to this struct
type EnvironmentSettings struct {
	PORT   string
	Heroku bool
	Ip     string
}

// Fetches the environment variables (Postgres connection string and hosting address)
func GetEnvironmentVariables(logger log.Logger) (EnvironmentSettings, bool) {

	env := EnvironmentSettings{}

	_, ok := os.LookupEnv("HEROKU")
	if !ok {
		env.Heroku = false
	} else {
		env.Heroku = true
	}

	ip, ok := os.LookupEnv("WAVE_API_IP")
	if !ok {
		// If we aren't hosting on heroku, we want an IP
		if env.Heroku != true {
			level.Error(logger).Log("missing_environment_variable", "D20_IP")
			return env, ok
		} else {
			env.Ip = ""
		}
	}

	portString, ok := os.LookupEnv("PORT")
	if !ok {
		level.Error(logger).Log("missing_environment_variable", "PORT")
		return env, ok
	}

	env.Ip = ip
	env.PORT = portString

	return env, true
}
