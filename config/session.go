package config

import "github.com/gorilla/sessions"

const SESSION_ID = "hris-idn-session"

var Store *sessions.CookieStore

func init() {
	Store = sessions.NewCookieStore([]byte("hris-idn"))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600 * 8, // 8 jam
		HttpOnly: true,
		Secure:   false, // Ubah jadi true jika pakai HTTPS
	}
}
