package server

import (
	"context"
	"log"
	"net/http"
	"strings"
	"student-observation/internal/database"

	"gitea.com/go-chi/session"
	"gorm.io/gorm"
)

func Opts() session.Options {
	return session.Options{
		// Name of provider. Default is "memory".
		Provider: "memory",
		// Provider configuration, it's corresponding to provider.
		ProviderConfig: "",
		// Cookie name to save session ID. Default is "MacaronSession".
		CookieName: "smsgwSession",
		// Cookie path to store. Default is "/".
		CookiePath: "/",
		// GC interval time in seconds. Default is 3600.
		Gclifetime: 3600,
		// Max life time in seconds. Default is whatever GC interval time is.
		Maxlifetime: 3600,
		// Use HTTPS only. Default is false.
		Secure: false,
		// Cookie life time. Default is 0.
		CookieLifeTime: 0,
		// Cookie domain name. Default is empty.
		Domain: "",
		// Session ID length. Default is 16.
		IDLength: 16,
		// Configuration section name. Default is "session".
	}

}

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	db := s.db.DB()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user database.User
		permissions := make(map[string]bool)
		whitelist := []string{"/login", "/assets", "/favicon.ico"}
		for _, str := range whitelist {
			if strings.Contains(r.URL.Path, str) {
				next.ServeHTTP(w, r)
				return
			}
		}
		sess := session.GetSession(r)
		username := sess.Get("loggedInUser")
		if username == nil {
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}
		db.Where("username =?", username).Preload("Permissions").Preload("Groups.Permissions").First(&user)
		if username != user.Username {
			sess.Delete("loggedInUser")
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}
		for _, permission := range user.Permissions {
			permissions[permission.Name] = true
		}
		for _, group := range user.Groups {
			//log.Println("Group name: " + group.Name)
			for _, permission := range group.Permissions {
				permissions[permission.Name] = true
			}
		}
		for key := range permissions {
			log.Println(key)
		}
		//log.Println("AuthMiddleware")
		ctx := r.Context()
		ctx = context.WithValue(ctx, "permissions", permissions)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Login(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			//redirect := r.URL.Query().Get("redirect")
			//ctx := map[string]interface{}{
			//	"redirect": redirect,

			// "csrf_value": c.Get("csrf").(string)
			//}
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte("Access Denied"))
		} else {
			var user database.User
			redirect := r.FormValue("redirect")
			username := r.FormValue("username")
			password := r.FormValue("password")
			// gorm select user where username
			db.Where("username =?", username).First(&user)
			if user.Username != username || user.Password != password {
				http.Redirect(w, r, "/login?redirect="+redirect, http.StatusSeeOther)
				return
			}
			sess := session.GetSession(r)
			sess.Set("loggedInUser", user.Username)
			http.Redirect(w, r, redirect, http.StatusSeeOther)
		}
	}
}
