package main

import (
	"flag"
	"log/slog"
	"os"

	"erli.ng/gotab/api"
	"erli.ng/gotab/storage"
)

func main() {
	options := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, &options)
	logger := slog.New(handler)

	disk := flag.String("disk", "./.tmp", "disk location")
	flag.Parse()

	store := storage.Disk{Name: *disk, Logger: logger}

	if err := store.Validate(); err != nil {
		logger.Error("Failed to validate disk", "disk", store, "error", err)
		os.Exit(1)
	}
	api.CreateServer(store).Run()

}
