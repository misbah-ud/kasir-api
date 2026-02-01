package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"
)

// =======================
// Config
// =======================
type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	var err error

	// ===== Viper config (SESUAI MATERI) =====
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	fmt.Println("DB_CONN:", config.DBConn)

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Println("DB not connected:", err)
		db = nil
	} else {
		defer db.Close()
	}

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	// ===== Routes + Dependency Injection =====
	if db == nil {
		// Return 503 Service Unavailable jika database tidak terhubung
		http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		})
		http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		})
	} else {
		// ===== Dependency Injection (SESUAI MATERI) =====
		productRepo := repositories.NewProductRepository(db)
		productService := services.NewProductService(productRepo)
		productHandler := handlers.NewProductHandler(productService)

		// ===== Routes (SESUAI MATERI) =====
		http.HandleFunc("/api/produk", productHandler.HandleProducts)
		http.HandleFunc("/api/produk/", productHandler.HandleProductByID)
	}

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}
}
