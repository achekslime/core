package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func handleSignals(cancelFunc context.CancelFunc) {
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-gracefulStop
	cancelFunc()
	logrus.Infof("STOP SERVICE %s \n", sig)
}
