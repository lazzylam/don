package utils

import (
    "regexp"
    "strings"
    "unicode"
)

// Deteksi upaya bypass
func DetectEvasion(text string) bool {
    // Normalisasi teks (hapus spasi berlebihan, ubah ke lowercase)
    normalized := strings.ToLower(strings.Join(strings.Fields(text), ""))

    // Cek manipulasi
    return containsInvisibleChars(text) ||
           usesHomoglyphs(normalized) ||
           hasSuspiciousSpacing(text) ||
           usesFontManipulation(text)
}

// Deteksi karakter tidak terlihat
func containsInvisibleChars(text string) bool {
    invisibleRegex := regexp.MustCompile(`[\u200B-\u200D\uFEFF]`)
    return invisibleRegex.MatchString(text)
}

// Deteksi homoglyph (contoh: 'е' Cyrillic vs 'e' Latin)
func usesHomoglyphs(text string) bool {
    homoglyphs := map[rune]rune{
        'е': 'e', // Cyrillic 'е'
        'а': 'a', // Cyrillic 'а'
        'о': 'o', // Cyrillic 'о'
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

// Deteksi spasi mencurigakan
func hasSuspiciousSpacing(text string) bool {
    // Contoh: "j o i n", "t m o"
    spacedWords := regexp.MustCompile(`(\b\w\s?){3,}\b`)
    return spacedWords.MatchString(text)
}

// Deteksi manipulasi font Unicode
func usesFontManipulation(text string) bool {
    // Range Unicode untuk font styling
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