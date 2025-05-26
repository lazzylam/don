package handlers

import (
    "anti-gcast-bot/database"
    "anti-gcast-bot/utils"
    "gopkg.in/telebot.v3"
)

func HandleOn(c telebot.Context) error {
    if !utils.IsAdmin(c.Sender().ID) && c.Chat().Type != telebot.ChatPrivate {
        return nil
    }

    chatID := c.Chat().ID
    settings := database.GroupSettings{
        ChatID:    chatID,
        AntiGCAST: true,
    }

    if err := database.DB.Save(&settings).Error; err != nil {
        utils.LogError(err, "HandleOn")
        return c.Send("Gagal mengaktifkan Anti-GCAST!")
    }

    return c.Send("Anti-GCAST telah diaktifkan!")
}

func HandleOff(c telebot.Context) error {
    if !utils.IsAdmin(c.Sender().ID) && c.Chat().Type != telebot.ChatPrivate {
        return nil
    }

    chatID := c.Chat().ID
    if err := database.DB.Model(&database.GroupSettings{}).
        Where("chat_id = ?", chatID).
        Update("anti_gcast", false).Error; err != nil {
        utils.LogError(err, "HandleOff")
        return c.Send("Gagal menonaktifkan Anti-GCAST!")
    }

    return c.Send("Anti-GCAST telah dinonaktifkan!")
}

func HandleAddWhite(c telebot.Context) error {
    if !utils.IsAdmin(c.Sender().ID) {
        return c.Send("Hanya admin yang bisa menggunakan command ini!")
    }

    text := c.Message().Payload
    if text == "" {
        return c.Send("Masukkan kata kunci yang ingin di-whitelist!")
    }

    chatID := c.Chat().ID
    var settings database.GroupSettings
    if err := database.DB.FirstOrCreate(&settings, database.GroupSettings{ChatID: chatID}).Error; err != nil {
        utils.LogError(err, "HandleAddWhite")
        return c.Send("Gagal menambahkan ke whitelist!")
    }

    settings.Whitelist = append(settings.Whitelist, text)
    if err := database.DB.Save(&settings).Error; err != nil {
        utils.LogError(err, "HandleAddWhite")
        return c.Send("Gagal menambahkan ke whitelist!")
    }

    return c.Send("Kata kunci berhasil ditambahkan ke whitelist!")
}

func HandleAddBL(c telebot.Context) error {
    if !utils.IsAdmin(c.Sender().ID) {
        return c.Send("Hanya admin yang bisa menggunakan command ini!")
    }

    text := c.Message().Payload
    if text == "" {
        return c.Send("Masukkan kata kunci yang ingin di-blacklist!")
    }

    chatID := c.Chat().ID
    var settings database.GroupSettings
    if err := database.DB.FirstOrCreate(&settings, database.GroupSettings{ChatID: chatID}).Error; err != nil {
        utils.LogError(err, "HandleAddBL")
        return c.Send("Gagal menambahkan ke blacklist!")
    }

    settings.Blacklist = append(settings.Blacklist, text)
    if err := database.DB.Save(&settings).Error; err != nil {
        utils.LogError(err, "HandleAddBL")
        return c.Send("Gagal menambahkan ke blacklist!")
    }

    return c.Send("Kata kunci berhasil ditambahkan ke blacklist!")
}