package handlers

import (
	"forum/models"
	"net/http"
)

func (h *handlers) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" &&
		r.URL.Path != "/Golang" &&
		r.URL.Path != "/AI" &&
		r.URL.Path != "/Python" &&
		r.URL.Path != "/JavaScript" &&
		r.URL.Path != "/Algorithms" {
		h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method == http.MethodGet {
		data := h.newTemplateData()
		data.UserID = h.GetUserIDBySession(r)
		data, err := h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		data, err = h.service.MainLoader(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		h.render(w, http.StatusOK, "index.html", data)
	} else if r.Method == http.MethodPost {
		data := h.newTemplateData()

		postID := r.FormValue("postID")

		filterData := &models.FilterData{
			Filter:      r.FormValue("filter"),
			Reaction:    r.FormValue("reaction"),
			Category:    r.FormValue("category"),
			CategoryURL: r.URL.Path[1:],
		}
		data.UserID = h.GetUserIDBySession(r)
		if data.UserID == 0 {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		data, err := h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		data, err = h.service.HomeUpdates(data, filterData, postID, r.URL.Path == "/")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/user/login", http.StatusSeeOther)
				return
			} else {

				h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}

		h.render(w, http.StatusOK, "index.html", data)
	}
}
