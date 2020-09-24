package internal

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func SetupGracefulShutdown(ctx context.Context, port string, engine *gin.Engine) {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}

	defer func() {
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("Server Shutdown: ", err)
		}
	}()

	signalForExit := make(chan os.Signal, 1)
	signal.Notify(signalForExit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal("Application failed", err)
		}
	}()
	logger.Info("Running application", logger.Params{"bind": port})

	stop := <-signalForExit
	logger.Info("Stop signal Received", stop)
	logger.Info("Waiting for all jobs to stop")
}

func SetupGracefulShutdownForObserver() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown ...")
	time.Sleep(time.Second * 5)
	logger.Info("Exiting  gracefully")
}
