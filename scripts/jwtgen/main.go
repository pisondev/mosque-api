package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env dari root folder (karena kita akan run perintah dari root)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, pastikan file ada di root")
	}

	// Ambil dari .env
	tenantID := os.Getenv("TESTING_TENANT_ID")
	secret := os.Getenv("JWT_SECRET")

	if tenantID == "" || secret == "" {
		log.Fatal("ERROR: TESTING_TENANT_ID atau JWT_SECRET masih kosong di file .env")
	}

	// Kita buat Claims tiruan
	claims := jwt.MapClaims{
		"user_id":   "user-uuid-testing-123",
		"tenant_id": tenantID,
		"email":     "pison.gm.dev@gmail.com",
		"role":      "admin",
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // Berlaku 7 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal("Error generating token:", err)
	}

	fmt.Println("========================================================")
	fmt.Println("🎉 TOKEN JWT GENERATED SUCCESSFULLY 🎉")
	fmt.Println("========================================================")
	fmt.Println(tokenString)
	fmt.Println("========================================================")
}
