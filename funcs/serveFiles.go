package funcs

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	if checkSession(w, r) {
		http.Redirect(w, r, "/mainPage", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "./html/login.html")
}
func VerifyEmailPage(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "./html/verifyEmailpage.html")
}

// Renders registration page if user is not logged in
func ServeRegister(w http.ResponseWriter, r *http.Request) {
	if checkSession(w, r) {
		http.Redirect(w, r, "/mainPage", http.StatusFound)
		return
	}
	http.ServeFile(w, r, "./html/register.html")
}
func ServeAdminPage(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "projectGo")
	userID, ok := session.Values["userID"].(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем пользователя из базы данных
	objectID, _ := primitive.ObjectIDFromHex(userID)
	var user User
	err := UserCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil || user.Status != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	http.ServeFile(w, r, "./html/adminPage.html")
}

// Renders the main page if the user is logged in
func ServeMain(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "projectGo")
	userID, ok := session.Values["userID"].(string)
	if !ok || userID == "" {
		http.Redirect(w, r, "/loginPage", http.StatusFound)
		return
	}

	// Retrieve user data from the database using userID
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	var user User
	err = UserCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Inject email into the page
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "./html/main.html")

	// Add script for email
	emailScript := fmt.Sprintf("<script>const userEmail = '%s';</script>", user.Email)
	w.Write([]byte(emailScript))
}
func ServeRecipesPage(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "projectGo")
	_, ok := session.Values["userID"]
	if !ok {
		http.Redirect(w, r, "/loginPage", http.StatusFound)
		return
	}

	http.ServeFile(w, r, "./html/recipes.html")
}
