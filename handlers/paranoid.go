package handlers

import (
    "don/utils"
    "fmt"
    "strings"
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

        // Log aktivitas paranoid scanner
        utils.LogInfo(fmt.Sprintf("Paranoid scanner running for chat %d", chatID))
        
        // Note: Tidak bisa mengambil pesan historis dengan Telebot v3
        // Scanner akan bekerja pada pesan baru yang masuk
    }
}

// Alternatif: Hapus pesan berdasarkan trigger umum
func checkAndDeleteSuspiciousContent(message *telebot.Message, bot *telebot.Bot) error {
    if message == nil || message.Text == "" {
        return nil
    }

    commonTriggers := []string{"link", "join", "group", "chat", "add", "http", "telegram.me", "t.me"}
    
    text := strings.ToLower(message.Text)
    for _, trigger := range commonTriggers {
        if strings.Contains(text, trigger) {
            // Hapus pesan mencurigakan
            if err := bot.Delete(message); err != nil {
                utils.LogError(err, "DeleteSuspiciousMessage")
                return err
            }
            
            // Log aktivitas
            utils.LogInfo(fmt.Sprintf("Deleted suspicious message containing '%s' from user %d", trigger, message.Sender.ID))
            break
        }
    }

    return nil
}

// Handler pesan untuk paranoid mode
func HandleParanoidMessage(c telebot.Context) error {
    chatID := c.Chat().ID
    
    // Cek apakah paranoid mode aktif
    if !paranoidMode[chatID] {
        return nil // Paranoid mode tidak aktif, lanjutkan normal
    }

    // Jika user adalah admin, izinkan pesan
    if utils.IsAdmin(c) {
        return nil
    }

    // Paranoid mode aktif dan user bukan admin
    message := c.Message()
    if message != nil {
        // Hapus pesan dari non-admin
        if err := c.Delete(); err != nil {
            utils.LogError(err, "DeleteParanoidMessage")
        }

        // Kirim peringatan ke user (private message)
        warningMsg := "🔒 Paranoid Mode aktif di grup ini! Hanya admin yang bisa mengirim pesan saat ini."
        
        // Coba kirim pesan private ke user
        if _, err := c.Bot().Send(c.Sender(), warningMsg); err != nil {
            // Jika gagal kirim private, kirim ke grup dan hapus setelah beberapa detik
            sent, err := c.Bot().Send(c.Chat(), fmt.Sprintf("@%s %s", c.Sender().Username, warningMsg))
            if err == nil {
                // Hapus pesan peringatan setelah 5 detik
                go func() {
                    time.Sleep(5 * time.Second)
                    c.Bot().Delete(sent)
                }()
            }
        }
        
        // Log aktivitas
        utils.LogInfo(fmt.Sprintf("Blocked message from user %d in paranoid mode (chat: %d)", c.Sender().ID, chatID))
    }

    return nil
}

// Fungsi untuk cek apakah paranoid mode aktif di chat tertentu
func IsParanoidModeActive(chatID int64) bool {
    return paranoidMode[chatID]
}

// Fungsi untuk mematikan paranoid mode secara manual
func DisableParanoidMode(chatID int64) {
    paranoidMode[chatID] = false
}

// Handler untuk pesan yang masuk saat paranoid mode (digunakan di main handler)
func HandleIncomingMessageParanoid(c telebot.Context) error {
    chatID := c.Chat().ID
    
    // Jika paranoid mode aktif, lakukan pemeriksaan ekstra
    if paranoidMode[chatID] {
        message := c.Message()
        if message != nil && !utils.IsAdmin(c) {
            // Periksa konten mencurigakan
            if err := checkAndDeleteSuspiciousContent(message, c.Bot()); err != nil {
                utils.LogError(err, "ParanoidContentCheck")
            }
        }
    }
    
    return nil
}