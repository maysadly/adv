package main

import (
    "FoodStore-AdvProg2/handler"
    "FoodStore-AdvProg2/infrastructure/postgres"
    "FoodStore-AdvProg2/usecase"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading from .env: %s", err)
    }

    dbHost := os.Getenv("DB")
    postgres.InitDB(dbHost)

    repo := postgres.NewProductPostgresRepo()
    uc := usecase.NewProductUseCase(repo)
    productHandler := handler.NewProductHandler(uc)

    router := mux.NewRouter()

    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./public/admin.html")
    }).Methods("GET")

    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./public/"))))

    api := router.PathPrefix("/api").Subrouter()
    api.HandleFunc("/products", productHandler.Create).Methods("POST")
    api.HandleFunc("/products/{id}", productHandler.Get).Methods("GET")
    api.HandleFunc("/products/{id}", productHandler.Update).Methods("PUT")
    api.HandleFunc("/products/{id}", productHandler.Delete).Methods("DELETE")
    api.HandleFunc("/products", productHandler.List).Methods("GET")


    log.Println("Server is running on http://localhost:8080/")
    http.ListenAndServe(":8080", router)
}