package main

import (
    "FoodStore-AdvProg2/handler"
    "FoodStore-AdvProg2/infrastructure/postgres"
    "FoodStore-AdvProg2/usecase"
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading from .env: %s", err)
    }

    dbHost := os.Getenv("DB")
    if dbHost == "" {
        log.Fatal("DB environment variable not set")
    }
    postgres.InitDB(dbHost)
    log.Println("Connected to PostgreSQL")

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

    server := &http.Server{
        Addr:         ":8080",
        Handler:      router,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    go func() {
        log.Println("Server is running on http://localhost:8080/")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed: %v", err)
        }
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan

    log.Println("Shutting down server...")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Server shutdown failed: %v", err)
    }
    log.Println("Server stopped")
}