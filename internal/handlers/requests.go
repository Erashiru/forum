package handlers

import (
	"errors"
	"forum/models"
	"net/http"
)

func (h *handlers) requestsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	data.Requests, err = h.service.RequestsList()
	if err != nil {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	h.render(w, http.StatusOK, "requests.html", data)
}


func (h *handlers) moderationRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	err = h.service.ModerationRequest(data.UserID)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateRequest) {
			h.ErrorHandler(w, r, http.StatusOK, "Request done!")
			return
		}
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
