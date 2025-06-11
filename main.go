package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	// Sesuaikan import path ini dengan module name Anda
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/databases" // Ini mungkin perlu disesuaikan jika nama repo berubah
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/handlers"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/middlewares"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/models"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/repositories"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/routes"
	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/services"

	"github.com/RehanAthallahAzhar/shopeezy-inventory-cart/pkg/authclient" // Impor authclient
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	dbPortStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(dbPortStr)
	if err != nil {
		panic("Invalid DB_PORT in .env file or not set: " + err.Error())
	}

	dbCredential := models.Credential{
		Host:         os.Getenv("DB_HOST"),
		Username:     os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
		Port:         port,
	}

	dbInstance := databases.NewDB()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := dbInstance.Connect(ctx, &dbCredential)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	err = conn.AutoMigrate(&models.Product{}, &models.Cart{})
	if err != nil {
		panic("Failed migrate to database: " + err.Error())
	}

	e := echo.New()

	// e.Use(middlewares.Logger())
	// e.Use(middlewares.Recover())

	// --- Initialize gRPC AuthClient ---
	grpcServerAddress := "localhost:50051" // Address of your account-app gRPC server
	authClient, err := authclient.NewAuthClient(grpcServerAddress)
	if err != nil {
		log.Fatalf("Failed to create auth gRPC client: %v", err)
	}
	defer authClient.Close() // Pastikan koneksi gRPC ditutup saat aplikasi mati

	// --- Initialize Repositories ---
	productsRepo := repositories.NewProductRepository(conn)
	cartsRepo := repositories.NewCartRepository(conn, productsRepo) // Pastikan CartRepository Anda ada

	// --- Initialize Services ---
	productService := services.NewProductService(productsRepo)      // <-- INISIALISASI PRODUCT SERVICE
	cartService := services.NewCartService(cartsRepo, productsRepo) // Pastikan Anda memiliki CartService

	// --- Initialize Handlers ---
	// Teruskan repositories, authClient, dan services ke handler
	handler := handlers.NewHandler(productsRepo, cartsRepo, authClient, productService, cartService) // <-- SESUAIKAN SIGNATURE

	// --- Initialize Routes with Middleware ---
	authMiddlewareOpts := middlewares.AuthMiddlewareOptions{
		AuthClient: authClient,
	}
	routes.InitRoutes(e, handler, authMiddlewareOpts)

	echoPort := ":1323"
	e.Logger.Fatal(e.Start(echoPort))
}
