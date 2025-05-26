package handlers

import (
    "context"

    "don/database"
    "don/utils"

    "go.mongodb.org/mongo-driver/bson"
    "gopkg.in/telebot.v3"
)

// Tambahkan fungsi-fungsi admin di sini jika diperlukan
func HandleAddAdmin(c telebot.Context) error {
    if !utils.IsAdmin(c) {
        return c.Send("Hanya admin yang bisa menggunakan command ini!")
    }

    if c.Message().ReplyTo == nil {
        return c.Send("Balas pesan user yang ingin dijadikan admin!")
    }

    // Tambahkan sebagai admin grup
    admin := database.Admin{
        UserID:   c.Message().ReplyTo.Sender.ID,
        ChatID:   c.Chat().ID,
        Username: c.Message().ReplyTo.Sender.Username,
    }

    collection := database.GetCollection("admins")
    _, err := collection.InsertOne(context.Background(), admin)
    if err != nil {
        utils.LogError(err, "HandleAddAdmin")
        return c.Send("Gagal menambahkan admin!")
    }

    return c.Send("Admin grup berhasil ditambahkan!")
}