package funcs

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"os"
)

func LoadLogger() {
	// Создаем или открываем файл для логирования
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		Logger.Fatal("Ошибка при открытии файла для записи логов: ", err)
	}
	// Устанавливаем логгер на запись как в файл, так и в консоль
	multiWriter := io.Writer(file)
	Logger.SetOutput(multiWriter)
	// Устанавливаем формат логов в JSON
	Logger.SetFormatter(&logrus.JSONFormatter{})
}
func performAction(action string, userId primitive.ObjectID, result string) error {
	Logger.WithFields(logrus.Fields{
		"user_id": userId,
		"action":  action,
	}).Info(result)
	return nil
}
