package main

import (
	"github.com/edkvm/ctrl/execute"
	"log"
	"net/http"

	"github.com/edkvm/ctrl/administrating"
	"github.com/edkvm/ctrl/inmem"
	"github.com/edkvm/ctrl/invoking"
)

func main() {

	mux := http.NewServeMux()


	actionRepo := inmem.NewActionRepo()
	schedRepo := inmem.NewScheduleRepo()
	statsRepo := inmem.NewStatsRepo()
	actionTimer := execute.NewActionTimer()
	actionProvider := execute.NewActionProvider()
	invkService := invoking.NewService(actionRepo, schedRepo, statsRepo, actionTimer, actionProvider)
	invokeHandler := invoking.MakeHandler(invkService)
	mux.Handle("/invoking/v1/", invokeHandler)


	adminService := administrating.NewService(actionRepo, schedRepo, statsRepo, administrating.NewEventHandler(invkService.AddActionSchedule))
	adminHandler := administrating.MakeHandler(adminService)
	mux.Handle("/admin/v1/", adminHandler)

	log.Println("starting")
	http.ListenAndServe(":6060", mux)

}
