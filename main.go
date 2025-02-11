package main

import (
	"awesomeProject3/funcs"
	"context"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	funcs.LoadLogger()
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables from system")
	}
	// Setup logging configuration
	funcs.Logger.SetFormatter(&logrus.JSONFormatter{})
	funcs.Logger.SetOutput(os.Stdout)
	funcs.Logger.SetLevel(logrus.InfoLevel)

	// Graceful shutdown setup
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-sigs
		funcs.Logger.Info("Graceful shutdown initiated")
		cancel()
	}()

	// Load environment variables from .env file

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		funcs.Logger.Fatal("SESSION_KEY environment variable is not set")
	}

	funcs.Store = sessions.NewCookieStore([]byte(sessionKey))
	funcs.Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	// MongoDB connection setup
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGOURL")).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		funcs.Logger.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		funcs.Logger.Fatal(err)
	}
	clientid := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SEC")

	redirectURL := os.Getenv("REDIRECT_URL")
	if redirectURL == "" {
		funcs.Logger.Fatal("REDIRECT_URL environment variable is not set")
	}
	funcs.Logger.Info("Using Redirect URL: ", redirectURL)

	conf := &oauth2.Config{
		ClientID:     clientid,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL, // Use the environment variable
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
	app := funcs.App{Config: conf}
	funcs.UserCollection = client.Database("projectGo").Collection("users")
	funcs.RecipeCollection = client.Database("projectGo").Collection("Recipes")
	funcs.Logger.Info("Connected to MongoDB!")

	// Start the HTTP server and handle routes

	http.HandleFunc("/register", funcs.RateLimitMiddleware(funcs.RegisterHandler))
	http.HandleFunc("/login", funcs.RateLimitMiddleware(funcs.Login))
	http.HandleFunc("/auth/oauth", funcs.RateLimitMiddleware(app.OAuthHandler))
	http.HandleFunc("/auth/callback", app.OAuthCallbackHandler)
	http.HandleFunc("/logout", funcs.RateLimitMiddleware(funcs.Logout))

	http.HandleFunc("/verifyEmail", funcs.RateLimitMiddleware(funcs.VerifyEmailHandler))
	http.HandleFunc("/sendEmail", funcs.RateLimitMiddleware(funcs.SendEmail))

	http.HandleFunc("/addRecipe", funcs.AddRecipeHandler)
	http.HandleFunc("/deleteRecipe", funcs.DeleteRecipeHandler)
	http.HandleFunc("/grantAdmin", funcs.GrantAdminHandler)

	http.HandleFunc("/checkLoginStatus", funcs.RateLimitMiddleware(funcs.CheckLoginStatus))
	http.HandleFunc("/recipes", funcs.RateLimitMiddleware(funcs.GetRecipesHandler))
	http.HandleFunc("/riskyOperation", funcs.RateLimitMiddleware(funcs.RiskyOperationHandler))

	http.HandleFunc("/updatePassword", funcs.UpdatePassword)
	http.HandleFunc("/updateName", funcs.UpdateName)
	http.HandleFunc("/updateAvatar", funcs.RateLimitMiddleware(funcs.UpdateAvatarHandler))

	//Pages
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect the root URL to the login page
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/loginPage", http.StatusFound)
		} else {
			http.ServeFile(w, r, "./404.html")
		}
	})
	http.HandleFunc("/verifyEmailPage", funcs.RateLimitMiddleware(funcs.VerifyEmailPage))
	http.HandleFunc("/loginPage", funcs.RateLimitMiddleware(funcs.ServeLogin))
	http.HandleFunc("/adminPage", funcs.RateLimitMiddleware(funcs.ServeAdminPage))
	http.HandleFunc("/mainPage", funcs.RateLimitMiddleware(funcs.ServeMain))
	http.HandleFunc("/registerPage", funcs.RateLimitMiddleware(funcs.ServeRegister))
	http.HandleFunc("/recipesPage", funcs.RateLimitMiddleware(funcs.ServeRecipesPage))
	http.HandleFunc("/userPage", funcs.UserPage)

	//GETTERS
	http.HandleFunc("/getUserRegMet", funcs.RateLimitMiddleware(funcs.GetUserRegMetHandler))
	http.HandleFunc("/getUserEmail", funcs.RateLimitMiddleware(funcs.GetUserEmailHandler))
	http.HandleFunc("/getUserStatus", funcs.RateLimitMiddleware(funcs.GetUserStatusHandler))
	http.HandleFunc("/getUserName", funcs.RateLimitMiddleware(funcs.GetUserNameHandler))
	http.HandleFunc("/getUserAvatar", funcs.RateLimitMiddleware(funcs.GetUserAvatar))

	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("./style"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/Recipes/", http.StripPrefix("/Recipes/", http.FileServer(http.Dir("./Recipes"))))
	port := os.Getenv("PORT") // Fetch the port from the environment variable
	if port == "" {           // If the PORT variable is not set, use the default port
		port = "8080"
	}
	link := ":" + port // Bind to all available network interfaces

	server := &http.Server{
		Addr:    link,
		Handler: http.DefaultServeMux,
	}

	funcs.Logger.Info("Server running on http://localhost:", port)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed && nil != err {
			funcs.Logger.Fatalf("Server failed: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		funcs.Logger.Errorf("Error during server shutdown: %v", err)
	}
	funcs.Logger.Info("Server gracefully stopped")
}
