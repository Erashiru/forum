package handlers

import (
	"fmt"
	"forum/models"
	"forum/pkg"
	"net/http"
	"strconv"
)

func (h *handlers) userPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/posts" {
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
		data, err = h.service.UserPosts(data)
		if err != nil {
			fmt.Println(err)
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		h.render(w, http.StatusOK, "userposts.html", data)
	} else if r.Method == http.MethodPost {
		data := h.newTemplateData()
		data.UserID = h.GetUserIDBySession(r)
		data, err := h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		postID, err := strconv.Atoi(r.FormValue("postID"))
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		data, err = h.service.UserPostsUpdate(data, postID, r.FormValue("reaction"))
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		h.render(w, http.StatusOK, "userposts.html", data)
	} else {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *handlers) likedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/posts/liked" {
		h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method == http.MethodGet {
		data := h.newTemplateData()
		data.UserID = h.GetUserIDBySession(r)
		data, err := h.service.GetUserData(data)
		if err != nil {
			fmt.Println(err)
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		data, err = h.service.LikedPostsLoader(data)
		if err != nil {
			fmt.Println(err, "aboba")
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		h.render(w, http.StatusOK, "likedposts.html", data)
	} else if r.Method == http.MethodPost {

		postID, err := strconv.Atoi(r.FormValue("postID"))
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		data := h.newTemplateData()
		data.UserID = h.GetUserIDBySession(r)
		data, err = h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		data, err = h.service.LikedPostsUpdate(data, postID, r.FormValue("reaction"))
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		h.render(w, http.StatusOK, "likedposts.html", data)
	} else {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *handlers) dislikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/user/posts/disliked" {
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
		data, err = h.service.DislikedPostsLoader(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		h.render(w, http.StatusOK, "disliked.html", data)
	} else if r.Method == http.MethodPost {

		postID, err := strconv.Atoi(r.FormValue("postID"))
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		data := h.newTemplateData()
		data.UserID = h.GetUserIDBySession(r)
		data, err = h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		data, err = h.service.DislikedPostsUpdate(data, postID, r.FormValue("reaction"))
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		h.render(w, http.StatusOK, "disliked.html", data)
	} else {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *handlers) updateUserPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		data := h.newTemplateData()
		// data, err := h.service.GetPostData(data, id)
		// if err != nil {
		// 	if errors.Is(err, models.NotFound) {
		// 		h.ErrorHandler(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		// 		return

		// 	} else {
		// 		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		// 		return
		// 	}
		// }
		data.Form = models.CommentForm{}
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
			Content: r.FormValue("content"),
			Title:   r.FormValue("title"),
		}
		data.PostID = pkg.SplitID(r.URL.Path)
		post, err := h.service.GetPostData(data, data.PostID)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		if post.Username != data.Username {
			h.ErrorHandler(w, r, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		err = h.service.UpdatePost(data, Form)
		if err != nil {
			if err == models.ErrNotValidPostForm {
				h.render(w, http.StatusBadRequest, "updatepost.html", data)
				return
			} else {
				h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}

		http.Redirect(w, r, fmt.Sprintf("/post/view?id=%d", data.PostID), http.StatusSeeOther)
	} else if r.Method == http.MethodGet {
		data := h.newTemplateData()
		data.Form = models.PostForm{}
		data.UserID = h.GetUserIDBySession(r)
		data.PostID = pkg.SplitID(r.URL.Path)
		data, err := h.service.GetUserData(data)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		post, err := h.service.GetPostData(data, data.PostID)
		if err != nil {
			h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		if post.Username != data.Username {
			h.ErrorHandler(w, r, http.StatusForbidden, http.StatusText(http.StatusForbidden))
			return
		}
		h.render(w, http.StatusOK, "updatepost.html", data)

	} else {
		h.ErrorHandler(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

func (h *handlers) deleteUserPost(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData()
	data.Form = models.PostForm{}
	data.UserID = h.GetUserIDBySession(r)
	data.PostID = pkg.SplitID(r.URL.Path)
	data, err := h.service.GetUserData(data)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = h.service.DeletePost(data.PostID)
	if err != nil {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
