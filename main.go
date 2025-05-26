package main

import (
    "anti-gcast-bot/config"
    "anti-gcast-bot/handlers"
    "anti-gcast-bot/repository"
    "anti-gcast-bot/utils"
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
    repository.InitDB(cfg)
    defer repository.CloseDB()

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