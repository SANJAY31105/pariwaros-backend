package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os" // Used to read environment variables

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- Database Connection ---
var DB *gorm.DB

func ConnectDatabase() {
	// Render provides the database connection string as an environment variable.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable not set!")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	// Auto-migrate the schema (creates tables if they don't exist)
	// In a real app, you'd use a more robust migration tool.
	database.AutoMigrate(&Family{}, &User{}, &Document{}, &Biller{}, &Bill{})

	DB = database
	fmt.Println("Database connection successfully opened")
}

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
	UserID    uint
	FamilyID  uint
	FileName  string
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
	DueDate  string // Using string for simplicity, would be time.Time in a full app
	IsPaid   bool
}


// --- API Handlers (Controllers) ---
// These functions handle incoming web requests.

func getBills(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// In a real app, you'd get the user ID from a JWT token.
	// For this demo, we'll just return some mock data.
	mockBills := []Bill{
		{Amount: 1245.00, DueDate: "2025-09-15", IsPaid: false, Biller: Biller{ProviderName: "Telangana Electricity"}},
		{Amount: 499.00, DueDate: "2025-08-20", IsPaid: true, Biller: Biller{ProviderName: "Airtel Postpaid"}},
	}
	json.NewEncoder(w).Encode(mockBills)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}


// --- Main Function ---
// This is the entry point of our application.

func main() {
	// For Render deployment, we must read the port from an environment variable.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port for local testing
	}

	// For now, we are not connecting the database in this simplified file.
	// ConnectDatabase() 

	r := mux.NewRouter()

	// API Routes
	r.HandleFunc("/api/v1/bills", getBills).Methods("GET")
	r.HandleFunc("/health", healthCheck).Methods("GET") // A simple health check for Render

	fmt.Println("Server starting on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
