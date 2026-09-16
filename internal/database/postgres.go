package database

import (
	"fmt"
	"log"

	"compro-tromagnon-service/internal/artist"
	"compro-tromagnon-service/internal/auth"
	"compro-tromagnon-service/internal/config"
	"compro-tromagnon-service/internal/hero"
	"compro-tromagnon-service/internal/live"
	"compro-tromagnon-service/internal/news"
	"compro-tromagnon-service/internal/release"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully")

	// Kolom genre dan format sudah di-migrate ke jsonb

	// Auto-migrate tables
	err = db.AutoMigrate(
		&artist.Artist{},
		&release.Release{},
		&news.News{},
		&live.LiveEvent{},
		&hero.HeroBanner{},
		&auth.Admin{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run database migration: %w", err)
	}

	log.Println("Database migration completed successfully")

	// Run initial seeds
	if err := Seed(db, cfg); err != nil {
		log.Printf("Warning: Seeder warning: %v\n", err)
	}

	return db, nil
}
