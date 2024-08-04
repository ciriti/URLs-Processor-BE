package main

import (
	"backend/internal/auth"
	"backend/internal/services"

	"github.com/sirupsen/logrus"
)

type application struct {
	authenticator auth.Authenticator
	logger        *logrus.Logger
	urlManager    services.URLManagerInterface
	taskQueue     services.TaskQueueInterface
}
