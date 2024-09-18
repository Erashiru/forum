package handlers

import (
	"forum/pkg/rate"
	"net/http"
)

func (h *handlers) Routes() http.Handler {
	mux := http.NewServeMux()
	newLimiter := rate.NewRateLimiter(5, 15)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/"))))

	// Main handlers
	mux.HandleFunc("/", h.sessionChecker(h.home))
	mux.HandleFunc("/post/view", h.postView)
	mux.HandleFunc("/post/create", h.isAuthenticated(h.createPost))

	// User handlers
	mux.HandleFunc("/user/login", h.isLogedIn(h.login))
	mux.HandleFunc("/user/signup", h.isLogedIn(h.signup))
	mux.HandleFunc("/user/logout", h.isAuthenticated(h.logout))
	mux.HandleFunc("/user/posts", h.isAuthenticated(h.userPosts))
	mux.HandleFunc("/user/posts/liked", h.isAuthenticated(h.likedPosts))
	mux.HandleFunc("/user/request", h.isAuthenticated(h.moderationRequest))
	mux.HandleFunc("/user/posts/disliked", h.isAuthenticated(h.dislikedPosts))

	// Authentication via google and github
	mux.HandleFunc("/auth/google/login", h.isLogedIn(h.googlelogin))
	mux.HandleFunc("/auth/google/callback", h.isLogedIn(h.googleCallback))
	mux.HandleFunc("/auth/github/login", h.isLogedIn(h.githublogin))
	mux.HandleFunc("/auth/github/callback", h.isLogedIn(h.githubCallback))

	// Admin H
	mux.HandleFunc("/admin/requests", h.isAuthenticated(h.requestsList))
	mux.HandleFunc("/admin/promote", h.isAuthenticated(h.promoteUser))
	mux.HandleFunc("/admin/panel", h.isAuthenticated(h.reportList))
	mux.HandleFunc("/admin/report", h.isAuthenticated(h.deletePost))

	// Post H
	mux.HandleFunc("/post/report", h.isAuthenticated(h.ReportPost))
	mux.HandleFunc("/post/update/", h.isAuthenticated(h.updateUserPost))
	mux.HandleFunc("/post/delete/", h.isAuthenticated(h.deleteUserPost))

	// Rate limiter
	return h.rateLimiterMiddleware(newLimiter)(h.recoverPanic(h.logRequest(mux)))
}
