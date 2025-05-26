package handlers

import (
    "anti-gcast-bot/database"
    "anti-gcast-bot/utils"
    "gopkg.in/telebot.v3"
)

func HandleAddAdmin(c telebot.Context) error {
    if !utils.IsAdmin(c.Sender().ID) {
        return c.Send("Hanya admin yang bisa menggunakan command ini!")
    }

    if c.Message().ReplyTo == nil {
        return c.Send("Balas pesan user yang ingin dijadikan admin!")
    }

    userID := c.Message().ReplyTo.Sender.ID
    username := c.Message().ReplyTo.Sender.Username

    user := database.User{
        UserID:   userID,
        IsAdmin:  true,
        Username: username,
    }

    if err := database.DB.Save(&user).Error; err != nil {
        utils.LogError(err, "HandleAddAdmin")
        return c.Send("Gagal menambahkan admin!")
    }

    return c.Send("Admin berhasil ditambahkan!")
}
