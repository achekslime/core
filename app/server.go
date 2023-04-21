package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

const GracefulShutdown = 9 * time.Second

type ServerStartTask func(ctx context.Context)

type serverScript struct {
	startTasks []ServerStartTask

	// Время ожидания сигнала
	timerBeforeShutDown time.Duration
}

func NewServerScript() *serverScript {
	return &serverScript{timerBeforeShutDown: GracefulShutdown}
}

// Tasks запуск задач.
func (s *serverScript) Tasks(settingServerConfiguration ...ServerStartTask) *serverScript {
	s.startTasks = append(s.startTasks, settingServerConfiguration...)
	return s
}

func StartServer(server *serverScript) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	for _, startTask := range server.startTasks {
		startTask(ctx)
	}

	logrus.Info("Service Started")
	handleSignals(cancelFunc)
	time.Sleep(GracefulShutdown)
	logrus.Println("Service Finished")
}
