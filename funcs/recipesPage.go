package funcs

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
)

func GetRecipesHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "projectGo")
	_, ok := session.Values["userID"]
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse query parameters for pagination, sorting, filtering, and price range
	page := r.URL.Query().Get("page")
	perPage := r.URL.Query().Get("perPage")
	sortBy := r.URL.Query().Get("sortBy")
	sortDirection := r.URL.Query().Get("sortDirection")
	levelFilter := r.URL.Query().Get("levelFilter")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")

	if page == "" {
		page = "1"
	}
	if perPage == "" {
		perPage = "6"
	}
	if sortDirection == "" {
		sortDirection = "asc"
	}

	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)
	if pageInt < 1 {
		pageInt = 1
	}
	if perPageInt < 1 {
		perPageInt = 6
	}

	filter := bson.M{}
	if levelFilter != "All" && levelFilter != "" {
		filter["level"] = levelFilter
	}
	if minPrice != "" && maxPrice != "" {
		minPriceFloat, _ := strconv.ParseFloat(minPrice, 64)
		maxPriceFloat, _ := strconv.ParseFloat(maxPrice, 64)
		filter["price"] = bson.M{"$gte": minPriceFloat, "$lte": maxPriceFloat}
	}

	sortOrder := 1
	if sortDirection == "desc" {
		sortOrder = -1
	}
	sortQuery := bson.M{sortBy: sortOrder}

	count, err := RecipeCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Failed to count Recipes", http.StatusInternalServerError)
		return
	}

	cursor, err := RecipeCollection.Find(context.TODO(), filter, options.Find().
		SetSort(sortQuery).
		SetSkip(int64((pageInt-1)*perPageInt)).
		SetLimit(int64(perPageInt)))
	if err != nil {

		http.Error(w, "Failed to fetch Recipes", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var Recipes []Recipe
	for cursor.Next(context.TODO()) {
		var Recipe Recipe
		if err := cursor.Decode(&Recipe); err != nil {
			http.Error(w, "Error decoding Recipe data", http.StatusInternalServerError)
			return
		}
		Recipes = append(Recipes, Recipe)
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Recipes    []Recipe `json:"Recipes"`
		TotalCount int      `json:"totalCount"`
	}{Recipes, int(count)}
	json.NewEncoder(w).Encode(response)
}
