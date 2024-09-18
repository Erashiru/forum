package handlers

import (
	"errors"
	"fmt"
	"forum/models"
	"net/http"
	"strconv"
)

func (h *handlers) postView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method == http.MethodGet {
		data := h.newTemplateData()
		data, err = h.service.GetPostData(data, id)

		if err != nil {
			if errors.Is(err, models.NotFound) {
				h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return

			} else {
				h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		data.Form = models.CommentForm{}
		data.UserID = h.GetUserIDBySession(r)
		data, err = h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		h.render(w, http.StatusOK, "viewpost.html", data)
	} else if r.Method == http.MethodPost {
		data := h.newTemplateData()
		data.UserID = h.GetUserIDBySession(r)
		if data.UserID == 0 {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		data, err = h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		var commentID int

		r.ParseForm()

		if r.Form.Has("comment_reaction") {
			commentID, err = strconv.Atoi(r.FormValue("commentID"))

			if err != nil {
				h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}

		}

		commentData := &models.UserReaction{

			Reaction:        r.FormValue("reaction"),
			CommentReaction: r.FormValue("comment_reaction"),
			Comment:         r.FormValue("comment"),
			CommentID:       commentID,
		}

		data, err := h.service.PostUpdate(data, commentData, id)
		if err != nil {
			if errors.Is(err, models.ErrNotValidPostForm) {
				h.render(w, http.StatusBadRequest, "viewpost.html", data)
				return
			} else if errors.Is(err, models.NotFound) {
				h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				return
			} else {
				h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}
		http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
		return
	}
}

func (h *handlers) createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if r.URL.Path != "/post/create" {
			h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}
		data := h.newTemplateData()
		data.Form = models.PostForm{}
		data.UserID = h.GetUserIDBySession(r)
		data, err := h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		h.render(w, http.StatusOK, "createpost.html", data)

	} else if r.Method == http.MethodPost {
		if r.URL.Path != "/post/create" {
			h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
			return
		}

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
			return
		}

		Form := &models.PostForm{
			Title:    r.FormValue("title"),
			Content:  r.FormValue("content"),
			Category: r.Form["category"],
		}

		data, id, err := h.service.PostCreate(data, Form)
		if err != nil {
			if err == models.ErrNotValidPostForm {
				h.render(w, http.StatusBadRequest, "createpost.html", data)
				return
			} else {
				h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}

		http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
	} else {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
