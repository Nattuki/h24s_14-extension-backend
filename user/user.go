package user

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Me struct {
	Name       string `json:"name"`
	IconBase64 string `json:"iconBase64"`
	IsVisitor  bool   `json:"isVisitor"`
}

func GetUserInformation(accessToken string) Me {
	req, err := http.NewRequest(http.MethodGet, "https://q.trap.jp/api/v3/users/me", nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	var myInformation Me
	if resp.StatusCode == 200 {
		myInformation.IsVisitor = false
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		json.Unmarshal(body, &myInformation)
		myInformation.IconBase64 = GetIcon(accessToken)
	} else {
		myInformation.IsVisitor = true
	}

	return myInformation
}

func GetIcon(accessToken string) string {
	req, err := http.NewRequest(http.MethodGet, "https://q.trap.jp/api/v3/users/me/icon", nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	base64Data := base64.StdEncoding.EncodeToString(body)
	return base64Data
}
