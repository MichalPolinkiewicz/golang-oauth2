package main

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "2325387494-6pjfob3cu88dh726dm0075ae480k27j3.apps.googleusercontent.com",
		ClientSecret: "Jer6VVpHDfEI5mZPyoCyljE5",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	state = "random"
)

func main() {

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil)

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html>
						<body>
							<a href="/login">Google Log In</a>
						</body>
					</html>`
	fmt.Fprintf(w, htmlIndex)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	//tworzony jest url do serwera autoryzujacego wraz z state
	url := googleOauthConfig.AuthCodeURL(state)

	//przekierowanie pod auth url
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	//po zalogowaniu serwer autoryzujacy przekierowuje do aplikacji, sprawdzany jest state
	//jesli jest inny niz przekazany to błąd
	//checking state - if different than given return error
	if r.FormValue("state") != state {
		fmt.Println("state is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//pobieramy kod z requestu i wymieniamy go na token
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Printf("could not get token: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//pobieramy dane usera z serwera autoryzujacego
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("could not create get request: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Println(resp)

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("could not parse response: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "response: %s", content)
}
