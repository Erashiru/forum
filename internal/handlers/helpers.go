package handlers

import (
	"bytes"
	"fmt"
	"forum/models"
	"net/http"
)

func (s *handlers) render(w http.ResponseWriter, status int, page string, data *models.TemplateData) {
	t, ok := s.app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		s.app.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		s.app.ServerError(w, err)
		return
	}
	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (s *handlers) newTemplateData() *models.TemplateData {
	return &models.TemplateData{}
}

func (s *handlers) GetUserIDBySession(r *http.Request) int {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		return 0
	}

	return s.service.GetUserIDBySessionToken(sessionCookie.Value)
}

func (h *handlers) ErrorHandler(w http.ResponseWriter, r *http.Request, errorCode int, msg string) {
	t, ok := h.app.TemplateCache["error.html"]

	if !ok {
		fmt.Println("the template error.html does not exist")
		return
	}
	errorki := models.CustomError{
		ErrorCode: errorCode,
		ErrorMsg:  msg,
	}

	w.WriteHeader(errorCode)
	t.Execute(w, errorki)
}
