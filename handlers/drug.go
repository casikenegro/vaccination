package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"vaccination-server/db"
	"vaccination-server/models"
	"vaccination-server/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type DrugDTO struct {
	Name        string `json:"name" validate:"required"`
	Approved    bool   `json:"approved" validate:"required"`
	MinDose     int64  `json:"min_dose" validate:"required"`
	MaxDose     int64  `json:"max_dose" validate:"required"`
	AvaliableAt string `json:"avaliable_at" validate:"required"`
}

type DrugResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Approved    bool      `json:"approved"`
	MinDose     int64     `json:"min_dose"`
	MaxDose     int64     `json:"max_dose"`
	AvaliableAt time.Time `json:"avaliable_at"`
}

func GetDrugsHandler(w http.ResponseWriter, r *http.Request) {
	var drugs []DrugResponse
	db.DB.Model(&models.Drug{}).Find(&drugs)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&drugs)
}

func GetDrugHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var drug DrugResponse
	db.DB.Model(&models.Drug{}).First(&drug, params["id"])

	if drug.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("drug not found"))
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&drug)
}

func UpdateDrugHandler(w http.ResponseWriter, r *http.Request) {
	var request = DrugDTO{}
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
	var drug models.Drug
	db.DB.Model(&models.Drug{}).First(&drug, params["id"])

	if drug.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("drug not found"))
		return
	}
	date, error := utils.ParseStringToTime(request.AvaliableAt)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("date format is not valid, valid format is YYYY-MM-DD"))
		return
	}
	db.DB.Save(&models.Drug{
		Id:          drug.Id,
		Name:        request.Name,
		Approved:    request.Approved,
		MinDose:     request.MinDose,
		MaxDose:     request.MaxDose,
		AvaliableAt: date,
	})
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&DrugResponse{
		ID:          drug.Id,
		Name:        request.Name,
		Approved:    request.Approved,
		MinDose:     request.MinDose,
		MaxDose:     request.MaxDose,
		AvaliableAt: date,
	})
}
func PostDrugHandler(w http.ResponseWriter, r *http.Request) {
	var request = DrugDTO{}
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
	date, error := utils.ParseStringToTime(request.AvaliableAt)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("date format is not valid, valid format is YYYY-MM-DD"))
		return
	}
	var drug = models.Drug{
		Name:        request.Name,
		Approved:    request.Approved,
		MinDose:     request.MinDose,
		MaxDose:     request.MaxDose,
		AvaliableAt: date,
	}
	create := db.DB.Create(&drug)
	errorCreate := create.Error

	if errorCreate != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorCreate.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&DrugResponse{
		ID:          drug.Id,
		Name:        drug.Name,
		Approved:    drug.Approved,
		MinDose:     drug.MinDose,
		MaxDose:     drug.MaxDose,
		AvaliableAt: drug.AvaliableAt,
	})
}

func DeleteDrugHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var drug models.Drug
	db.DB.First(&drug, params["id"])

	if drug.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("drug not found"))
		return
	}

	db.DB.Unscoped().Delete(&drug)
	w.WriteHeader(http.StatusOK)
}
