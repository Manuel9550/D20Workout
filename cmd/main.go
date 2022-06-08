package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Manuel9550/d20-workout/pkg/dal"
	"github.com/Manuel9550/d20-workout/pkg/environment"
	"github.com/Manuel9550/d20-workout/pkg/health"
	"github.com/Manuel9550/d20-workout/pkg/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

// Where to keep the logs
const logPath = "./../log"
const logFilePath = "./../log/D20Workout.log"

func main() {

	// Pre-App setup
	// Setting up the logger
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	logger.Debug("Logrus started!")

	// Get the environment variable we need
	env, ok := environment.GetEnvironmentVariables(logger)

	if !ok {
		os.Exit(-1)
	}

	// Only attempt to log to file if we aren't using Heroku
	if !env.Heroku {
		// If Logging folder doesn't exist, create it
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			os.Mkdir(logPath, os.ModeDir)
		}

		// Setting up the log file
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		defer logFile.Close()
		if err != nil {
			logger.Error("exit", err)
			os.Exit(-1)
		}

		// Want to write to terminal and file, if possible
		mw := io.MultiWriter(os.Stdout, logFile)

		logger.SetOutput(mw)
	}

	// Create the database connection and db manager
	dataManager, err := dal.NewDBManager(env.ConnectionString, logger)
	if err != nil {
		logger.Error("exit", err)
		os.Exit(-1)
	}
	defer dataManager.DB.Close()

	// Create the service
	service := service.NewService(dataManager, logger)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/test", health.GetTest)

	// Assigning the endpoints to the service methods. Make this more clean later!
	r.Get("/user", service.CheckUser)

	// If we are running on Heroku, it will listen on any interface
	fulladdress := env.Ip + ":" + env.PORT
	err = http.ListenAndServe(fulladdress, r)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("Shutting down")
}
