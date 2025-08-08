package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"net/http"

	"github.com/giulio-alfieri/toq_server/internal/core/config"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
	usermodel "github.com/giulio-alfieri/toq_server/internal/core/model/user_model"
)

func main() {
	//creates the main context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//create the infos and requestID to be used in the configuratio process
	infos := usermodel.UserInfos{}
	infos.ID = usermodel.SystemUserID
	ctx = context.WithValue(ctx, globalmodel.TokenKey, infos)
	ctx = context.WithValue(ctx, globalmodel.RequestIDKey, "configuration_process")

	//create the config object
	config := config.NewConfig(ctx)

	//load the environment variables
	config.LoadEnv()

	//initialize the log
	config.InitializeLog()

	slog.Info("TOQ_API initialized with version: ", "version", globalmodel.AppVersion)

	//initialize the database
	config.InitializeDatabase()
	defer func() {
		err := config.GetDatabase().Close()
		if err != nil {
			slog.Error("error closing mysql", "error", err)
			os.Exit(1)
		} else {
			slog.Debug("mysql closed")
		}
	}()

	// // Initialize OpenTelemetry
	shutdownOtel := config.InitializeTelemetry()
	defer shutdownOtel()

	// Initialize Activity Tracker early (before gRPC middleware)
	config.InitializeActivityTracker()

	//initialize the grpc server
	config.InitializeGRPC()

	//initialize the dependencies
	gcsClose := config.InjectDependencies()
	defer func() {
		err := gcsClose()
		if err != nil {
			slog.Error("error closing gcs", "error", err)
			os.Exit(1)
		}
	}()

	// Set user service in activity tracker after dependencies are initialized
	config.SetActivityTrackerUserService()

	//verify the database is new and should be initialized
	config.VerifyDatabase()

	//initilize the memory cache
	// config.InitializeMemoryCache()

	// initialize the goroutines
	config.InitializeGoRoutines()

	//get infos about the services and methods
	serviceQty, methodQty := config.GetInfos()
	slog.Info("Server started", "serviceQty", serviceQty, "methodQty", methodQty)

	//start the server
	if err := config.GetGRPCServer().Serve(config.GetListener()); err != nil {
		slog.Error("failed to serve", " err:", err)
		os.Exit(1)
	}

	// wait for the goroutines to finish
	config.GetWG().Wait()
}
