package utils

import "regexp"

var (
    LinkRegex    = regexp.MustCompile(`(?i)(https?://|www\.)[^\s]+`)
    MentionRegex = regexp.MustCompile(`(?i)@[a-zA-Z0-9_]+`)
    TMORegEx     = regexp.MustCompile(`(?i)\b(tmo|telegram messenger offline)\b`)
    JoinRegex    = regexp.MustCompile(`(?i)\b(join|gabung|add)\b`)
)

func IsGCASTMessage(text string) bool {
    return LinkRegex.MatchString(text) || 
           MentionRegex.MatchString(text) || 
           TMORegEx.MatchString(text) || 
           JoinRegex.MatchString(text)
}
