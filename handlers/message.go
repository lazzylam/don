package handlers

import (
    "anti-gcast-bot/database"
    "anti-gcast-bot/utils"
    "gopkg.in/telebot.v3"
    "strings"
)

func HandleMessage(c telebot.Context) error {
    // Skip if message is from admin
    if utils.IsAdmin(c.Sender().ID) {
        return nil
    }

    // Skip if message is empty or from private chat
    if c.Message().Text == "" || c.Chat().Type == telebot.ChatPrivate {
        return nil
    }

    chatID := c.Chat().ID
    var settings database.GroupSettings
    if err := database.DB.First(&settings, "chat_id = ?", chatID).Error; err != nil {
        // If no settings found, assume anti-gcast is off
        return nil
    }

    // Check if anti-gcast is enabled
    if !settings.AntiGCAST {
        return nil
    }

    text := strings.ToLower(c.Message().Text)

    // Check whitelist
    for _, w := range settings.Whitelist {
        if strings.Contains(text, strings.ToLower(w)) {
            return nil
        }
    }

    // Check blacklist
    for _, b := range settings.Blacklist {
        if strings.Contains(text, strings.ToLower(b)) {
            _ = c.Delete()
            return nil
        }
    }

    // Check for GCAST patterns
    if utils.IsGCASTMessage(text) {
        _ = c.Delete()
        return nil
    }

    return nil
}