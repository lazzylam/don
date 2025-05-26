package utils

import (
    "don/database"
    "context"
    "strconv"

    "gopkg.in/telebot.v3"
    "go.mongodb.org/mongo-driver/bson"
)

func IsAdmin(c telebot.Context) bool {
    // Periksa apakah pengguna adalah admin grup
    member, err := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
    if err == nil && (member.Role == telebot.Administrator || member.Role == telebot.Creator) {
        return true
    }

    // Periksa database untuk admin global
    collection := database.GetCollection("admins")
    filter := bson.M{"user_id": c.Sender().ID, "$or": []bson.M{
        {"chat_id": 0},
        {"chat_id": c.Chat().ID},
    }}
    
    count, err := collection.CountDocuments(context.Background(), filter)
    if err != nil {
        LogError(err, "IsAdmin")
        return false
    }

    return count > 0
}

func Int64ToString(i int64) string {
    return strconv.FormatInt(i, 10)
}