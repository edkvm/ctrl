package main

import (
	"context"

	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"



	"github.com/edkvm/ctrl"
	"github.com/edkvm/ctrl/administrating"
	"github.com/edkvm/ctrl/inmem"
	"github.com/edkvm/ctrl/invoke"
	"github.com/edkvm/ctrl/invoking"
	"github.com/edkvm/ctrl/packing"
	"github.com/edkvm/ctrl/pkg/gitsrv"

)

type ConfigEnv struct {
	Port int
	RootDir string
}

func (c *ConfigEnv) GitDir() string {
	return fmt.Sprintf("%s/git", c.RootDir)
}

func main() {

	cfg, err := loadConfigEnv()
	if err != nil {

	}

	logger := NewLogger()


	serviceLoc := ctrl.NewServeLoc(cfg.RootDir)
	mux := http.NewServeMux()



	actionRepo := inmem.NewActionRepo()
	schedRepo := inmem.NewScheduleRepo()
	statsRepo := inmem.NewStatsRepo()
	actionTimer := invoke.NewActionTimer()
	actionProvider := invoke.NewActionProvider(serviceLoc)
	invkService := invoking.NewService(actionRepo, schedRepo, statsRepo, actionTimer, actionProvider)
	invokeHandler := invoking.MakeHandler(invkService)
	mux.Handle("/invoking/v1/", invokeHandler)

	actionPacker := packing.NewActionPack(serviceLoc)
	adminService := administrating.NewService(
		actionRepo,
		schedRepo,
		statsRepo,
		actionPacker,
		administrating.NewEventHandler(
			invkService.AddActionSchedule,
			invkService.RemoveActionSchedule,
		),
	)
	adminHandler := administrating.MakeHandler(adminService)
	mux.Handle("/admin/v1/", adminHandler)


	// Expose git endpoints
	mux.Handle("/git/", gitsrv.GitServer(cfg.GitDir(), "/git/", gitsrv.NewEventHadler(adminService.ActionCodeModified)))

	logger.Log("level", "info", "msg", fmt.Sprintf("staring on port %v", cfg.Port))

	srv := &http.Server{
		Handler:      NewLoggerMiddelware(logger, mux),
		Addr:         fmt.Sprintf(":%v", cfg.Port),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Log("level", "error", "msg", err)
		}
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	<- term

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		os.Exit(1)
		return
	}
	os.Exit(0)
}

func loadConfigEnv() (*ConfigEnv,error) {
	cfg := &ConfigEnv{}

	flag.IntVar(&cfg.Port,"port",6060,"specify the port, defaults to 6060")
	flag.StringVar(&cfg.RootDir, "dir", "/usr/local/var/ctrl", "specify a directory for the service")
	flag.Parse()

	return cfg, nil
}