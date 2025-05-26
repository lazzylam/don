package utils

import "log"

func LogError(err error, context string) {
    if err != nil {
        log.Printf("[ERROR] %s: %v", context, err)
    }
}

func LogInfo(message string) {
    log.Printf("[INFO] %s", message)
}