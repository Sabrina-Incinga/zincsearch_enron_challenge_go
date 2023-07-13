package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zincsearch_enron_challenge_go/zincSearchApp/models"
	"github.com/zincsearch_enron_challenge_go/zincSearchApp/services"
)


func GetMails(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	from := r.URL.Query().Get("from")
	maxResults := r.URL.Query().Get("max")

	fromInt, err := strconv.Atoi(from)
	if err != nil {
		http.Error(w, "Error al convertir a entero, la variable from debe ser un número entero: "+err.Error(), http.StatusBadRequest)
		return
	}
	maxResultsInt, err := strconv.Atoi(maxResults)
	if err != nil {
		http.Error(w, "Error al convertir a entero, la variable max debe ser un número entero: "+err.Error(), http.StatusBadRequest)
		return
	}

	var data models.SearchRequest = models.SearchRequest{
		SearchType: "match",
		Query: models.Query{
			Term: term,
			Field: "_all",
		},
		From: fromInt,
		MaxResults: maxResultsInt,
		Source: make([]string, 0),
	}

	response, err := services.GetMailsByQuery(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}