package handlers

import (
	"fmt"
	"forum/models"
	"net/http"
	"strconv"
)

func (h *handlers) ReportPost(w http.ResponseWriter, r *http.Request) {
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
	if data.Role != "moderator" {
		h.ErrorHandler(w, r, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	report := *&models.Report{
		Text: r.FormValue("text"),
	}
	report.PostID, err = strconv.Atoi(r.FormValue("postID"))
	if err != nil {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	report.ModerID = data.UserID

	err = h.service.CreateReport(&report)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *handlers) reportList(w http.ResponseWriter, r *http.Request) {
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

	if data.Role != "admin" {
		h.ErrorHandler(w, r, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	data.Reports, err = h.service.ReportsList()
	if err != nil {
		fmt.Println(err)
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	h.render(w, http.StatusOK, "adminpanel.html", data)
}
