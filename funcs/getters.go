package funcs

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

// Handler to get user avatar
func GetUserAvatar(w http.ResponseWriter, r *http.Request) {
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

	// Возвращаем email в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"avatar": user.AvatarPath,
	})
}
func GetUserEmailHandler(w http.ResponseWriter, r *http.Request) {
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

	// Возвращаем email в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"email": user.Email,
	})
}
func GetUserStatusHandler(w http.ResponseWriter, r *http.Request) {
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
	// Возвращаем email в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": user.Status,
	})
}
func GetUserNameHandler(w http.ResponseWriter, r *http.Request) {
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

	// Возвращаем email в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"name": user.Name,
	})
}
func GetUserRegMetHandler(w http.ResponseWriter, r *http.Request) {
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

	// Возвращаем email в формате JSON
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"reg_met": user.RegMet,
	})
}
