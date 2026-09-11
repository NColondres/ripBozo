package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type AppAccessToken struct {
	Token     string   `json:"access_token"`
	ExpiresIn int      `json:"expires_in"`
	TokenType string   `json:"token_type"`
	Scope     []string `json:"scope"`
}

type UserAccessToken struct {
	Token        string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	TokenType    string   `json:"token_type"`
}

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	ProfileImageUrl string `json:"profile_image_url"`
	CreatedAt       string `json:"created_at"`
	Description     string `json:"description,omitempty"`
	UserAccessToken UserAccessToken
}

func (token *UserAccessToken) ValidateToken() (string, error) {
	var err error
	login := ""
	client := &http.Client{}
	url := "https://id.twitch.tv/oauth2/validate"

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		log.Println(err)
		return login, err
	}

	req.Header.Set("Authorization", "OAuth "+token.Token)

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return login, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var jsonLogin struct {
		Login string `json:"login"`
	}

	err = json.Unmarshal(body, &jsonLogin)

	if err != nil {
		fmt.Println(err)
	}

	login = jsonLogin.Login

	return login, err
}

func getUserAccessToken(authorizationCode string, redirect_uri_path string) UserAccessToken {
	queryParams := url.Values{
		"client_id":     {twitchClientID},
		"client_secret": {twitchClientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {authorizationCode},
		"redirect_uri":  {fmt.Sprintf("%s/%s", twitchRedirectURL, redirect_uri_path)},
	}

	url := "https://id.twitch.tv/oauth2/token?" + queryParams.Encode()

	resp, err := http.Post(url, "application/x-www-form-urlencoded", http.NoBody)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))
	var token = new(UserAccessToken)

	json.Unmarshal(body, token)

	return *token
}
