package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TwitchBot struct {
	UserId          string `json:"user_id,omitempty"`
	UserName        string `json:"user_name,omitempty"`
	AppAccessToken  AppAccessToken
	UserAccessToken UserAccessToken
	Configured      bool
	Secret          string
}

type TwitchUserInfoResponse struct {
	Data []User `json:"data"`
}

func (bot TwitchBot) printUserAccessTokenURL() {

	scopes := []string{"user:bot", "user:read:chat", "user:write:chat"}
	queryParams := url.Values{
		"response_type": {"code"},
		"client_id":     {twitchClientID},
		"redirect_uri":  {fmt.Sprintf("%s/botoauth", twitchRedirectURL)},
		"scope":         {strings.Join(scopes, " ")},
	}

	url := "https://id.twitch.tv/oauth2/authorize?" + queryParams.Encode()

	fmt.Printf("\033[32mIMPORTANT!\033[0m. You must use this URL to fully authenticate the Twitch Bot.\nPaste into a browser that is signed into the %s account\n\n\033[32m%s\033[0m\n\n", os.Getenv("TWITCH_USERNAME"), url)

}

func NewTwitchBot() *TwitchBot {

	var bot = new(TwitchBot)

	bot.Secret = "random_string_for_now"

	bot.printUserAccessTokenURL()

	scopes := []string{"user:bot", "user:read:chat", "user:write:chat"}
	queryParams := url.Values{
		"client_id":     {twitchClientID},
		"client_secret": {twitchClientSecret},
		"grant_type":    {"client_credentials"},
		"scope":         {strings.Join(scopes, " ")},
	}

	fullURL := "https://id.twitch.tv/oauth2/token?" + queryParams.Encode()

	resp, err := http.Post(fullURL, "application/json", http.NoBody)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &bot.AppAccessToken)

	user := bot.getUserInfo(os.Getenv("TWITCH_USERNAME"))

	bot.UserId = user.ID
	bot.UserName = user.DisplayName

	return bot
}

func (bot *TwitchBot) getUserInfo(login string) User {
	queryParams := url.Values{
		"login": {login},
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users?"+queryParams.Encode(), http.NoBody)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bot.AppAccessToken.Token))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Client-Id", os.Getenv("TWITCH_CLIENT_ID"))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var userInfoResponse = new(TwitchUserInfoResponse)

	json.Unmarshal(body, userInfoResponse)

	fmt.Printf("Twitch User Info: %+v\n\n", *userInfoResponse)

	return userInfoResponse.Data[0]
}

func (bot TwitchBot) NewUser(authorizationCode string) (*User, error) {

	var err error
	token := getUserAccessToken(authorizationCode, "oauth")
	fmt.Printf("New User Access Token Struct Created: %+v\n\n", token)

	login, err := token.ValidateToken()

	if err != nil {
		log.Println(err)
	}

	user := bot.getUserInfo(login)

	user.UserAccessToken = token

	return &user, err

}

func (bot *TwitchBot) SubscribeToTwitchChat(broadcasterUserID string) error {

	url := "https://api.twitch.tv/helix/eventsub/subscriptions"

	payload := map[string]any{
		"type":    "channel.chat.message",
		"version": "1",
		"condition": map[string]string{
			"broadcaster_user_id": broadcasterUserID,
			"user_id":             bot.UserId,
		},
		"transport": map[string]string{
			"method":   "webhook",
			"callback": fmt.Sprintf("%s/event", twitchCallbackURL),
			"secret":   bot.Secret,
		},
	}

	bodyBytes, err := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))

	req.Header.Set("Authorization", "Bearer "+bot.AppAccessToken.Token)
	req.Header.Set("Client-Id", twitchClientID)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(body))

	return err
}

func (bot TwitchBot) VerifyEventSignature(messageToCheck string, signature string) bool {

	h := hmac.New(sha256.New, []byte(bot.Secret))

	h.Write([]byte(messageToCheck))

	localSignature := h.Sum(nil)

	_, signature, _ = strings.Cut(signature, "=")

	receivedSignature, err := hex.DecodeString(signature)

	if err != nil {
		log.Println("Unable to decode received signature")
		return false
	}

	if hmac.Equal(localSignature, receivedSignature) {
		return true
	}

	return false
}
