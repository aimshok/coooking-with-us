package funcs

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"os"
	"strconv"
)

func DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ID string `json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Delete from database
	result, err := RecipeCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, "Error deleting from database", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe deleted successfully"})
}
func AddRecipeHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" { // Handle CORS preflight request
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	Logger.Info("addRecipe Handler ")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Parse form data
	err := r.ParseMultipartForm(0 << 20) // Max 10MB file size
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Get fields
	name := r.FormValue("name")
	caloriesStr := r.FormValue("calories")
	level := r.FormValue("level")
	description := r.FormValue("description")
	priceStr := r.FormValue("price") // Get price as string
	file, header, err := r.FormFile("image")

	fmt.Println("Received Form Data:")
	fmt.Println("Name:", name)
	fmt.Println("Calories:", caloriesStr)
	fmt.Println("Level:", level)
	fmt.Println("Description:", description)
	fmt.Println("Price:", priceStr)
	fmt.Println("File Name:", header.Filename)

	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}
	calories, err := strconv.ParseFloat(caloriesStr, 64)
	// Save the file (for example, to a "uploads" folder)
	filePath := "Recipes/" + header.Filename
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Save to database
	Recipe := Recipe{
		Name:     name,
		Desc:     description,
		Calories: calories,
		Price:    price,
		Level:    level,
		Path:     header.Filename,
	}
	Logger.Info("Inserting recipe into database...")
	_, err = RecipeCollection.InsertOne(context.TODO(), Recipe)
	if err != nil {
		Logger.Error("Error saving to database: ", err) // Log the actual error for debugging
		http.Error(w, "Error saving to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	outFile, err := os.Create("../Recipes/" + header.Filename) // save to the 'uploads' folder
	if err != nil {
		http.Error(w, "Error creating file on server", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe added successfully"})
	Logger.Info(json.NewEncoder(w).Encode(map[string]string{"message": "Recipe added successfully"}))
}
func GrantAdminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update user's role in the database
	filter := bson.M{"email": request.Email}
	update := bson.M{"$set": bson.M{"status": "admin"}}

	result, err := UserCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(w, "Error updating user role", http.StatusInternalServerError)
		return
	}
	if result.ModifiedCount == 0 {
		http.Error(w, "User not found or already an admin", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User granted admin rights successfully"})
}
