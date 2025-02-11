package tests

import (
	"awesomeProject3/funcs"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"net/http"
	"net/http/httptest"
	"testing"
)

// Структура для хранения данных о Pokémon
type Recipe struct {
	ID    string  `json:"id" bson:"_id"`
	Name  string  `json:"name" bson:"name"`
	Price float64 `json:"price" bson:"price"`
}

// Функция для обработки запроса
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:Delson_action@cluster0.lpci1.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0").SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		funcs.Logger.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	funcs.RecipeCollection = client.Database("projectGo").Collection("Recipes")

	var recipe Recipe
	err = funcs.RecipeCollection.FindOne(ctx, bson.M{"name": "Caesar Salad 5"}).Decode(&recipe)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Рецепт не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка при поиске в базе данных", http.StatusInternalServerError)
		}
		return
	}

	// Установка заголовков и отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func TestGetRecipe(t *testing.T) {
	fmt.Println("TESTING INTEGRATION TEST")
	// Создание запроса
	request, err := http.NewRequest("GET", "/recipes/679f62e62bf1adc5a36d6073", nil)
	if err != nil {
		t.Fatalf("Не удалось создать запрос: %v", err)
	}

	// Создание фиктивного респондера
	response := httptest.NewRecorder()

	// Вызов обработчика эндпоинта
	GetRecipe(response, request)

	// Проверка кода состояния ответа
	if response.Code != http.StatusOK {
		t.Errorf("Неверный код состояния. Ожидалось: %d, Получено: %d", http.StatusOK, response.Code)
	}

	// Проверка тела ответа
	var actual Recipe
	err = json.Unmarshal(response.Body.Bytes(), &actual)
	if err != nil {
		t.Fatalf("Не удалось разобрать тело ответа: %v", err)
	}

	expected := Recipe{
		ID:    "679f62e62bf1adc5a36d6073",
		Name:  "Caesar Salad 5",
		Price: 12.97,
	}

	if actual != expected {
		t.Errorf("Неверное тело ответа. Ожидалось: %+v, Получено: %+v", expected, actual)
	}
}
