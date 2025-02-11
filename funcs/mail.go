package funcs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

// Handles sending email requests
func SendEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data (multipart form)
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for file uploads
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	var emailData struct {
		Subject string `json:"subject"`
		Body    string `json:"body"`
		ToEmail string `json:"toEmail"`
	}

	// Decode email data
	err = json.NewDecoder(strings.NewReader(r.FormValue("emailData"))).Decode(&emailData)
	if err != nil {
		http.Error(w, "Invalid email data", http.StatusBadRequest)
		return
	}

	// Get user email from session
	session, _ := Store.Get(r, "projectGo")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "User not logged in", http.StatusUnauthorized)
		return
	}

	// Retrieve user data from the database
	var user User
	objectID, _ := primitive.ObjectIDFromHex(userID)
	err = UserCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if a file is attached
	var fileBytes []byte
	var fileName string
	file, _, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
		fileBytes, err = ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		fileName = r.FormValue("fileName")
	}

	// Send email
	if fileBytes != nil {
		err = sendEmailWithAttachment(emailData.ToEmail, user.Email, emailData.Subject, emailData.Body, fileName, fileBytes)
	} else {
		err = sendEmailWithoutAttachment(emailData.ToEmail, user.Email, emailData.Subject, emailData.Body)
	}

	if err != nil {
		log.Printf("Error sending email: %v", err)
		http.Error(w, "Failed to send email. Please check server logs.", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Email sent successfully"})
}

// Send email with attachment
func sendEmailWithAttachment(toEmail, fromEmail, subject, body, fileName string, fileBytes []byte) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")

	// Create email boundary
	boundary := "myboundary"

	// Create the email body with the text and attachment in a simplified MIME structure
	bodyContent := fmt.Sprintf(`Subject: %s
Content-Type: multipart/mixed; boundary="%s"

--%s
Content-Type: text/plain; charset=UTF-8

%s

--%s
Content-Type: application/octet-stream; name="%s"
Content-Disposition: attachment; filename="%s"
Content-Transfer-Encoding: base64

%s
--%s--`, subject, boundary, boundary, body, boundary, fileName, fileName, base64.StdEncoding.EncodeToString(fileBytes), boundary)

	// Setup the SMTP authentication
	auth := smtp.PlainAuth("", username, password, smtpHost)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{toEmail}, []byte(bodyContent))
	return err
}

// Send email without attachment
func sendEmailWithoutAttachment(toEmail, fromEmail, subject, text string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	body := fmt.Sprintf("Subject: %s\n\n%s\n\nFrom: %s", subject, text, fromEmail)

	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{toEmail}, []byte(body))
	return err
}
