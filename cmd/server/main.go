package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/williamjparra/go_VideoStreaming/internal/handler"
	"github.com/williamjparra/go_VideoStreaming/internal/repository"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("No se encontraron variables de entorno: %v", err)
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH no está configurado en el archivo .env")
	}

	//Inicializar la base de datos (creará app.db en la raíz)
	db, err := repository.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Error inicializando la base de datos: %v", err)
	}
	defer db.Close()
	log.Println("Base de datos SQLite conectada correctamente.")

	// iniciar repositorios y handlers
	userRepo := repository.NewUserRepository(db)
	authHandler := handler.NewAuthHandler(userRepo)

	// definir rutas
	http.HandleFunc("/api/register", authHandler.Register)
	http.HandleFunc("/api/login", authHandler.Login)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Iniciar Servidor HTTP
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT no está configurado en el archivo .env")
	}

	log.Printf("Servidor corriendo en http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
