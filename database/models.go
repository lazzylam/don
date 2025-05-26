package database

type GroupSettings struct {
    ChatID       int64 `gorm:"primaryKey"`
    AntiGCAST    bool
    Whitelist    []string `gorm:"serializer:json"`
    Blacklist    []string `gorm:"serializer:json"`
}

type User struct {
    UserID    int64 `gorm:"primaryKey"`
    IsAdmin   bool
    Username  string
}
