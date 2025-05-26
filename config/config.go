package config
import "os"

type Config struct {
    BotToken       string
    MongoDBURI     string
    DBName         string
    MaxGoroutines  int
}

func LoadConfig() *Config {
    return &Config{
        BotToken:      "7930223337:AAFkF8FO4ijKNYHd4zApmY4XPEodXSJbZAI",
        MongoDBURI:    "mongodb+srv://itsmerick184:dummylove@cluster0.b3tg7rz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0",
        DBName:        "anti_gcast_bot",
        MaxGoroutines: 100,
    }
}