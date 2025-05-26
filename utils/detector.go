package utils

import (
    "regexp"
    "strings"
    "unicode"
)

// Variabel regex global (di-export jika diperlukan di package lain)
var (
    LinkRegex    = regexp.MustCompile(`(?i)(https?://|www\.)[^\s]+`)
    MentionRegex = regexp.MustCompile(`@[a-zA-Z0-9_]{5,32}`)
    TMORegEx     = regexp.MustCompile(`(?i)\b(tmo|telegram\s+messenger\s+offline)\b`)
    JoinRegex    = regexp.MustCompile(`(?i)\b(join|gabung|daftar|add|masuk|tmo|vcs|sange|anjing)\b`)
)

// Deteksi utama yang mengkombinasikan semua pola
func IsGCASTMessage(text string) bool {
    text = strings.ToLower(text)
    
    // Deteksi dasar menggunakan regex global
    if LinkRegex.MatchString(text) || 
       MentionRegex.MatchString(text) || 
       TMORegEx.MatchString(text) || 
       JoinRegex.MatchString(text) {
        return true
    }
    
    // Deteksi lanjutan
    return detectSuspiciousPattern(text) ||
           detectHiddenUnicode(text) ||
           detectFontManipulation(text)
}

// Deteksi pola mencurigakan
func detectSuspiciousPattern(text string) bool {
    patterns := []string{
        `(\w)\1{3,}`,       // Karakter berulang (contoh: "joiiiiin")
        `[a-z]+\d+[a-z]+`,  // Kata dengan angka (contoh: "j0in", "t3legram")
        `\b\w{15,}\b`,      // Kata sangat panjang
        `\S\s\S\s\S\s\S`,   // Spasi tidak wajar
    }

    for _, pattern := range patterns {
        if matched, _ := regexp.MatchString(pattern, text); matched {
            return true
        }
    }
    return false
}

// Deteksi Unicode tersembunyi
func detectHiddenUnicode(text string) bool {
    for _, r := range text {
        if r > unicode.MaxASCII {
            return true // Karakter non-ASCII ditemukan
        }
    }
    return false
}

// Deteksi manipulasi font
func detectFontManipulation(text string) bool {
    suspiciousRanges := []*unicode.RangeTable{
        {R16: []unicode.Range16{{0x1D400, 0x1D7FF, 1}}}, // Mathematical Alphanumeric
    }
    
    for _, r := range text {
        if unicode.IsOneOf(suspiciousRanges, r) {
            return true
        }
    }
    return false
}