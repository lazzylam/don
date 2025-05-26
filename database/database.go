package database

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "log"
)

var DB *gorm.DB

func InitDB() {
    var err error
    DB, err = gorm.Open(sqlite.Open("anti-gcast.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto migrate models
    err = DB.AutoMigrate(&GroupSettings{}, &User{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
}
