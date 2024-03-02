package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"spotle-backend/model"

	"github.com/gorilla/mux"
)

func (a *App) GetArtists(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}

	if start < 0 {
		start = 0
	}

	artists, err := model.GetArtists(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, artists)
}

func (a *App) GetArtist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid recipe ID")
		return
	}

	recipe := model.Artist{ID: id}
	if err := recipe.GetArtist(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Artist not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, recipe)
}

func (a *App) CreateArtist(w http.ResponseWriter, r *http.Request) {
	var recipe model.Artist

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recipe); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	defer r.Body.Close()

	if err := recipe.CreateArtist(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, recipe)
}

func (a *App) UpdateArtist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	var recipe model.Artist
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&recipe); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	defer r.Body.Close()
	recipe.ID = id

	if err := recipe.UpdateArtist(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, recipe)
}

func (a *App) DeleteArtist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	recipe := model.Artist{ID: id}
	if err := recipe.DeleteArtist(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "successfully deleted"})
}
