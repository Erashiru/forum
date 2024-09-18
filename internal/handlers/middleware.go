package handlers

import (
	"fmt"
	"forum/pkg/cookie"
	"forum/pkg/rate"
	"net/http"
)

// Middleware to apply rate limiting
func (h *handlers) rateLimiterMiddleware(rl *rate.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.Allow() {
				h.ErrorHandler(w, r, http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (h *handlers) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (h *handlers) isLogedIn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cooka := cookie.GetSessionCookie(r)
		if cooka != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *handlers) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				h.app.ServerError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *handlers) isAuthenticated(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *handlers) sessionChecker(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cook := cookie.GetSessionCookie(r)
		if cook != nil {
			isValid, err := h.service.IsValidToken(cook.Value)
			if err != nil {
				h.app.ServerError(w, err)
				return
			}

			if !isValid {
				cookie.ExpireSessionCookie(w)
				h.service.DeleteSession(cook.Value)
				http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
