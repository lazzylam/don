package config

type Config struct {
    BotToken       string
    MongoDBURI     string
    DBName         string
    MaxGoroutines  int
}

func LoadConfig() *Config {
    return &Config{
        BotToken:      "YOUR_BOT_TOKEN",
        MongoDBURI:    "mongodb://localhost:27017",
        DBName:        "anti_gcast_bot",
        MaxGoroutines: 100,
    }
}