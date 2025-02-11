package funcs

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var emailVerificationCodes = sync.Map{} // Простое хранение для примера

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	LoadLogger()
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	avatars := []string{
		"../avatars/001.png",
		"../avatars/002.png",
		"../avatars/003.png",
		"../avatars/004.png",
	}
	rand.Seed(time.Now().UnixNano())
	avatar := avatars[rand.Intn(len(avatars))]
	var verifyData struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Code     string `json:"code"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&verifyData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	storedCode, ok := emailVerificationCodes.Load(verifyData.Email)
	if !ok || storedCode != verifyData.Code {
		http.Error(w, "Invalid or expired code", http.StatusBadRequest)
		return
	}

	// Удаляем код из хранилища
	emailVerificationCodes.Delete(verifyData.Email)

	// Хешируем пароль
	hashedPassword, err := hashPassword(verifyData.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Сохраняем пользователя в базу данных
	user := User{Email: verifyData.Email, Password: hashedPassword, Name: verifyData.Name, Status: "user", RegMet: "email", AvatarPath: avatar}
	_, err = UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}
	err = performAction("registration", user.ID, "success")
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Registration successful"})
}
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	LoadLogger()
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var userData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {

		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверка наличия email в базе данных
	var existingUser User
	err := UserCollection.FindOne(context.TODO(), bson.M{"email": userData.Email}).Decode(&existingUser)
	if err == nil {
		// Email уже существует
		performAction("register", existingUser.ID, "fail")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "Email already exists"})
		return
	}

	// Генерация кода подтверждения
	verificationCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	emailVerificationCodes.Store(userData.Email, verificationCode)

	// Отправка кода подтверждения на почту
	go func() {
		err := sendEmailWithoutAttachment(
			userData.Email,
			"no-reply@Cooking-with-us.com",
			"Email Verification Code",
			"Your verification code is: "+verificationCode,
		)
		if err != nil {
			Logger.Errorf("Failed to send verification email: %v", err)
		}
	}()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Verification code sent to email"})
}

// Handles user login by validating credentials and saving session data

// Renders login page if user is not logged in

// Hashes the password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compares the plain password with the hashed on
