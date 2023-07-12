package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/zincsearch_enron_challenge_go/zincSearchApp/controllers"
	"github.com/zincsearch_enron_challenge_go/zincSearchApp/zincsearchIndexer"
)

func main(){
	err := zincsearchIndexer.RunFilesIndexer()
	if err != nil {
		fmt.Printf("Error al crear el índice: %v\n", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(customCorsMiddleware)

	r.Get("/search", controllers.GetMails)

	http.ListenAndServe(":9000", r)
	fmt.Print("App escuchando en puerto :9000")
}

func customCorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}