package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt"
	"github.com/rs/cors"
)

// Define the secret key to sign the token
var secretKey = []byte("s3cr3t_JWT_Key_123456")

// Claims struct to define what information will be included in the JWT token
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Message struct for the hello response
type Message struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

// Generate JWT token
func generateJWT() string {
	// Define the claims (e.g., username and expiration)
	claims := Claims{
		Username: "default",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
			Issuer:    "your-issuer",
		},
	}

	// Create a new token with the defined claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return ""
	}

	// Return the signed token
	return signedToken
}

// JWTMiddleware to validate the JWT token
func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if err := validateJWT(r, tokenString); err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// validateJWT function to validate the token
func validateJWT(r *http.Request, tokenString string) error {
	// Extract token from "Authorization" header
	if tokenString == "" {
		return fmt.Errorf("Missing Authorization token")
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure token is signed with the correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("Invalid or expired token")
	}

	return nil
}

// helloHandler for a test endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg := Message{Status: "success", Data: "Hello from Go backend!"}
	json.NewEncoder(w).Encode(msg)
}

// Your custom handler functions (e.g., IngestHandler, PreviewHandler)
func IngestHandler(w http.ResponseWriter, r *http.Request) {
	// Example handler logic
	w.Write([]byte("Ingest data"))
}

func PreviewHandler(w http.ResponseWriter, r *http.Request) {
	// Example handler logic
	w.Write([]byte("Preview data"))
}

func SchemaHandler(w http.ResponseWriter, r *http.Request) {
	// Example handler logic
	w.Write([]byte("Schema data"))
}

func TablesHandler(w http.ResponseWriter, r *http.Request) {
	// Example handler logic
	w.Write([]byte("Tables data"))
}

func main() {
	// Generate and print the JWT token
	token := generateJWT()
	fmt.Println("Generated JWT Token: ", token)

	// Setup routes and handlers
	mux := http.NewServeMux()

	// Public route (no JWT needed)
	mux.HandleFunc("/api/hello", helloHandler)

	// Protected routes with JWT middleware
	mux.HandleFunc("/ingest", JWTMiddleware(IngestHandler))
	mux.HandleFunc("/preview", JWTMiddleware(PreviewHandler))
	mux.HandleFunc("/schema", JWTMiddleware(SchemaHandler))
	mux.HandleFunc("/tables", JWTMiddleware(TablesHandler))

	// Static file server for frontend (e.g., React)
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs) // Serves index.html at root

	// Apply CORS settings for the frontend
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:5500"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	}).Handler(mux)

	// Start the server
	log.Println("✅ Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
