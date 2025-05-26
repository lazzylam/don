package handlers

import (
    "don/utils"
    "gopkg.in/telebot.v3"
    "time"
)

var paranoidMode = make(map[int64]bool) // Map chatID -> status

// Command untuk mengaktifkan paranoid mode
func HandleParanoidMode(c telebot.Context) error {
    if !utils.IsAdmin(c) {
        return c.Send("Hanya admin yang bisa mengaktifkan paranoid mode!")
    }

    chatID := c.Chat().ID
    paranoidMode[chatID] = !paranoidMode[chatID]

    status := "dinonaktifkan"
    if paranoidMode[chatID] {
        status = "diaktifkan"
        // Jalankan scanner real-time
        go startParanoidScanner(chatID, c.Bot())
    }

    return c.Send(fmt.Sprintf("🕵️‍♂️ Paranoid Mode %s!", status))
}

// Scanner real-time untuk paranoid mode
func startParanoidScanner(chatID int64, bot *telebot.Bot) {
    for paranoidMode[chatID] {
        // Lakukan pemeriksaan ekstra setiap 30 detik
        time.Sleep(30 * time.Second)
        
        // Contoh: Hapus semua pesan yang mengandung kata umum
        if err := deleteCommonTriggerMessages(chatID, bot); err != nil {
            utils.LogError(err, "ParanoidScanner")
        }
    }
}

// Hapus pesan dengan kata umum saat paranoid mode
func deleteCommonTriggerMessages(chatID int64, bot *telebot.Bot) error {
    commonTriggers := []string{"link", "join", "group", "chat", "add", "http"}
    
    // Contoh implementasi: Hapus pesan terakhir dan cek
    messages, err := bot.ChatByID(chatID)
    if err != nil {
        return err
    }

    for _, msg := range messages.Messages {
        text := strings.ToLower(msg.Text)
        for _, trigger := range commonTriggers {
            if strings.Contains(text, trigger) {
                _ = bot.Delete(msg) // Hapus pesan
                break
            }
        }
    }
    
    return nil
}

// Handler pesan untuk paranoid mode
func HandleParanoidMessage(c telebot.Context) error {
    if paranoidMode[c.Chat().ID] && !utils.IsAdmin(c) {
        // Blokir semua pesan non-admin
        _ = c.Delete()
        
        // Kirim peringatan
        _, _ = c.Bot().Send(
            c.Sender(),
            "🔒 Paranoid Mode aktif! Hanya admin yang bisa mengirim pesan.",
        )
        
        return nil
    }
    return nil
}