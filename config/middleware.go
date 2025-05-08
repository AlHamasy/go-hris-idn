package config

import (
	"hris-idn/helpers"
	"net/http"
)

// func AllowAll(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		session, _ := Store.Get(r, SESSION_ID)
// 		loggedIn, ok := session.Values["loggedIn"].(bool)
// 		if !ok || !loggedIn {
// 			redirectTo := url.QueryEscape(r.URL.Path)
// 			http.Redirect(w, r, "/login?redirect="+redirectTo, http.StatusSeeOther)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	}
// }

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, SESSION_ID)
		loggedIn, ok := session.Values["loggedIn"].(bool)
		if !ok || !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		isAdmin, ok := session.Values["isAdmin"].(bool)
		if !ok || !isAdmin {
			data := map[string]interface{}{
				"isAdmin": isAdmin,
			}
			helpers.RenderTemplate(w, "views/static/forbidden/forbidden.html", data)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func NonAdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, SESSION_ID)
		loggedIn, ok := session.Values["loggedIn"].(bool)
		if !ok || !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		isAdmin, ok := session.Values["isAdmin"].(bool)
		if ok && isAdmin {
			data := map[string]interface{}{
				"isAdmin": isAdmin,
			}
			helpers.RenderTemplate(w, "views/static/forbidden/forbidden.html", data)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// // Middleware untuk user yang harus login
// func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		session, _ := Store.Get(r, SESSION_ID)
// 		if session.Values["loggedIn"] != true {
// 			http.Redirect(w, r, "/login", http.StatusSeeOther)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	}
// }

// Middleware untuk guest (belum login)
func GuestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, SESSION_ID)
		if session.Values["loggedIn"] == true {
			if session.Values["isAdmin"] == true {
				http.Redirect(w, r, "/home-admin", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/home", http.StatusSeeOther)
			}
			return
		}
		next.ServeHTTP(w, r)
	}
}
