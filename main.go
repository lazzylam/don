package main

import (
    "anti-gcast-bot/config"
    "anti-gcast-bot/database"
    "anti-gcast-bot/handlers"
    "log"
    "time"

    "gopkg.in/telebot.v3"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize database
    database.InitDB()

    // Create bot
    bot, err := telebot.NewBot(telebot.Settings{
        Token:  cfg.BotToken,
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    })
    if err != nil {
        log.Fatal(err)
    }

    // Register handlers
    bot.Handle("/on", handlers.HandleOn)
    bot.Handle("/off", handlers.HandleOff)
    bot.Handle("/addwhite", handlers.HandleAddWhite)
    bot.Handle("/addbl", handlers.HandleAddBL)
    bot.Handle("/addadmin", handlers.HandleAddAdmin)
    
    // Handle all messages
    bot.Handle(telebot.OnText, handlers.HandleMessage)
    bot.Handle(telebot.OnPhoto, handlers.HandleMessage)
    bot.Handle(telebot.OnVideo, handlers.HandleMessage)
    bot.Handle(telebot.OnDocument, handlers.HandleMessage)

    log.Println("Bot is running...")
    bot.Start()
}