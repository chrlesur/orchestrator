package utils

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "strings"
    "time"
)

// GenerateID génère un identifiant unique de la longueur spécifiée
func GenerateID(length int) string {
    bytes := make([]byte, length/2)
    if _, err := rand.Read(bytes); err != nil {
        panic(err)
    }
    return hex.EncodeToString(bytes)
}

// FormatDuration formate une durée en une chaîne lisible
func FormatDuration(d time.Duration) string {
    d = d.Round(time.Second)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    d -= m * time.Minute
    s := d / time.Second

    parts := []string{}
    if h > 0 {
        parts = append(parts, fmt.Sprintf("%dh", h))
    }
    if m > 0 {
        parts = append(parts, fmt.Sprintf("%dm", m))
    }
    if s > 0 || len(parts) == 0 {
        parts = append(parts, fmt.Sprintf("%ds", s))
    }

    return strings.Join(parts, "")
}

// TruncateString tronque une chaîne à la longueur maximale spécifiée
func TruncateString(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen-3] + "..."
}

// IsValidEmail vérifie si une chaîne est une adresse email valide
func IsValidEmail(email string) bool {
    // Cette implémentation est très basique et ne couvre pas tous les cas
    // Pour une validation plus robuste, vous pourriez utiliser une bibliothèque dédiée
    return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ParseKeyValuePairs parse une chaîne de paires clé-valeur
func ParseKeyValuePairs(s string) (map[string]string, error) {
    result := make(map[string]string)
    pairs := strings.Split(s, ",")
    for _, pair := range pairs {
        kv := strings.SplitN(pair, "=", 2)
        if len(kv) != 2 {
            return nil, fmt.Errorf("format invalide pour la paire clé-valeur: %s", pair)
        }
        result[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
    }
    return result, nil
}