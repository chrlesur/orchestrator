package api

import (
    "net/http"
    "crypto/subtle"
)

var apiKeys = map[string]string{
    "admin": "admin_key",
    "user":  "user_key",
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        key := r.Header.Get("X-API-Key")
        if key == "" {
            http.Error(w, "Missing API key", http.StatusUnauthorized)
            return
        }

        valid := false
        for _, validKey := range apiKeys {
            if subtle.ConstantTimeCompare([]byte(key), []byte(validKey)) == 1 {
                valid = true
                break
            }
        }

        if !valid {
            http.Error(w, "Invalid API key", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    }
}