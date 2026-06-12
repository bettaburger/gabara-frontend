package pkce
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/google/uuid"
	//"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

// redirectURI is the OAuth redirect URI for the application.
const redirectURI = "http://127.0.0.1:8080/callback"

var (
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate))
	ch = make(chan *spotify.Client)
	state string
	codeVerifier string
	codeChallenge string
)

func NewAuth() {
	//r := mux.NewRouter()
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	// check .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	//clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	// generate random state seed, verifier and challenge 
	state = uuid.New().String() 
	codeVerifier = oauth2.GenerateVerifier()
	codeChallenge = oauth2.S256ChallengeFromVerifier(codeVerifier)

	auth = spotifyauth.New(
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithRedirectURL(redirectURI))
	
	// the url to OAuth2.0 that asks for permission for the required parameters
	url := auth.AuthURL(state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("client_id", clientID),
	)

	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	// wait for auth to complete
	client := <-ch
	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CALLBACK HIT")
	token, err := auth.Token(r.Context(), state, r,
	oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatal("State mismatch: %s != %s\n", st, state)
	}
	client := spotify.New(auth.Client(r.Context(), token))
	fmt.Printf("token: %v", token)
	fmt.Fprintf(w, "login completed")
	ch <- client
}