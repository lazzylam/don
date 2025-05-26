package utils

import (
    "log"
    "regexp"
    "strings"
    "time"
    "unicode"
)

// FilterConfig konfigurasi untuk filter
type FilterConfig struct {
    DeleteThreshold  float64  // Ambang batas untuk menghapus pesan (0.0 - 1.0)
    EnableAutoDelete bool     // Aktifkan penghapusan otomatis
    EnableLogging    bool     // Aktifkan logging
    WhitelistUsers   []string // Daftar user yang dikecualikan
}

// DetectionResult hasil analisis pesan
type DetectionResult struct {
    ShouldDelete  bool    // Apakah pesan harus dihapus
    Confidence    float64 // Tingkat keyakinan deteksi (0.0 - 1.0)
    DetectedTypes string  // Jenis deteksi yang triggered
    Reason        string  // Alasan penghapusan
}

// DefaultConfig memberikan konfigurasi default
func DefaultConfig() *FilterConfig {
    return &FilterConfig{
        DeleteThreshold:  0.6,    // Hapus jika confidence >= 60%
        EnableAutoDelete: true,
        EnableLogging:    true,
        WhitelistUsers:   []string{},
    }
}

// Variabel regex global
var (
    LinkRegex    = regexp.MustCompile(`(?i)(https?://|www\.)[^\s]+`)
    MentionRegex = regexp.MustCompile(`@[a-zA-Z0-9_]{5,32}`)
    TMORegEx     = regexp.MustCompile(`(?i)\b(tmo|telegram\s+messenger\s+offline)\b`)
    JoinRegex    = regexp.MustCompile(`(?i)\b(join|gabung|daftar|add|masuk|tmo|vcs|sange|anjing)\b`)
    
    // Regex untuk deteksi evasion  
    InvisibleCharsRegex = regexp.MustCompile("[\u200B\u200C\u200D\uFEFF]")
    SpacedWordsRegex    = regexp.MustCompile(`(\b\w\s?){3,}\b`)
    SuspiciousPatterns  = []*regexp.Regexp{
        regexp.MustCompile(`[a-z]+[0-9]+[a-z]+`), // Kata dengan angka (j0in)
        regexp.MustCompile(`\b\w{15,}\b`),      // Kata sangat panjang
        regexp.MustCompile(`\S\s\S\s\S\s\S`),   // Spasi tidak wajar
    }
)

// MessageFilter struktur utama untuk filtering
type MessageFilter struct {
    config *FilterConfig
}

// NewMessageFilter membuat instance baru MessageFilter
func NewMessageFilter(config *FilterConfig) *MessageFilter {
    if config == nil {
        config = DefaultConfig()
    }
    return &MessageFilter{config: config}
}

// CheckMessage fungsi utama untuk mengecek apakah pesan harus dihapus
func (mf *MessageFilter) CheckMessage(text, userID string) *DetectionResult {
    // Cek whitelist terlebih dahulu
    if mf.isWhitelisted(userID) {
        return &DetectionResult{
            ShouldDelete:  false,
            Confidence:    0.0,
            DetectedTypes: "whitelisted",
            Reason:        "User is whitelisted",
        }
    }

    // Hitung skor spam
    confidence := mf.calculateSpamScore(text)
    shouldDelete := confidence >= mf.config.DeleteThreshold && mf.config.EnableAutoDelete
    
    result := &DetectionResult{
        ShouldDelete:  shouldDelete,
        Confidence:    confidence,
        DetectedTypes: mf.getDetectionTypes(text),
        Reason:        mf.getDeleteReason(text, confidence),
    }

    // Log jika diaktifkan
    if mf.config.EnableLogging {
        mf.logDetection(userID, text, result)
    }

    return result
}

// calculateSpamScore menghitung skor spam (0.0 - 1.0)
func (mf *MessageFilter) calculateSpamScore(text string) float64 {
    var score float64 = 0.0
    
    // Deteksi dasar dengan bobot
    if LinkRegex.MatchString(text) {
        score += 0.3
    }
    if MentionRegex.MatchString(text) {
        score += 0.2
    }
    if TMORegEx.MatchString(text) {
        score += 0.4
    }
    if JoinRegex.MatchString(text) {
        score += 0.3
    }
    
    // Deteksi pola mencurigakan
    if mf.detectSuspiciousPattern(text) {
        score += 0.25
    }
    
    // Deteksi evasion techniques
    if mf.detectEvasion(text) {
        score += 0.35
    }
    
    // Deteksi Unicode tersembunyi
    if mf.detectHiddenUnicode(text) {
        score += 0.4
    }
    
    // Deteksi manipulasi font
    if mf.detectFontManipulation(text) {
        score += 0.3
    }
    
    // Pastikan skor tidak melebihi 1.0
    if score > 1.0 {
        score = 1.0
    }
    
    return score
}

// getDetectionTypes mengembalikan jenis deteksi yang triggered
func (mf *MessageFilter) getDetectionTypes(text string) string {
    detections := []string{}
    
    if LinkRegex.MatchString(text) {
        detections = append(detections, "link")
    }
    if MentionRegex.MatchString(text) {
        detections = append(detections, "mention")
    }
    if TMORegEx.MatchString(text) {
        detections = append(detections, "tmo")
    }
    if JoinRegex.MatchString(text) {
        detections = append(detections, "join_keyword")
    }
    if mf.detectSuspiciousPattern(text) {
        detections = append(detections, "suspicious_pattern")
    }
    if mf.detectEvasion(text) {
        detections = append(detections, "evasion")
    }
    if mf.detectHiddenUnicode(text) {
        detections = append(detections, "hidden_unicode")
    }
    if mf.detectFontManipulation(text) {
        detections = append(detections, "font_manipulation")
    }
    
    if len(detections) == 0 {
        return "clean"
    }
    
    return strings.Join(detections, ",")
}

// getDeleteReason memberikan alasan penghapusan
func (mf *MessageFilter) getDeleteReason(text string, confidence float64) string {
    if confidence < mf.config.DeleteThreshold {
        return "Message appears clean"
    }
    
    reasons := []string{}
    if LinkRegex.MatchString(text) {
        reasons = append(reasons, "contains suspicious links")
    }
    if TMORegEx.MatchString(text) {
        reasons = append(reasons, "contains TMO references")
    }
    if JoinRegex.MatchString(text) {
        reasons = append(reasons, "contains join/invite keywords")
    }
    if mf.detectEvasion(text) {
        reasons = append(reasons, "uses evasion techniques")
    }
    if mf.detectHiddenUnicode(text) {
        reasons = append(reasons, "contains hidden Unicode characters")
    }
    if mf.detectFontManipulation(text) {
        reasons = append(reasons, "uses font manipulation")
    }
    
    if len(reasons) == 0 {
        return "suspicious pattern detected"
    }
    
    return strings.Join(reasons, "; ")
}

// isWhitelisted cek apakah user ada di whitelist
func (mf *MessageFilter) isWhitelisted(userID string) bool {
    for _, id := range mf.config.WhitelistUsers {
        if id == userID {
            return true
        }
    }
    return false
}

// logDetection mencatat hasil deteksi
func (mf *MessageFilter) logDetection(userID, text string, result *DetectionResult) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    action := "KEEP"
    if result.ShouldDelete {
        action = "DELETE"
    }
    
    log.Printf("[SPAM_FILTER] %s | User: %s | Action: %s | Confidence: %.2f | Types: %s | Text: %q",
        timestamp, userID, action, result.Confidence, result.DetectedTypes, text)
}

// METODE DETEKSI INTERNAL

// detectSuspiciousPattern deteksi pola mencurigakan
func (mf *MessageFilter) detectSuspiciousPattern(text string) bool {
    text = strings.ToLower(text)
    
    // Cek regex patterns
    for _, pattern := range SuspiciousPatterns {
        if pattern.MatchString(text) {
            return true
        }
    }
    
    // Manual check untuk karakter berulang (karena Go tidak support backreference)
    return mf.hasRepeatedChars(text)
}

// hasRepeatedChars deteksi karakter berulang manual
func (mf *MessageFilter) hasRepeatedChars(text string) bool {
    if len(text) < 4 {
        return false
    }
    
    count := 1
    for i := 1; i < len(text); i++ {
        if text[i] == text[i-1] {
            count++
            if count >= 4 { // 4 atau lebih karakter sama berturut-turut
                return true
            }
        } else {
            count = 1
        }
    }
    return false
}

// detectEvasion deteksi upaya bypass
func (mf *MessageFilter) detectEvasion(text string) bool {
    return mf.containsInvisibleChars(text) ||
           mf.usesHomoglyphs(text) ||
           mf.hasSuspiciousSpacing(text)
}

// containsInvisibleChars deteksi karakter tidak terlihat
func (mf *MessageFilter) containsInvisibleChars(text string) bool {
    return InvisibleCharsRegex.MatchString(text)
}

// usesHomoglyphs deteksi homoglyph (Cyrillic vs Latin)
func (mf *MessageFilter) usesHomoglyphs(text string) bool {
    homoglyphs := map[rune]rune{
        'е': 'e', // Cyrillic 'е'
        'а': 'a', // Cyrillic 'а'
        'о': 'o', // Cyrillic 'о'
        'р': 'p', // Cyrillic 'р'
        'с': 'c', // Cyrillic 'с'
        'х': 'x', // Cyrillic 'х'
    }

    for _, r := range text {
        if _, ok := homoglyphs[r]; ok {
            if unicode.Is(unicode.Cyrillic, r) {
                return true
            }
        }
    }
    return false
}

// hasSuspiciousSpacing deteksi spasi mencurigakan (contoh: "j o i n")
func (mf *MessageFilter) hasSuspiciousSpacing(text string) bool {
    return SpacedWordsRegex.MatchString(text)
}

// detectHiddenUnicode deteksi karakter Unicode tersembunyi
func (mf *MessageFilter) detectHiddenUnicode(text string) bool {
    for _, r := range text {
        if r > unicode.MaxASCII {
            return true
        }
    }
    return false
}

// detectFontManipulation deteksi manipulasi font Unicode
func (mf *MessageFilter) detectFontManipulation(text string) bool {
    styledRanges := []*unicode.RangeTable{
        {
            R32: []unicode.Range32{{0x1D400, 0x1D7FF, 1}}, // Mathematical Alphanumeric
        },
        {
            R16: []unicode.Range16{{0xFE70, 0xFEFF, 1}}, // Arabic Presentation Forms
        },
    }

    for _, r := range text {
        if unicode.IsOneOf(styledRanges, r) {
            return true
        }
    }
    return false
}

// UTILITY FUNCTIONS

// AddToWhitelist menambah user ke whitelist
func (mf *MessageFilter) AddToWhitelist(userID string) {
    for _, id := range mf.config.WhitelistUsers {
        if id == userID {
            return // Sudah ada
        }
    }
    mf.config.WhitelistUsers = append(mf.config.WhitelistUsers, userID)
}

// RemoveFromWhitelist menghapus user dari whitelist
func (mf *MessageFilter) RemoveFromWhitelist(userID string) {
    for i, id := range mf.config.WhitelistUsers {
        if id == userID {
            mf.config.WhitelistUsers = append(
                mf.config.WhitelistUsers[:i], 
                mf.config.WhitelistUsers[i+1:]...)
            return
        }
    }
}

// UpdateThreshold mengubah threshold penghapusan
func (mf *MessageFilter) UpdateThreshold(threshold float64) {
    if threshold >= 0.0 && threshold <= 1.0 {
        mf.config.DeleteThreshold = threshold
    }
}

// GetStats mendapatkan statistik filter
func (mf *MessageFilter) GetStats() map[string]interface{} {
    return map[string]interface{}{
        "delete_threshold":   mf.config.DeleteThreshold,
        "auto_delete":        mf.config.EnableAutoDelete,
        "logging_enabled":    mf.config.EnableLogging,
        "whitelisted_users":  len(mf.config.WhitelistUsers),
    }
}

// FUNGSI KOMPATIBILITAS DENGAN KODE LAMA

// IsGCASTMessage fungsi kompatibilitas - menggabungkan semua deteksi
func IsGCASTMessage(text string) bool {
    filter := NewMessageFilter(DefaultConfig())
    result := filter.CheckMessage(text, "")
    return result.ShouldDelete
}

// DetectEvasion fungsi kompatibilitas
func DetectEvasion(text string) bool {
    filter := NewMessageFilter(DefaultConfig())
    return filter.detectEvasion(text)
}

// FUNGSI HELPER SEDERHANA

// ShouldDeleteMessage fungsi helper sederhana untuk cek apakah pesan harus dihapus
func ShouldDeleteMessage(text, userID string) bool {
    filter := NewMessageFilter(DefaultConfig())
    result := filter.CheckMessage(text, userID)
    return result.ShouldDelete
}

// AnalyzeMessage fungsi helper untuk analisis lengkap
func AnalyzeMessage(text, userID string) (shouldDelete bool, confidence float64, reason string) {
    filter := NewMessageFilter(DefaultConfig())
    result := filter.CheckMessage(text, userID)
    return result.ShouldDelete, result.Confidence, result.Reason
}