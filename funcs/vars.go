package funcs

import (
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/time/rate"
)

var (
	UserCollection   *mongo.Collection
	RecipeCollection *mongo.Collection
	Store            *sessions.CookieStore
	rateLimiter      = rate.NewLimiter(1, 5)
	Logger           = logrus.New()
)

type Recipe struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Desc     string             `bson:"desc" json:"desc"`
	Calories float64            `bson:"calories" json:"calories"`
	Price    float64            `bson:"price" json:"price"`
	Level    string             `bson:"level" json:"level"`
	Path     string             `bson:"path" json:"path"`
}

type User struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Status     string             `bson:"status" json:"status"`
	RegMet     string             `bson:"reg_met" json:"reg_met"`
	Name       string             `bson:"name" json:"name"`
	Email      string             `bson:"email" json:"email"`
	Password   string             `bson:"password" json:"password"`
	AvatarPath string             `bson:"avatarPath" json:"avatarPath,omitempty"`
}
