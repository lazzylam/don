package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type GroupSettings struct {
    ID         primitive.ObjectID `bson:"_id,omitempty"`
    ChatID     int64             `bson:"chat_id"`
    AntiGCAST  bool              `bson:"anti_gcast"`
    Whitelist  []string          `bson:"whitelist"`
    Blacklist  []string          `bson:"blacklist"`
}

type Admin struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    UserID   int64             `bson:"user_id"`
    ChatID   int64             `bson:"chat_id"` // 0 untuk admin global
    Username string            `bson:"username"`
}