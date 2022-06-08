package environment

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// If other environment variable are added in the future, they can be added to this struct
type EnvironmentSettings struct {
	ConnectionString string
	PORT             string
	Heroku           bool
	Ip               string
}

// Fetches the environment variables (Postgres connection string and hosting address)
func GetEnvironmentVariables(logger *log.Logger) (EnvironmentSettings, bool) {

	env := EnvironmentSettings{}

	_, ok := os.LookupEnv("HEROKU")
	if !ok {
		env.Heroku = false
	} else {
		env.Heroku = true
	}

	ip, ok := os.LookupEnv("D20WORKOUT_API_IP")
	if !ok {
		// If we aren't hosting on heroku, we want an IP
		if env.Heroku != true {
			logger.Error("missing_environment_variable: D20_IP")
			return env, ok
		} else {
			env.Ip = ""
		}
	}

	connectionString, ok := os.LookupEnv("D20_WORKOUT_CAPI_CONNECTION_STRING")
	if !ok {
		logger.Error("missing_environment_variable: D20_WORKOUT_CAPI_CONNECTION_STRING")
		return env, ok
	}

	portString, ok := os.LookupEnv("PORT")
	if !ok {
		logger.Error("missing_environment_variable: PORT")
		return env, ok
	}

	env.ConnectionString = connectionString
	env.Ip = ip
	env.PORT = portString

	return env, true
}
