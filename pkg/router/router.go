package router

import (
	"github.com/gorilla/mux"
	"github.com/haroldlomo15/scheduler/pkg/scheduler"
	"net/http"
)

const version = "v1"


func NewRouter() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/"+version+"/getappointments", scheduler.GetAvailableAppointment)
	m.HandleFunc("/"+version+"/getscheduled", scheduler.GetScheduledAppointment)
	m.HandleFunc("/"+version+"/postappointment", scheduler.PostAppointment)
	m.HandleFunc("/"+version+"/health", health)
	return m
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("OK"))
}
