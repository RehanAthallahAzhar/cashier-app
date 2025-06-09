package databases

import (
	"context" // Import package context
	"fmt"
	"time" // Import time for PingContext

	model "github.com/RehanAthallahAzhar/shopeezy-inventory-cart/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct{}

func (p *Postgres) Connect(ctx context.Context, creds *model.Credential) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		creds.Host, creds.Username, creds.Password, creds.DatabaseName, creds.Port)

	// Gunakan context saat membuka koneksi GORM
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Ping database dengan context untuk memverifikasi koneksi awal
	sqlDB, err := dbConn.DB()
	if err != nil {
		return nil, err
	}
	// PingContext akan menghormati timeout atau pembatalan dari ctx
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database with context: %w", err)
	}

	// Set connection pool settings (optional, but good practice)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return dbConn, nil
}

func NewDB() *Postgres {
	return &Postgres{}
}

// Tambahkan context.Context ke fungsi Reset
func (p *Postgres) Reset(ctx context.Context, db *gorm.DB, table string) error {
	// Gunakan WithContext(ctx) untuk memastikan transaksi menghormati context
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("TRUNCATE " + table).Error; err != nil {
			return err
		}

		if err := tx.Exec("ALTER SEQUENCE " + table + "_id_seq RESTART WITH 1").Error; err != nil {
			return err
		}

		return nil
	})
}
