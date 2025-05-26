package utils

import "regexp"

var (
    LinkRegex    = regexp.MustCompile(`(?i)(https?://|www\.)[^\s]+`)
    MentionRegex = regexp.MustCompile(`@[a-zA-Z0-9_]+`)
    TMORegEx     = regexp.MustCompile(`(?i)\b(tmo|telegram messenger offline)\b`)
    JoinRegex    = regexp.MustCompile(`(?i)\b(join|gabung|add|tmo|vcs|sange|anjing)\b`)
)