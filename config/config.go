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
        BotToken:      os.Getenv("BOT_TOKEN"),
        MongoDBURI:    os.Getenv("MONGODB_URI"),
        DBName:        "anti_gcast_bot",
        MaxGoroutines: 100,
    }
}