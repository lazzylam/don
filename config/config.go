package config

type Config struct {
    BotToken      string
    AdminIDs      []int64
    MaxGoroutines int
}

func LoadConfig() *Config {
    return &Config{
        BotToken:      "YOUR_BOT_TOKEN",
        AdminIDs:      []int64{123456789}, // Ganti dengan ID admin Anda
        MaxGoroutines: 100,
    }
}
