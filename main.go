package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"erli.ng/gotab/api"
	"erli.ng/gotab/storage"
)

func init() {
	fmt.Println(`
				 ┓ 
	┏┓┏┓╋┏┓┣┓
	┗┫┗┛┗┗┻┗┛
	 ┛
	 `)
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

	api.CreateServer(store).Run()
}

func initLogger(verbosity int) *slog.Logger {
	level := slog.Level(verbosity)

	fmt.Println("verbosity output;", level.String())

	options := slog.HandlerOptions{
		Level: slog.Level(level),
	}

	handler := slog.NewTextHandler(os.Stdout, &options)

	logger := slog.New(handler)

	return logger
}
