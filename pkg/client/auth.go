package client

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/fi3te/sptfio/pkg/config"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2"

	"github.com/zmb3/spotify/v2"
)

const alphanumericCharacters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var authenticator *spotifyauth.Authenticator
var scChan chan *spotify.Client
var state string
var codeVerifier string

func Login(ctx context.Context, cfg *config.Config) (*spotify.Client, error) {
	authenticator = spotifyauth.New(
		spotifyauth.WithClientID(cfg.ClientID),
		spotifyauth.WithRedirectURL(cfg.RedirectURL),
		spotifyauth.WithScopes(spotifyauth.ScopePlaylistReadPrivate, spotifyauth.ScopePlaylistModifyPrivate))
	scChan = make(chan *spotify.Client)
	state = generateRandomAlphanumericString(20)
	codeVerifier = generateVerifier()

	addr := fmt.Sprintf(":%d", cfg.GetPort())
	server := &http.Server{Addr: addr}
	log.Printf("Starting server on address %s...", addr)
	http.HandleFunc(cfg.GetCallback(), completeAuth)
	go server.ListenAndServe()

	codeChallenge := generateCodeChallenge(codeVerifier)
	url := authenticator.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)
	log.Printf("Please visit the following url:\n%s\n", url)

	sc := <-scChan

	log.Println("Shutting down server...")
	serverErr := server.Shutdown(ctx)
	if serverErr != nil {
		log.Printf("failed shutting down server: %v\n", serverErr)
	}

	var err error
	if sc == nil {
		err = errors.New("login unsuccessful")
	}

	return sc, err
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	defer close(scChan)
	token, err := authenticator.Token(r.Context(), state, r,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		message := fmt.Sprintf("couldn't obtain token %v", err)
		http.Error(w, message, http.StatusForbidden)
		return
	}
	if usedState := r.FormValue("state"); usedState != state {
		message := fmt.Sprintf("invalid state: %s != %s", state, usedState)
		http.Error(w, message, http.StatusForbidden)
		return
	}
	sc := spotify.New(authenticator.Client(r.Context(), token))
	fmt.Fprintf(w, "login successful")
	scChan <- sc
}

func generateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	return strings.ReplaceAll(base64.URLEncoding.EncodeToString(hash[:]), "=", "")
}

func generateVerifier() string {
	return generateRandomAlphanumericString(128)
}

func generateRandomAlphanumericString(length int) string {
	var text string
	for a := 0; a < length; a++ {
		i := rand.Int() % len(alphanumericCharacters)
		text += string(alphanumericCharacters[i])
	}
	return text
}
