package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/models"
	"forum/pkg"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	ACCESS_TOKEN         = "access_token"
	GOOGLE_USER_INFO_URL = "https://www.googleapis.com/oauth2/v2/userinfo"
	GOOGLE_TOKEN_URL     = "https://accounts.google.com/o/oauth2/token"
	GITHUB_USER_INFO_URL = "https://api.github.com/user"
	GITHUB_TOKEN_URL     = "https://github.com/login/oauth/access_token"
)

func (h *handlers) Load(w http.ResponseWriter, r *http.Request, form *models.UserSignupForm) {
	id, err := h.service.SaveUser(form)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, models.ErrDuplicateEmail) || errors.Is(err, models.ErrDuplicateName) {
			user, err := h.service.GetUserByEmail(form.Email)
			if err != nil {
				h.app.ServerError(w, err)
				return
			}

			id = user.ID
		} else {
			h.app.ServerError(w, err)
			return
		}
	}

	err = h.service.CreateSession(w, r, id)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// GOOGLE

func (h *handlers) googlelogin(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(h.conf.ExternalAuth.GoogleClientID, h.conf.ExternalAuth.GoogleRedirectURL)
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile",
		h.conf.ExternalAuth.GoogleClientID, h.conf.ExternalAuth.GoogleRedirectURL)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *handlers) googleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	clientData := fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code, h.conf.ExternalAuth.GoogleClientID, h.conf.ExternalAuth.GoogleClientSecret, h.conf.ExternalAuth.GoogleRedirectURL)

	resp, err := http.Post(GOOGLE_TOKEN_URL, "application/x-www-form-urlencoded", strings.NewReader(clientData))
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		h.app.ServerError(w, err)
		return
	}

	accessToken, ok := tokenResponse[ACCESS_TOKEN].(string)
	if !ok {
		h.ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	newReq, err := http.NewRequest("GET", GOOGLE_USER_INFO_URL, nil)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	newReq.Header.Add("Authorization", "Bearer "+accessToken)

	userInfoResp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	defer userInfoResp.Body.Close()

	var googleInfo struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&googleInfo); err != nil {
		h.app.ServerError(w, err)
		return
	}

	form := models.UserSignupForm{
		Name:  googleInfo.Name,
		Email: googleInfo.Email,
	}
	form.Password, err = pkg.GenerateRandomPassword(8)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	h.Load(w, r, &form)
}

// GITHUB

func (h *handlers) githublogin(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		h.conf.ExternalAuth.GithubClientID, h.conf.ExternalAuth.GithubRedirectURL,
	)

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func (h *handlers) githubCallback(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	payload := url.Values{
		"client_id":     {h.conf.ExternalAuth.GithubClientID},
		"client_secret": {h.conf.ExternalAuth.GithubClientSecret},
		"code":          {code},
		"redirect_uri":  {h.conf.ExternalAuth.GithubRedirectURL},
	}
	resp, err := http.PostForm(GITHUB_TOKEN_URL, payload)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	tokenResp, err := url.ParseQuery(string(body))
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	accessToken := tokenResp.Get(ACCESS_TOKEN)

	newreq, err := http.NewRequest("GET", GITHUB_USER_INFO_URL, nil)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	newreq.Header.Set("Authorization", "token "+accessToken)
	client := &http.Client{}
	resp, err = client.Do(newreq)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}
	defer resp.Body.Close()

	var githubInfo struct {
		Login string `json:"login"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubInfo); err != nil {
		h.app.ServerError(w, err)
		return
	}

	if githubInfo.Email == "" {
		githubInfo.Email = githubInfo.Login + "@github.com"
	}
	form := models.UserSignupForm{
		Name:  githubInfo.Login,
		Email: githubInfo.Email,
	}
	form.Password, err = pkg.GenerateRandomPassword(8)
	if err != nil {
		h.app.ServerError(w, err)
		return
	}

	h.Load(w, req, &form)
}
