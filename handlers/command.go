package handlers

import (
    "don/database"
    "don/utils"
    "context"

    "gopkg.in/telebot.v3"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func HandleOn(c telebot.Context) error {
    if !utils.IsAdmin(c) && c.Chat().Type != telebot.ChatPrivate {
        return nil
    }

    collection := repository.GetCollection("group_settings")
    filter := bson.M{"chat_id": c.Chat().ID}
    update := bson.M{"$set": bson.M{"anti_gcast": true}}
    opts := options.Update().SetUpsert(true)

    _, err := collection.UpdateOne(context.Background(), filter, update, opts)
    if err != nil {
        utils.LogError(err, "HandleOn")
        return c.Send("Gagal mengaktifkan Anti-GCAST!")
    }

    return c.Send("Anti-GCAST telah diaktifkan!")
}

func HandleOff(c telebot.Context) error {
    if !utils.IsAdmin(c) && c.Chat().Type != telebot.ChatPrivate {
        return nil
    }

    collection := database.GetCollection("group_settings")
    filter := bson.M{"chat_id": c.Chat().ID}
    update := bson.M{"$set": bson.M{"anti_gcast": false}}

    _, err := collection.UpdateOne(context.Background(), filter, update)
    if err != nil {
        utils.LogError(err, "HandleOff")
        return c.Send("Gagal menonaktifkan Anti-GCAST!")
    }

    return c.Send("Anti-GCAST telah dinonaktifkan!")
}

func HandleAddWhite(c telebot.Context) error {
    if !utils.IsAdmin(c) {
        return c.Send("Hanya admin yang bisa menggunakan command ini!")
    }

    text := c.Message().Payload
    if text == "" {
        return c.Send("Masukkan kata kunci yang ingin di-whitelist!")
    }

    collection := database.GetCollection("group_settings")
    filter := bson.M{"chat_id": c.Chat().ID}
    update := bson.M{"$addToSet": bson.M{"whitelist": text}}
    opts := options.Update().SetUpsert(true)

    _, err := database.UpdateOne(context.Background(), filter, update, opts)
    if err != nil {
        utils.LogError(err, "HandleAddWhite")
        return c.Send("Gagal menambahkan ke whitelist!")
    }

    return c.Send("Kata kunci berhasil ditambahkan ke whitelist!")
}

func HandleAddBL(c telebot.Context) error {
    if !utils.IsAdmin(c) {
        return c.Send("Hanya admin yang bisa menggunakan command ini!")
    }

    text := c.Message().Payload
    if text == "" {
        return c.Send("Masukkan kata kunci yang ingin di-blacklist!")
    }

    collection := database.GetCollection("group_settings")
    filter := bson.M{"chat_id": c.Chat().ID}
    update := bson.M{"$addToSet": bson.M{"blacklist": text}}
    opts := options.Update().SetUpsert(true)

    _, err := collection.UpdateOne(context.Background(), filter, update, opts)
    if err != nil {
        utils.LogError(err, "HandleAddBL")
        return c.Send("Gagal menambahkan ke blacklist!")
    }

    return c.Send("Kata kunci berhasil ditambahkan ke blacklist!")
}