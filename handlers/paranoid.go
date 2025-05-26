package handlers

import (
    "don/utils"
    "fmt"
    "strings"
    "gopkg.in/telebot.v3"
    "time"
)

var (
    paranoidMode = make(map[int64]bool)           // Map chatID -> status
    warnedUsers  = make(map[int64]map[int64]bool) // Map chatID -> userID -> warned
)

// Command untuk mengaktifkan paranoid mode
func HandleParanoidMode(c telebot.Context) error {
    if !utils.IsAdmin(c) {
        return c.Send("Hanya admin yang bisa mengaktifkan paranoid mode!")
    }

    chatID := c.Chat().ID
    paranoidMode[chatID] = !paranoidMode[chatID]

    if paranoidMode[chatID] {
        // Clear warning cache saat mengaktifkan
        clearWarnings(chatID)
        
        // Jalankan scanner real-time
        go startParanoidScanner(chatID, c.Bot())
        
        return c.Send("🕵️‍♂️ PARANOID MODE DIAKTIFKAN!\n\n🚨 SEMUA pesan non-admin akan DIHAPUS OTOMATIS\n⚡ Mode ini akan menghapus pesan secara real-time\n🔒 Hanya admin yang bisa mengirim pesan\n\n✅ Scanner aktif!")
    } else {
        // Clear warning cache saat menonaktifkan
        clearWarnings(chatID)
        
        return c.Send("🕵️‍♂️ Paranoid Mode DINONAKTIFKAN!\n\n✅ Pesan normal kembali diizinkan")
    }
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

// Handler pesan untuk paranoid mode - HAPUS SEMUA PESAN NON-ADMIN
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

    // PARANOID MODE: HAPUS SEMUA PESAN NON-ADMIN TANPA KECUALI
    message := c.Message()
    if message != nil {
        // HAPUS PESAN LANGSUNG
        if err := c.Delete(); err != nil {
            utils.LogError(err, "DeleteParanoidMessage")
        } else {
            utils.LogInfo(fmt.Sprintf("🗑️ DELETED: Message from user %d (paranoid mode)", c.Sender().ID))
        }

        // Beri warning sekali saja (cek apakah user sudah pernah diberi warning)
        userID := c.Sender().ID
        if !hasBeenWarned(chatID, userID) {
            markAsWarned(chatID, userID)
            
            warningMsg := "🔒 PARANOID MODE AKTIF!\n\n❌ Semua pesan Anda akan dihapus otomatis\n✅ Hanya admin yang bisa mengirim pesan\n\n⚠️ Pesan ini akan terhapus dalam 10 detik"
            
            // Kirim warning ke grup dan hapus setelah 10 detik
            sent, err := c.Bot().Send(c.Chat(), fmt.Sprintf("👤 @%s\n%s", c.Sender().Username, warningMsg))
            if err == nil {
                go func() {
                    time.Sleep(10 * time.Second)
                    c.Bot().Delete(sent)
                }()
            }
        }
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
    clearWarnings(chatID) // Clear warning cache
}

// Handler untuk pesan yang masuk saat paranoid mode - HAPUS LANGSUNG!
func HandleIncomingMessageParanoid(c telebot.Context) error {
    chatID := c.Chat().ID
    
    // Jika paranoid mode aktif, HAPUS SEMUA pesan non-admin
    if paranoidMode[chatID] && !utils.IsAdmin(c) {
        message := c.Message()
        if message != nil {
            // HAPUS PESAN LANGSUNG - TIDAK ADA TOLERANSI
            if err := c.Delete(); err != nil {
                utils.LogError(err, "ParanoidAutoDelete")
            } else {
                utils.LogInfo(fmt.Sprintf("🚨 AUTO-DELETED: Message from user %d in paranoid mode", c.Sender().ID))
            }
            
            // Juga periksa konten mencurigakan dan lakukan tindakan tambahan
            if utils.IsGCASTMessage(message.Text) {
                // Ban user jika mengirim GCAST
                member := &telebot.ChatMember{
                    User: c.Sender(),
                    RestrictedUntil: time.Now().Add(24 * time.Hour).Unix(),
                }
                if err := c.Bot().Restrict(c.Chat(), member); err != nil {
                    utils.LogError(err, "BanGCASTUser")
                } else {
                    utils.LogInfo(fmt.Sprintf("🔨 BANNED: User %d for GCAST attempt", c.Sender().ID))
                }
            }
        }
    }
    
    return nil
}

// Fungsi untuk cek apakah user sudah diberi warning
func hasBeenWarned(chatID, userID int64) bool {
    if warnedUsers[chatID] == nil {
        return false
    }
    return warnedUsers[chatID][userID]
}

// Fungsi untuk mark user sebagai sudah diberi warning
func markAsWarned(chatID, userID int64) {
    if warnedUsers[chatID] == nil {
        warnedUsers[chatID] = make(map[int64]bool)
    }
    warnedUsers[chatID][userID] = true
}

// Fungsi untuk clear warning cache ketika paranoid mode dimatikan
func clearWarnings(chatID int64) {
    delete(warnedUsers, chatID)
}