package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"erli.ng/gotab/api"
	"erli.ng/gotab/storage"
)

func init() {
	fmt.Println()
	fmt.Println("       ┓")
	fmt.Println("┏┓┏┓╋┏┓┣┓")
	fmt.Println("┗┫┗┛┗┗┻┗┛")
	fmt.Println(" ┛")
	fmt.Println()
}
func main() {

	disk := flag.String("disk", "./.tmp", "disk location")
	verbosity := flag.Int("verbosity", -4, "assign verbosity")
	flag.Parse()

	logger := initLogger(*verbosity)

	store := storage.Disk{Name: *disk, Logger: logger}

	if err := store.Validate(); err != nil {
		logger.Error("Failed to validate disk", "disk", store, "error", err)
		os.Exit(1)
	}

	srv := api.CreateServer(store)
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen:", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 10)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server Shutdown:", "error", err)
	}

	<-ctx.Done()

	logger.Info("Server shutdown gracefully")
}

func initLogger(verbosity int) *slog.Logger {
	level := slog.Level(verbosity)

	fmt.Println("verbosity:", level.String())

	options := slog.HandlerOptions{
		Level: slog.Level(level),
	}

	handler := slog.NewTextHandler(os.Stdout, &options)

	logger := slog.New(handler)

	return logger
}
