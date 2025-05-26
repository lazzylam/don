package utils

import (
    "regexp"
    "strings"
    "unicode"
)

// Detektor utama
func IsGCASTMessage(text string) bool {
    text = strings.ToLower(text)
    
    return detectTMO(text) || 
           detectMention(text) || 
           detectJoin(text) || 
           detectLink(text) ||
           detectSuspiciousPattern(text) ||
           detectHiddenUnicode(text) ||
           detectFontManipulation(text) ||
           detectSpamBehavior(text)
}

// Deteksi TMO
func detectTMO(text string) bool {
    tmoRegex := regexp.MustCompile(`(?i)\b(tmo|telegram\s+messenger\s+offline)\b`)
    return tmoRegex.MatchString(text)
}

// Deteksi mention
func detectMention(text string) bool {
    mentionRegex := regexp.MustCompile(`@[a-zA-Z0-9_]{5,32}`)
    return mentionRegex.MatchString(text)
}

// Deteksi kata join/gabung
func detectJoin(text string) bool {
    joinRegex := regexp.MustCompile(`(?i)\b(join|gabung|daftar|add|masuk)\b`)
    return joinRegex.MatchString(text)
}

// Deteksi link
func detectLink(text string) bool {
    linkRegex := regexp.MustCompile(`(?i)(https?://|www\.)[^\s]+`)
    return linkRegex.MatchString(text)
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
        matched, _ := regexp.MatchString(pattern, text)
        if matched {
            return true
        }
    }
    return false
}

// Deteksi Unicode tersembunyi
func detectHiddenUnicode(text string) bool {
    for _, r := range text {
        if r > unicode.MaxASCII {
            // Karakter non-ASCII ditemukan (mungkin Unicode palsu)
            return true
        }
    }
    return false
}

// Deteksi manipulasi font (contoh: 𝖒𝖆𝖓𝖎𝖕𝖚𝖑𝖆𝖘𝖎)
func detectFontManipulation(text string) bool {
    // Range Unicode untuk font-style manipulation
    suspiciousRanges := []*unicode.RangeTable{
        unicode.RangeTable{R16: []unicode.Range16{{0x1D400, 0x1D7FF, 1}}}, // Mathematical Alphanumeric Symbols
    }
    
    for _, r := range text {
        if unicode.IsOneOf(suspiciousRanges, r) {
            return true
        }
    }
    return false
}