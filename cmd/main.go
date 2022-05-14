package main

import (
	"fmt"
	"net/http"

	"github.com/Manuel9550/d20-workout/pkg/health"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/test", health.GetTest)

	err := http.ListenAndServe(":4000", r)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("Shutting down")
}
