package main

import (
	"fmt"
	"log"
	"net/http"
)

// CORS with credentials cannot use "*" for Allow-Origin — the browser
// rejects it. So we echo back whatever Origin the request came from.
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func setCookieHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "testcookie",
		Value:    "hello",
		Path:     "/",
		HttpOnly: false, // false so the page can also show it via document.cookie
		Secure:   true,  // required when SameSite=None
		SameSite: http.SameSiteNoneMode,
	})
	fmt.Fprintln(w, "cookie set")
}

func checkCookieHandler(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("testcookie")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "no cookie received:", err)
		return
	}
	fmt.Fprintf(w, "received cookie: %s=%s\n", c.Name, c.Value)
}

func main() {
	http.HandleFunc("/set-cookie", withCORS(setCookieHandler))
	http.HandleFunc("/check-cookie", withCORS(checkCookieHandler))

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
