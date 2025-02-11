package funcs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
)

// UpdatePassword handles the password update request
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "projectGo")
	userID, ok := session.Values["userID"].(string)
	if !ok || userID == "" {
		http.Redirect(w, r, "/loginPage", http.StatusFound)
		return
	}

	// Decode the request body
	var passwordData struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	err := json.NewDecoder(r.Body).Decode(&passwordData)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Fetch the user from the database
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	var user User
	err = UserCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching user", http.StatusInternalServerError)
		}
		return
	}

	// Compare the old password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordData.OldPassword))
	if err != nil {
		http.Error(w, "Old password is incorrect", http.StatusUnauthorized)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Update the user's password in the database
	_, err = UserCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"password": string(hashedPassword)}},
	)
	if err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	// Respond with success
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Password updated successfully",
	})
}

// UserPage handles the user's page request.
func UserPage(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	session, _ := Store.Get(r, "projectGo")
	userID, ok := session.Values["userID"].(string)
	if !ok || userID == "" {
		http.Redirect(w, r, "/loginPage", http.StatusFound)
		return
	}

	// Fetch user from database
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

	// Serve user page HTML
	http.ServeFile(w, r, "./html/userPage.html")
}

// UpdateName handles the update of the user's name.
func UpdateName(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "projectGo")
	userID, ok := session.Values["userID"].(string)
	if !ok || userID == "" {
		http.Redirect(w, r, "/loginPage", http.StatusFound)
		return
	}

	var newName string
	if err := json.NewDecoder(r.Body).Decode(&newName); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Retrieve user and update name
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

	user.Name = newName

	// Update user in database
	_, err = UserCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": user})
	if err != nil {
		http.Error(w, "Error while updating name", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Name changed successfully"})
}

// UpdateAvatarHandler handles the update of the user's avatar.
func UpdateAvatarHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "projectGo")
	userID, ok := session.Values["userID"].(string)
	if !ok || userID == "" {
		http.Redirect(w, r, "/loginPage", http.StatusFound)
		return
	}

	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Handle avatar upload
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		http.Error(w, `{"status": "error", "message": "Unable to parse form"}`, http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Error retrieving file"}`, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure the avatars directory exists
	avatarDir := "./avatars"
	if _, err := os.Stat(avatarDir); os.IsNotExist(err) {
		if mkdirErr := os.Mkdir(avatarDir, os.ModePerm); mkdirErr != nil {
			http.Error(w, `{"status": "error", "message": "Server error: unable to create directory"}`, http.StatusInternalServerError)
			return
		}
	}

	// Generate unique file name
	uniqueFileName := fmt.Sprintf("%s-%s", uuid.New().String(), handler.Filename)
	filePath := fmt.Sprintf("%s/%s", avatarDir, uniqueFileName)

	// Save file
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Unable to save file"}`, http.StatusInternalServerError)
		return
	}
	fmt.Println("Saving file to:", filePath)
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		fmt.Println("Error copying file:", err)
		http.Error(w, `{"status": "error", "message": "Error saving file"}`, http.StatusInternalServerError)
		return
	}
	fmt.Printf("File received: %s, size: %d bytes\n", handler.Filename, handler.Size)
	// Update user avatar in the database
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Invalid user ID"}`, http.StatusInternalServerError)
		return
	}

	_, err = UserCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"avatarPath": filePath}})
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Error updating avatar"}`, http.StatusInternalServerError)
		return
	}

	// Respond with the new avatar URL
	response := map[string]interface{}{
		"status":    "success",
		"message":   "Avatar updated successfully",
		"avatarUrl": fmt.Sprintf("./avatars/%s", uniqueFileName),
	}

	// Encode response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"status": "error", "message": "Error encoding JSON"}`, http.StatusInternalServerError)
	}
}
