package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/torniker/codenotary/immudb"
	"github.com/torniker/codenotary/model"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleAny)
	fs := http.FileServer(http.Dir("./web/dist/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("GET /api/accounting", handleRetrieve)
	mux.HandleFunc("GET /api/accounting/count", handleCount)
	mux.HandleFunc("POST /api/accounting", handleAdd)
	println("Server is running on http://0.0.0.0:5656")
	log.Fatal(http.ListenAndServe(":5656", mux))
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	var accounting model.Accounting
	err := json.NewDecoder(r.Body).Decode(&accounting)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	immudb.New().CreateDocument(r.Context(), "accounting", accounting)
	w.Write([]byte("accounting created successfully!"))
}

func handleAny(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/dist")
}

func handleCount(w http.ResponseWriter, r *http.Request) {
	query := immudb.Query{
		Limit: 0,
	}
	result, err := immudb.New().CountDocument(r.Context(), "accounting", query)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(strconv.Itoa(result.Count)))
}

func handleRetrieve(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	perPage, err := strconv.Atoi(r.URL.Query().Get("perPage"))
	if err != nil {
		perPage = 10
	}
	result, err := immudb.New().SearchDocument(r.Context(), "accounting", immudb.Query{}, page, perPage)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	accounting, err := model.FromSearchResult(result)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(accounting)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
