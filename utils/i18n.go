package utils

var (
    langStrings = map[string]map[string]string{
        "id": {
            "warn_message":   "⚠️ Jangan kirim pesan GCAST!",
            "mute_message":  "🔇 Kamu di-mute karena melanggar aturan.",
            "ban_message":   "🚨 Kamu dibanned karena spam!",
            "flood_warning": "⏳ Jangan spam!",
        },
        "en": {
            "warn_message":   "⚠️ Don't send GCAST messages!",
            "mute_message":  "🔇 You've been muted for rule violation.",
            "ban_message":   "🚨 You've been banned for spamming!",
            "flood_warning": "⏳ Don't flood!",
        },
    }
    
    defaultLang = "en"
)

// Get localized string
func GetString(lang, key string) string {
    if strings, ok := langStrings[lang]; ok {
        if str, ok := strings[key]; ok {
            return str
        }
    }
    return langStrings[defaultLang][key] // Fallback ke default
}

// Deteksi bahasa user (contoh sederhana)
func DetectUserLanguage(userID int64) string {
    // Bisa diimplementasikan dengan:
    // 1. Simpan preference di database
    // 2. Deteksi dari pesan/user data
    return "id" // Contoh default ke Indonesia
}