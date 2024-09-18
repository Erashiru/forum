package handlers

import (
	"forum/models"
	"forum/pkg/cookie"
	"net/http"
	"strings"
)

func (h *handlers) signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := h.newTemplateData()
		data.Form = models.UserSignupForm{}
		data.UserID = h.GetUserIDBySession(r)
		data, err := h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		h.render(w, http.StatusOK, "signup.html", data)
	} else if r.Method == http.MethodPost {
		data := h.newTemplateData()
		data.UserID = h.GetUserIDBySession(r)
		data, err := h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		err = r.ParseForm()
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		form := &models.UserSignupForm{
			Name:          strings.ToLower(r.FormValue("username-signup")),
			Email:         strings.ToLower(r.FormValue("email-signup")),
			Password:      r.FormValue("password-signup"),
			PasswordAgain: r.PostFormValue("password-again"),
		}

		data, err = h.service.UserSignUp(data, form)
		if err != nil {
			if err == models.ErrNotValidPostForm {
				h.render(w, http.StatusBadRequest, "signup.html", data)
				return
			} else {
				h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}

		http.Redirect(w, r, "/user/login", http.StatusSeeOther) // redirect to login page
	} else {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *handlers) login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := h.newTemplateData()
		data.Form = models.UserLoginForm{}
		data.UserID = h.GetUserIDBySession(r)
		data, err := h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		h.render(w, http.StatusOK, "login.html", data)
	} else if r.Method == http.MethodPost {
		data := h.newTemplateData()
		data.UserID = h.GetUserIDBySession(r)
		data, err := h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		err = r.ParseForm()
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}

		form := &models.UserLoginForm{
			Email:    r.FormValue("email-login"),
			Password: r.FormValue("password-login"),
		}

		data, id, err := h.service.UserLogin(data, form)
		if err != nil {
			if err == models.ErrNotValidPostForm {
				h.render(w, http.StatusBadRequest, "login.html", data)
				return
			}
			h.ErrorHandler(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		err = h.service.CreateSession(w, r, id)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther) // redirect to home page
	} else {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *handlers) logout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")

	if err == nil {
		err := h.service.DeleteSession(sessionCookie.Value)
		if err != nil {
			h.app.ServerError(w, err)
			return
		}
	}

	cookie.ExpireSessionCookie(w)
	// http.SetCookie(w, &http.Cookie{
	// 	Name:    "session_token",
	// 	Value:   "",
	// 	Expires: time.Unix(0, 0),
	// 	Path:    "/",
	// 	MaxAge:  -1,
	// })
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
