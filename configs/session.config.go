package configs

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
)

var SessionStore *session.Store
var StateStore *session.Store

func InitSession() {
	CLIENT_DOMAIN := os.Getenv("CLIENT_DOMAIN")

	SessionStore = session.New(session.Config{
		Expiration:     2 * 24 * time.Hour, // 2 days
		CookieHTTPOnly: true,
		CookieSecure:   true,
		CookieDomain:   CLIENT_DOMAIN,
		CookiePath:     "/",
	})
}

func InitStateSession() {
	CLIENT_DOMAIN := os.Getenv("CLIENT_DOMAIN")

	StateStore = session.New(session.Config{
		KeyLookup:      "cookie:session_state",
		Expiration:     5 * time.Minute, // 5 minutes
		CookieHTTPOnly: true,
		CookieSecure:   true,
		CookieDomain:   CLIENT_DOMAIN,
		CookiePath:     "/",
	})
}
