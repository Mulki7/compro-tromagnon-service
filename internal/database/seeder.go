package database

import (
	"errors"
	"log"
	"time"

	"compro-tromagnon-service/internal/auth"
	"compro-tromagnon-service/internal/config"
	"compro-tromagnon-service/internal/hero"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB, cfg *config.Config) error {
	// 1. Seed Admin
	var admin auth.Admin
	err := db.Where("email = ?", cfg.AdminEmail).First(&admin).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		newAdmin := auth.Admin{
			ID:           uuid.New(),
			Email:        cfg.AdminEmail,
			PasswordHash: string(hashedPassword),
			Name:         cfg.AdminName,
			CreatedAt:    time.Now(),
		}

		if err := db.Create(&newAdmin).Error; err != nil {
			return err
		}
		log.Printf("Admin seeded successfully: %s\n", cfg.AdminEmail)
	}

	// 2. Seed Hero Banners (hanya jika belum ada data sama sekali)
	var count int64
	db.Model(&hero.HeroBanner{}).Count(&count)
	if count == 0 {
		banners := []hero.HeroBanner{
			{
				ID:        uuid.New(),
				Title:     "TROMAGNON RECORDS",
				Subtitle:  "Independent Music Label & Creative Collective",
				ImageURL:  "",
				Order:     1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        uuid.New(),
				Title:     "Our Artists",
				Subtitle:  "Distorting Boundaries Since 2020",
				ImageURL:  "",
				Order:     2,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		for _, b := range banners {
			if err := db.Create(&b).Error; err != nil {
				return err
			}
		}
		log.Println("Hero banners seeded successfully (2 banners)")
	}

	return nil
}
