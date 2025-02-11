package funcs

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
)

func Login(w http.ResponseWriter, r *http.Request) {
	LoadLogger()
	var reqUser User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var dbUser User
	err := UserCollection.FindOne(context.TODO(), bson.M{"email": reqUser.Email}).Decode(&dbUser)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"status": "fail", "message": "User not found"})
		return
	}

	if !checkPasswordHash(reqUser.Password, dbUser.Password) {
		json.NewEncoder(w).Encode(map[string]string{"status": "fail", "message": "Invalid password"})
		return
	}

	// Store session information
	session, _ := Store.Get(r, "projectGo")
	session.Values["userID"] = dbUser.ID.Hex()
	session.Save(r, w)
	err = performAction("login", dbUser.ID, "success")
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Login successful", "session_key": os.Getenv("SESSION_KEY")})
}

// Checks the Login status by looking for the user session
func CheckLoginStatus(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "projectGo")
	userID, ok := session.Values["userID"].(string)
	if ok && userID != "" {
		json.NewEncoder(w).Encode(map[string]string{"status": "logged_in"})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "not_logged_in"})
	}
}
func checkPasswordHash(password, hash string) bool {
	// bcrypt.CompareHashAndPassword возвращает ошибку, если пароли не совпадают
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Checks if the user is logged in based on the session
func checkSession(w http.ResponseWriter, r *http.Request) bool {
	session, _ := Store.Get(r, "projectGo")
	userID, ok := session.Values["userID"].(string)
	return ok && userID != ""
}

// Logs out the user by clearing their session
func Logout(w http.ResponseWriter, r *http.Request) {
	LoadLogger()

	session, _ := Store.Get(r, "projectGo")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}
	var user User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err := UserCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	performAction("Logout", user.ID, "success")
	session.Values = nil
	session.Options.MaxAge = -1
	session.Save(r, w)

	// Respond with success message
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Logged out"})
}
