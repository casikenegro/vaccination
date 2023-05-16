package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"vaccination-server/db"
	"vaccination-server/models"
	"vaccination-server/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type VaccinationDTO struct {
	Name   string `json:"name" validate:"required"`
	Dose   int64  `json:"dose" validate:"required"`
	Date   string `json:"date" validate:"required"`
	DrugId uint   `json:"drug_id" validate:"required"`
}

type VaccinationResponse struct {
	ID     uint      `json:"id"`
	Name   string    `json:"name"`
	Dose   int64     `json:"dose"`
	Date   time.Time `json:"date"`
	DrugId uint      `json:"drug_id"`
}

func GetVaccinationsHandler(w http.ResponseWriter, r *http.Request) {
	var vaccinations []VaccinationResponse
	db.DB.Model(&models.Vaccination{}).Find(&vaccinations)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&vaccinations)
}

func GetVaccinationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var vaccination VaccinationResponse
	db.DB.Model(&models.Vaccination{}).First(&vaccination, params["id"])

	if vaccination.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Vaccination not found"))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&vaccination)
}

func UpdateVaccinationHandler(w http.ResponseWriter, r *http.Request) {
	var request = VaccinationDTO{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)
	var vaccination models.Vaccination
	db.DB.Model(&models.Vaccination{}).First(&vaccination, params["id"])

	if vaccination.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Vaccination not found"))
		return
	}
	date, error := utils.ParseStringToTime(request.Date)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("date format is not valid, valid format is YYYY-MM-DD"))
		return
	}

	err = ValidateDrug(request.DrugId, request.Dose, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db.DB.Save(&models.Vaccination{
		Id:     vaccination.Id,
		Name:   request.Name,
		Date:   date,
		Dose:   request.Dose,
		DrugId: request.DrugId,
	})
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&VaccinationResponse{
		ID:     vaccination.Id,
		Name:   request.Name,
		Date:   date,
		Dose:   request.Dose,
		DrugId: request.DrugId,
	})
}
func PostVaccinationHandler(w http.ResponseWriter, r *http.Request) {
	var request = VaccinationDTO{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	date, error := utils.ParseStringToTime(request.Date)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("date format is not valid, valid format is YYYY-MM-DD"))
		return
	}
	err = ValidateDrug(request.DrugId, request.Dose, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var vaccination = models.Vaccination{
		Name:   request.Name,
		Date:   date,
		Dose:   request.Dose,
		DrugId: request.DrugId,
	}
	create := db.DB.Create(&vaccination)
	errorCreate := create.Error

	if errorCreate != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorCreate.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&VaccinationResponse{
		ID:     vaccination.Id,
		Name:   vaccination.Name,
		Date:   vaccination.Date,
		Dose:   vaccination.Dose,
		DrugId: vaccination.DrugId,
	})
}

func DeleteVaccinationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var vaccination models.Vaccination
	db.DB.First(&vaccination, params["id"])

	if vaccination.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Vaccination not found"))
		return
	}

	db.DB.Unscoped().Delete(&vaccination)
	w.WriteHeader(http.StatusOK)
}

func validateDose(dose int64, MinDose int64, MaxDose int64) bool {
	return dose < MinDose || dose > MaxDose
}

func ValidateDrug(drugId uint, dose int64, date time.Time) error {
	var drug models.Drug
	db.DB.Model(&models.Drug{}).First(&drug, drugId)
	if drug.Id == 0 {
		return errors.New("Drug not found")
	}
	if validateDose(dose, drug.MinDose, drug.MaxDose) {
		return errors.New("Dose is not valid")
	}
	if date.Before(drug.AvaliableAt) {
		return errors.New("Date is not valid")
	}
	return nil
}
