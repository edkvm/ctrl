package main

import (
	"flag"
	"fmt"
	"os"

	"net/http"

	kitlog "github.com/go-kit/kit/log"

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

	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "instance_id", 123)

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

	err = http.ListenAndServe(fmt.Sprintf(":%v", cfg.Port), NewLoggerMiddelware(logger, mux))
	if err != nil {
		logger.Log("level", "error", "msg", err)
	}

}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

type loggingResponseWriter struct {

}

func NewResponseWriterWrapper(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{w, http.StatusOK}
}

func (lrw *responseWriterWrapper) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

type loggerHandler struct {
	h http.Handler
	logger kitlog.Logger
}

func NewLoggerMiddelware(logger kitlog.Logger, handler http.Handler) http.Handler {
	return &loggerHandler{
		h: handler,
		logger: logger,
	}
}

func (l *loggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ww := NewResponseWriterWrapper(w)
	l.h.ServeHTTP(ww, r)
	l.logger.Log("method", r.Method, "uri", r.URL, "status", ww.statusCode)
}

func loadConfigEnv() (*ConfigEnv,error) {
	cfg := &ConfigEnv{
		RootDir: "/usr/local/var/ctrl",
	}

	flag.IntVar(&cfg.Port,"port",6060,"specify the port, defaults to 6060")
	flag.Parse()

	return cfg, nil
}
