package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var (
	twitchClientID     string
	twitchClientSecret string
	twitchRedirectURL  string
	twitchCallbackURL  string
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	twitchClientID = os.Getenv("TWITCH_CLIENT_ID")
	twitchClientSecret = os.Getenv("TWITCH_CLIENT_SECRET")
	twitchRedirectURL = os.Getenv("TWITCH_REDIRECT_URL")
	twitchCallbackURL = os.Getenv("TWITCH_CALLBACK_URL")
}

func main() {

	twitchBot := NewTwitchBot()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve the index.html file for the root URL "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "static/index.html")
	})

	// HealthCheck
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := struct {
			Status string
		}{
			Status: "Healthy",
		}

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("error encoding response struct: %v", err)
		}

	})

	// TODO:
	// Fix the /.*oauth endpoints to pass the scopes into the getUserAccessToken
	// so that scopes stay consitent throughout the whole oauth process.

	// This is only for the twitchbot account and gettin the UserAccessToken at startup
	http.HandleFunc("/botoauth", func(w http.ResponseWriter, r *http.Request) {
		urlParameters := r.URL.Query()
		for k, v := range urlParameters {
			fmt.Fprintf(w, "Param %s => %v\n", k, v)
			if k == "code" {

				twitchBot.UserAccessToken = getUserAccessToken(v[0], "botoauth")

				if twitchBot.UserAccessToken.Token != "" {
					twitchBot.Configured = true
				} else {
					log.Fatalln("Unsuccessful in intializing TwitchBot")
				}
			}

		}
	})

	http.HandleFunc("/oauth", func(w http.ResponseWriter, r *http.Request) {

		if twitchBot.Configured == false {

			error_message := "TwitchBot has not been initalized yet. Cannot complete authentication\n\n"

			log.Print(error_message)
			fmt.Fprint(w, error_message)

		} else {
			urlParameters := r.URL.Query()
			for k, v := range urlParameters {
				fmt.Fprintf(w, "Param %s => %v\n", k, v)
				if k == "code" {

					user, _ := twitchBot.NewUser(v[0])
					fmt.Printf("New User Struct Created: %+v\n\n", *user)

					twitchBot.SubscribeToTwitchChat(user.ID)
				}

			}
		}

	})

	http.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {

		// TODO: verify each twitch event message
		// How to: https://dev.twitch.tv/docs/eventsub/handling-webhook-events/

		twitchMessageID := r.Header.Get("Twitch-Eventsub-Message-Id")
		twitchMessageTimestamp := r.Header.Get("Twitch-Eventsub-Message-Timestamp")
		twitchMessageSignature := r.Header.Get("Twitch-Eventsub-Message-Signature")

		// If we see all three of these headers, we know its a twitch chat message
		if twitchMessageID != "" && twitchMessageTimestamp != "" && twitchMessageSignature != "" {

			defer r.Body.Close()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading body: %v\n", err)
			}
			// Now we verify the signature to make sure this is actually coming from twitch.

			if !twitchBot.VerifyEventSignature(twitchMessageID+twitchMessageTimestamp+string(body), twitchMessageSignature) {
				log.Println("Twitch Verfication Failed")
				return
			}

			fmt.Println(string(body))

			if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {

				var reqChallenge struct {
					Challenge string `json:"challenge"`
				}

				if err := json.Unmarshal(body, &reqChallenge); err != nil {
					log.Print("Unable to unmarshall the challenge value")
					return
				} else {

					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(http.StatusOK)
					io.WriteString(w, reqChallenge.Challenge)
					return

				}

			}

		} else {
			io.WriteString(w, "/event is reachable")
			return
		}

	})

	// Start the web server on port 8000
	log.Printf("Listening on \033[32m http://localhost:8080\033[0m\n\n")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
