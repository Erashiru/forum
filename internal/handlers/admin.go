package handlers

import (
	"errors"
	"forum/models"
	"net/http"
	"strconv"
)

func (h *handlers) promoteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	data := h.newTemplateData()
	data.UserID = h.GetUserIDBySession(r)
	data, err := h.service.GetUserData(data)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	data.ReqUsrID, err = strconv.Atoi(r.FormValue("senderID"))
	if err != nil {
		h.ErrorHandler(w, r, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	err = h.service.UpdateUserRole(data, "moderator")
	if err != nil {
		if errors.Is(err, models.NoPermission) {
			h.ErrorHandler(w, r, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		} else {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	http.Redirect(w, r, "/admin/requests", http.StatusSeeOther)
}

func (h *handlers) deletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	data := h.newTemplateData()
	data.UserID = h.GetUserIDBySession(r)
	data, err := h.service.GetUserData(data)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if data.Role != "admin" {
		h.ErrorHandler(w, r, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	data.Post.ID, err = strconv.Atoi(r.FormValue("postID"))
	if err != nil {
		h.ErrorHandler(w, r, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	err = h.service.PostDelete(data.Post.ID)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/admin/panel", http.StatusSeeOther)
}
