package utils

import (
    "sync"
    "time"
)

type FloodControl struct {
    userMessages map[int64][]time.Time
    mu           sync.Mutex
}

var floodControl = &FloodControl{
    userMessages: make(map[int64][]time.Time),
}

// Cek apakah user flooding
func IsFlooding(userID int64) bool {
    floodControl.mu.Lock()
    defer floodControl.mu.Unlock()

    const (
        maxMessages = 5      // Max 5 pesan
        timeWindow  = 10     // Dalam 10 detik
        banThreshold = 3     // Ban setelah 3x flooding
    )

    now := time.Now()
    messages := floodControl.userMessages[userID]

    // Hapus timestamp yang sudah lewat
    var validMessages []time.Time
    for _, t := range messages {
        if now.Sub(t) <= timeWindow*time.Second {
            validMessages = append(validMessages, t)
        }
    }

    // Cek apakah melebihi limit
    if len(validMessages) >= maxMessages {
        // Simpan count violation
        violations := incrementViolationCount(userID)
        
        if violations >= banThreshold {
            resetViolationCount(userID)
            return true // Trigger ban
        }
        return true // Trigger mute
    }

    // Tambahkan pesan baru
    validMessages = append(validMessages, now)
    floodControl.userMessages[userID] = validMessages
    
    return false
}

// Fungsi helper untuk violation count (simpan di Redis/MongoDB)
func incrementViolationCount(userID int64) int {
    // Implementasi database
    return 1 // Contoh sederhana
}

func resetViolationCount(userID int64) {
    // Implementasi database
}