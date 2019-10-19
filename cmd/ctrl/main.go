package main

import (
	"flag"
	"fmt"
	"github.com/edkvm/ctrl"
	"github.com/edkvm/ctrl/packing"
	"log"
	"net/http"

	"github.com/edkvm/ctrl/administrating"
	"github.com/edkvm/ctrl/inmem"
	"github.com/edkvm/ctrl/invoking"
	"github.com/edkvm/ctrl/invoke"
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


	mux.Handle("/git/", http.StripPrefix("/git/", gitsrv.GitServer(cfg.GitDir())))

	log.Println("starting")
	http.ListenAndServe(":6060", mux)

}

func loadConfigEnv() (*ConfigEnv,error) {
	cfg := &ConfigEnv{
		RootDir: "/usr/local/var/ctrl",
	}

	flag.IntVar(&cfg.Port,"port",6060,"specify the port, defaults to 6060")
	flag.Parse()

	return cfg, nil
}
