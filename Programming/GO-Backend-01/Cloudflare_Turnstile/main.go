/**
Turnstile Widgets -> Hostnames = localhost, 127.0.0.1, tunnel.trycloudflare.com
cloudflared tunnel --url http://localhost:8080
**/

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

type turnstileResp struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

type loginReq struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	TurnstileToken string `json:"turnstileToken"`
}

func verifyTurnstile(token, remoteIP string) (bool, error) {
	form := url.Values{}
	// -- SECRET KEY -- //
	// form.Set("secret", "1x0000000000000000000000000000000AA")
	form.Set("secret", "0x4AA**********************Vg")
	form.Set("response", token)
	form.Set("remoteip", remoteIP)

	resp, err := http.PostForm(
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		form,
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result turnstileResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	return result.Success, nil
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}
		next(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	ok, err := verifyTurnstile(req.TurnstileToken, ip)
	if err != nil || !ok {
		http.Error(w, "captcha verification failed", http.StatusForbidden)
		return
	}

	fmt.Println("YOOOH, GOOT ITT...")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func main() {
	http.HandleFunc("/login", corsMiddleware(loginHandler))
	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)
}
