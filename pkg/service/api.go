package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Manuel9550/d20-workout/pkg/dal"
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
	userName := r.URL.Query().Get("user")
	if userName == "" {
		service.respondWithError(w, 404, "Blank user passed")
	}

	user, err := service.DM.GetUser(ctx, userName)
	if err != nil {
		service.respondWithError(w, 500, "An internal error occured")
	} else {
		service.respondWithJSON(w, 200, user)
	}

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