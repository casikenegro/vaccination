package router

import (
	"net/http"
	"vaccination-server/handlers"
	"vaccination-server/middleware"

	"github.com/gorilla/mux"
)

func BindRoutes(r *mux.Router) {
	r.HandleFunc("/", handlers.HomeHandler).Methods(http.MethodGet)
	// API
	r.HandleFunc("/login", handlers.LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/signup", handlers.SignUpHandler).Methods(http.MethodPost)
	a := r.PathPrefix("/api").Subrouter()
	a.Use(middleware.CheckAuthMiddleware)

	a.HandleFunc("/drugs", handlers.GetDrugsHandler).Methods(http.MethodGet)
	a.HandleFunc("/drugs/{id}", handlers.GetDrugHandler).Methods(http.MethodGet)
	a.HandleFunc("/drugs", handlers.PostDrugHandler).Methods(http.MethodPost)
	a.HandleFunc("/drugs/{id}", handlers.DeleteDrugHandler).Methods(http.MethodDelete)
	a.HandleFunc("/drugs/{id}", handlers.UpdateDrugHandler).Methods(http.MethodPut)
	a.HandleFunc("/vaccinations", handlers.GetVaccinationsHandler).Methods(http.MethodGet)
	a.HandleFunc("/vaccinations/{id}", handlers.GetVaccinationHandler).Methods(http.MethodGet)
	a.HandleFunc("/vaccinations", handlers.PostVaccinationHandler).Methods(http.MethodPost)
	a.HandleFunc("/vaccinations/{id}", handlers.DeleteVaccinationHandler).Methods(http.MethodDelete)
	a.HandleFunc("/vaccinations/{id}", handlers.UpdateVaccinationHandler).Methods(http.MethodPut)
}
