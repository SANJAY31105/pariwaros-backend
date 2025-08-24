// This is the complete backend in a single file for easy setup.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- Data Models ---
// These structs define the structure of our database tables.
type Family struct {
	gorm.Model
	Name  string
	Users []User
}

type User struct {
	gorm.Model
	PhoneNumber string `gorm:"unique"`
	FamilyID    uint
}

type Document struct {
	gorm.Model
	UserID     uint
	FamilyID   uint
	FileName   string
	StorageKey string
}

type Biller struct {
	gorm.Model
	UserID       uint
	FamilyID     uint
	ProviderName string
	ConsumerID   string
}

type Bill struct {
	gorm.Model
	BillerID uint
	Amount   float64
	DueDate  string
	IsPaid   bool
}

// --- Database Connection ---
var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL not set. Running without database for demo.")
		return
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	log.Println("Database connection successful. Migrating tables...")
	database.AutoMigrate(&Family{}, &User{}, &Document{}, &Biller{}, &Bill{})
	DB = database
	log.Println("Database migrated.")
}

// --- API Handlers (Controllers) ---
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}

func GetBillsHandler(w http.ResponseWriter, r *http.Request) {
	mockBills := []map[string]interface{}{
		{"ProviderName": "Telangana Electricity", "Amount": 1245.00, "DueDate": "2025-09-15", "IsPaid": false},
		{"ProviderName": "Airtel Postpaid", "Amount": 499.00, "DueDate": "2025-08-20", "IsPaid": true},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockBills)
}

// --- Main Function ---
func main() {
	log.Println("Starting PariwarOS server...")

	// For this local test, we will skip the database connection
	// to ensure the server starts without needing the DATABASE_URL.
	// ConnectDatabase()

	r := mux.NewRouter()

	// API Routes
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/api/v1/bills", GetBillsHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
