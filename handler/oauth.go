package handler

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"h24s_14-extension-backend/util"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

const (
	lenState        = 30
	lenCodeVerifier = 64
	AuthorizeURL    = "https://q.trap.jp/api/v3/oauth2/authorize"
	GetTokenURL     = "https://q.trap.jp/api/v3/oauth2/token"
)

var (
	ClientID string
)

func init() {
	ClientID = os.Getenv("TRAQ_CLIENT_ID")
}

func HandleGetOAuthUrl(c echo.Context) error {
	sess, err := session.Get("note_session", c)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to get the session.")
	}

	state, _ := util.RandomString(lenState)
	codeVerifier, _ := util.RandomString(lenCodeVerifier)

	sess.Values["state"] = state
	sess.Values["code_verifier"] = codeVerifier
	sess.Save(c.Request(), c.Response())

	b := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b[:])

	u, _ := url.Parse(AuthorizeURL)
	q := u.Query()
	q.Add("response_type", "code")
	q.Add("client_id", ClientID)
	q.Add("state", state)
	q.Add("code_challenge", codeChallenge)
	q.Add("code_challenge_method", "S256")
	u.RawQuery = q.Encode()

	return c.String(http.StatusOK, u.String())
}

func HandleGetToken(c echo.Context) error {
	sess, err := session.Get("note_session", c)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Failed to get the session.")
	}
	u := c.Request().URL
	if u.Query().Get("state") != sess.Values["state"] {
		return c.String(http.StatusInternalServerError, "Invalid state.")
	}

	codeVerifier := sess.Values["code_verifier"].(string)
	sess.Values["state"] = ""
	sess.Values["code_verifier"] = ""
	sess.Save(c.Request(), c.Response())

	code := u.Query().Get("code")

	q := url.Values{}
	q.Add("grant_type", "authorization_code")
	q.Add("client_id", ClientID)
	q.Add("code", code)
	q.Add("code_verifier", codeVerifier)
	req, err := http.NewRequest(http.MethodPost, GetTokenURL, strings.NewReader(q.Encode()))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to write into the new request.")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get the response from the authorization server.")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var token TokenResponse
	json.Unmarshal(body, &token)
	sess.Values["access_token"] = token.AccessToken
	sess.Values["token_type"] = token.TokenType
	sess.Values["expires_in"] = token.ExpiresIn
	sess.Values["refresh_token"] = token.RefreshToken
	sess.Values["scope"] = token.Scope
	sess.Values["id_token"] = token.IDToken
	sess.Save(c.Request(), c.Response())

	return c.String(http.StatusAccepted, "success")
}
