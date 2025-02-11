package funcs

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

type App struct {
	Config *oauth2.Config
}

func (a *App) OAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := a.Config.AuthCodeURL("hello world", oauth2.AccessTypeOffline)
	Logger.Info("Generated OAuth URL: ", url) // Added this line for logging
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Registers a new user by saving their email and password in the database
func (a *App) OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	LoadLogger()
	code := r.URL.Query().Get("code")

	// Exchanging the code for an access token
	t, err := a.Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Creating an HTTP client to make authenticated request using the access key.
	// This client method also regenerate the access key using the refresh key.
	client := a.Config.Client(context.Background(), t)

	// Getting the user public details from google API endpoint
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Closing the request body when this function returns.
	// This is a good practice to avoid memory leak
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Извлекаем нужные значения
	email, emailOk := userInfo["email"].(string)
	name, nameOk := userInfo["name"].(string)
	if !emailOk || !nameOk {
		http.Error(w, "Invalid response from Google API", http.StatusInternalServerError)
		return
	}
	var existingUser User
	err = UserCollection.FindOne(context.TODO(), bson.M{"email": email, "reg_met": "google"}).Decode(&existingUser)

	if err == nil {
		// Email уже существует
		session, _ := Store.Get(r, "projectGo")
		session.Values["userID"] = existingUser.ID.Hex()
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, "Failed to save session", http.StatusInternalServerError)
			return
		}

		err = performAction("login", existingUser.ID, "success")
		if err != nil {
			fmt.Println(err)
			return
		}
		http.Redirect(w, r, "/mainPage", http.StatusTemporaryRedirect)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"redirect": "/mainPage",
		})
		return
	}
	user := User{Email: email, Name: name, Status: "user", RegMet: "google"}
	_, err = UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}
	err = performAction("google login", user.ID, "success")
	if err != nil {
		fmt.Println(err)
		return
	}
	session, _ := Store.Get(r, "projectGo")
	session.Values["userID"] = user.ID.Hex()
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	fmt.Println("Redirect URL:", os.Getenv("REDIRECT_URL"))
	err = performAction("login", user.ID, "success")
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Redirect(w, r, "/mainPage", http.StatusTemporaryRedirect)

	json.NewEncoder(w).Encode(map[string]string{
		"redirect": "/mainPage",
	})
	return

}
