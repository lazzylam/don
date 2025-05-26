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
        BotToken:      "YOUR_BOT_TOKEN",
        MongoDBURI:    "mongodb+srv://myuser:mypassword@cluster0.mongodb.net/?retryWrites=true&w=majority",
        DBName:        "anti_gcast_bot",
        MaxGoroutines: 100,
    }
}