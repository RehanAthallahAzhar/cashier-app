package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/rehanazhar/cashier-app/databases"
	"github.com/rehanazhar/cashier-app/handlers"
	"github.com/rehanazhar/cashier-app/models"
	"github.com/rehanazhar/cashier-app/repositories"
	"github.com/rehanazhar/cashier-app/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Muat variabel lingkungan dari file .env
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	portStr := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("Invalid DB_PORT in .env file or not set: " + err.Error())
	}

	dbCredential := models.Credential{
		Host:         os.Getenv("DB_HOST"),
		Username:     os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
		Port:         port,
		// Schema:       os.Getenv("DB_SCHEMA"),
	}

	dbInstance := databases.NewDB() // Mengganti nama variabel 'db' agar tidak bentrok dengan package 'db'

	// Buat context dengan timeout untuk koneksi database
	// Ini akan membatalkan upaya koneksi jika melebihi 10 detik
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Pastikan cancel dipanggil untuk membersihkan resource context

	// Teruskan context ke fungsi Connect
	conn, err := dbInstance.Connect(ctx, &dbCredential) // <--- PERUBAHAN DI SINI
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	err = conn.AutoMigrate(&models.Product{}, &models.Cart{}) // &models.User{}, &models.Session{},
	if err != nil {
		panic("Failed migrate to database: " + err.Error())
	}

	e := echo.New()

	// usersRepo := repositories.NewUserRepository(conn)
	// sessionsRepo := repositories.NewSessionsRepository(conn)
	productsRepo := repositories.NewProductRepository(conn)
	cartsRepo := repositories.NewCartRepository(conn, productsRepo)

	handler := handlers.NewHandler(productsRepo, cartsRepo) // usersRepo, sessionsRepo coming soon

	// Inisialisasi rute
	routes.InitRoutes(e, handler)

	// Jalankan server Echo
	e.Logger.Fatal(e.Start(":1323"))
}
