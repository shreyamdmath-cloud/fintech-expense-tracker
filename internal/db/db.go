package db

import (
	"github.com/user/fintech-expense-tracker/internal/model"
	"gorm.io/driver/postgres"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"log"
)

func Init(dsn string) *gorm.DB {
	var db *gorm.DB
	var err error

	// Attempt connection: If it fails, fallback to SQLite for demonstration stability.
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("PostgreSQL connection failed: %v. Initializing with local SQLite instead.", err)
		db, err = gorm.Open(sqlite.Open("fintech_tracker.db"), &gorm.Config{})
	} else {
		log.Println("PostgreSQL connection established")
	}

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&model.User{}, &model.Group{}, &model.Expense{}, &model.ExpenseSplit{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")
	return db
}

// Ping checks if the database is reachable.
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
