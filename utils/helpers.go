package utils

import (
    "anti-gcast-bot/database"
    "strconv"
)

func IsAdmin(userID int64) bool {
    var user database.User
    if err := database.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
        return false
    }
    return user.IsAdmin
}

func Int64ToString(i int64) string {
    return strconv.FormatInt(i, 10)
}