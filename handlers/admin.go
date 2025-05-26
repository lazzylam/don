package handlers

import (
    "don/database"
    "don/utils"
    "context"

    "gopkg.in/telebot.v3"
    "go.mongodb.org/mongo-driver/bson"
)

func HandleAddAdmin(c telebot.Context) error {
    if !utils.IsAdmin(c) {
        return c.Send("Hanya admin yang bisa menggunakan command ini!")
    }

    if c.Message().ReplyTo == nil {
        return c.Send("Balas pesan user yang ingin dijadikan admin!")
    }

    // Tambahkan sebagai admin grup
    admin := repository.Admin{
        UserID:   c.Message().ReplyTo.Sender.ID,
        ChatID:   c.Chat().ID,
        Username: c.Message().ReplyTo.Sender.Username,
    }

    collection := repository.GetCollection("admins")
    _, err := collection.InsertOne(context.Background(), admin)
    if err != nil {
        utils.LogError(err, "HandleAddAdmin")
        return c.Send("Gagal menambahkan admin!")
    }

    return c.Send("Admin grup berhasil ditambahkan!")
}