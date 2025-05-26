package main

import (
    "don/config"
    "don/handlers"
    "don/database"
    "don/utils"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "gopkg.in/telebot.v3"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize database
    database.InitDB(cfg)
    defer database.CloseDB()

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
bot.Handle("/paranoid", handlers.HandleParanoidMode)

    // Handle all messages
    bot.Handle(telebot.OnText, handlers.HandleMessage)
    bot.Handle(telebot.OnPhoto, handlers.HandleMessage)
    bot.Handle(telebot.OnVideo, handlers.HandleMessage)
    bot.Handle(telebot.OnDocument, handlers.HandleMessage)
bot.Handle(telebot.OnText, handlers.HandleParanoidMessage)
    

    // Graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-sigChan
        utils.LogInfo("Shutting down bot...")
        bot.Stop()
    }()

    log.Println("Bot is running...")
    bot.Start()
}