package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Manuel9550/d20-workout/pkg/dal"
	"github.com/Manuel9550/d20-workout/pkg/entities"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type D20Service struct {
	DM     *dal.DBManager
	logger *logrus.Logger
}

func NewService(dm *dal.DBManager, logger *logrus.Logger) D20Service {
	return D20Service{
		DM:     dm,
		logger: logger,
	}
}

func (service *D20Service) CheckUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(context.Background(), "APIEndpoint", "CheckUser")

	// Get the name of the user
	userName := chi.URLParam(r, "username")
	if userName == "" {
		service.respondWithError(w, 404, "Blank user passed")
		return
	}

	user, err := service.DM.GetUser(ctx, userName)
	if err != nil {
		resourceNotFoundError, ok := err.(*dal.ResourceNotFoundError)
		if ok {
			service.respondWithError(w, 404, resourceNotFoundError.Error())
			return
		} else {
			service.respondWithError(w, 500, "An internal error occured")
			return
		}
	} else {
		service.respondWithJSON(w, 200, user)
	}

}

func (service *D20Service) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(context.Background(), "APIEndpoint", "CreateUser")

	// Get the name of the user
	userName := chi.URLParam(r, "username")
	if userName == "" {
		service.respondWithError(w, 404, "Blank user passed")
		return
	}

	user, err := service.DM.CreateUser(ctx, userName)
	if err != nil {
		resourceDuplicateError, ok := err.(*dal.ResourceDuplicateError)
		if ok {
			service.respondWithError(w, 409, resourceDuplicateError.Error())
			return
		} else {
			service.respondWithError(w, 500, "An internal error occured")
			return
		}
	} else {
		service.respondWithJSON(w, 201, user)
		return
	}

}

func (service *D20Service) AddPoint(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(context.Background(), "APIEndpoint", "CreateUser")

	// Get the point
	var exercisePoint entities.Point
	err := json.NewDecoder(r.Body).Decode(&exercisePoint)

	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"Error Decoding JSON Point": err,
		}).Error()

		service.respondWithError(w, 400, "Invalid JSON for Point")
		return
	}

	// Take the current timestamp and add it to the point
	exercisePoint.Timestamp = time.Now().UTC()

	err = service.DM.AddUserPoint(ctx, &exercisePoint)
	if err != nil {
		resourceNotFoundError, ok := err.(*dal.ResourceNotFoundError)
		if ok {
			service.respondWithError(w, 404, resourceNotFoundError.Error())
			return
		} else {
			service.respondWithError(w, 500, "An internal error occured")
			return
		}
	}

	service.respondWithJSON(w, 200, exercisePoint)
	return

}

func (service *D20Service) GetUserPoints(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(context.Background(), "APIEndpoint", "GetUserPoints")

	// Get the name of the user
	userName := chi.URLParam(r, "username")
	if userName == "" {
		service.respondWithError(w, 404, "Blank user passed")
		return
	}

	// Get the start time
	startTimeString := chi.URLParam(r, "start")
	if startTimeString == "" {
		service.respondWithError(w, 404, "No start time entered")
		return
	}

	startTime, err := time.Parse("2006-01-02T15:04:05-0700", startTimeString)
	if err != nil {
		service.respondWithError(w, 400, "Incorrect time format for start time")
		return
	}

	// Get the end time
	endTimeString := chi.URLParam(r, "end")
	if endTimeString == "" {
		service.respondWithError(w, 404, "No end time entered")
		return
	}

	endTime, err := time.Parse("2006-01-02T15:04:05-0700", endTimeString)
	if err != nil {
		service.respondWithError(w, 400, "Incorrect time format for end time")
		return
	}

	// Convert the timestamps found to UTC
	startTime = startTime.UTC()
	endTime = endTime.UTC()

	points, err := service.DM.GetUserPoints(ctx, userName, startTime, endTime)
	if err != nil {
		resourceNotFoundError, ok := err.(*dal.ResourceNotFoundError)
		if ok {
			service.respondWithError(w, 404, resourceNotFoundError.Error())
			return
		} else {
			service.respondWithError(w, 500, "An internal error occured")
			return
		}
	}

	returnedPoints := entities.ExercisePoints{}
	returnedPoints.Username = userName
	returnedPoints.Points = points
	service.respondWithJSON(w, 200, returnedPoints)
	return

}

func (service *D20Service) respondWithError(w http.ResponseWriter, code int, message string) {
	err := service.respondWithJSON(w, code, map[string]string{"error": message})
	if err == nil {
		service.logger.WithFields(logrus.Fields{
			"err":  message,
			"code": code,
		}).Info("sent response to client")
	}
}

func (service *D20Service) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)

	if err != nil {
		service.logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed to send response to client")
	} else {
		service.logger.WithFields(logrus.Fields{
			"payload": payload,
			"code":    code,
		}).Info("sent response to caller")
	}

	return err
}
