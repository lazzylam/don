package handlers

import (
    "don/database"
    "don/utils"
    "context"
    "strings"

    "gopkg.in/telebot.v3"
    "go.mongodb.org/mongo-driver/bson"
)

func HandleMessage(c telebot.Context) error {
    // Skip jika pesan dari admin
    if utils.IsAdmin(c) {
        return nil
    }

    // Skip jika pesan kosong atau dari private chat
    if c.Message().Text == "" || c.Chat().Type == telebot.ChatPrivate {
        return nil
    }

    collection := database.GetCollection("group_settings")
    var settings database.GroupSettings
    
    filter := bson.M{"chat_id": c.Chat().ID}
    err := collection.FindOne(context.Background(), filter).Decode(&settings)
    
    // Jika tidak ada settings, anggap anti-gcast mati
    if err != nil || !settings.AntiGCAST {
        return nil
    }

    text := strings.ToLower(c.Message().Text)

    // Cek whitelist
    for _, w := range settings.Whitelist {
        if strings.Contains(text, strings.ToLower(w)) {
            return nil
        }
    }

    // Cek blacklist
    for _, b := range settings.Blacklist {
        if strings.Contains(text, strings.ToLower(b)) {
            _ = c.Delete()
            return nil
        }
    }

    // Cek pola GCAST
    if utils.IsGCASTMessage(text) {
        _ = c.Delete()
        return nil
    }

    return nil
}