package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// TOTPProvider defines an interface for storing and validating OTP secrets.
type TOTPProvider interface {
	StoreSecret(username, secret string) error
	Validate(username, code string) (bool, error)
}

// SecretVault is an in-memory store for OTP secrets.
type SecretVault map[string]string

// TOTPService struct
type TOTPService struct {
	userCache SecretVault
}

// NewTOTPService initializes the OTP service.
func NewTOTPService() TOTPService {
	return TOTPService{
		userCache: SecretVault{},
	}
}

// StoreSecret stores a user's OTP secret.
func (t *TOTPService) StoreSecret(username, secret string) error {
	t.userCache[username] = secret
	return nil
}

// Validate checks if the OTP code is correct.
func (t *TOTPService) Validate(username, code string) (bool, error) {

	secret, exists := t.userCache[username]

	if !exists {
		return false, fmt.Errorf("unknown user")
	}

	valid := totp.Validate(code, secret)
	if !valid {
		return false, fmt.Errorf("invalid code")
	}
	return true, nil
}

// OnBoard handler for registering users with OTP.
func OnBoard(t TOTPProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			username := r.FormValue("username")
			if username == "" {
				http.Error(w, "Username is required", http.StatusBadRequest)
				return
			}

			// Generate a secure OTP secret
			key, err := totp.Generate(totp.GenerateOpts{
				Issuer:      "MySecureApp",
				AccountName: username,
			})
			if err != nil {
				http.Error(w, "Failed to generate secret", http.StatusInternalServerError)
				return
			}

			// Store the secret
			if err := t.StoreSecret(username, key.Secret()); err != nil {
				http.Error(w, "Failed to store secret", http.StatusInternalServerError)
				return
			}

			// Generate a QR Code
			qrBytes, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
			if err != nil {
				http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
				return
			}

			qrBase64 := base64.StdEncoding.EncodeToString(qrBytes)

			// Render the QR Code page
			tmplParsed, err := template.New("confirmation").Parse(onboardTemplate)
			if err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			tmplParsed.Execute(w, qrBase64)
			return
		}

		// Render the onboarding form
		fmt.Fprint(w, onboardForm)
	}
}

// HomePage handler for OTP validation.
func HomePage(t TOTPProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			username := r.FormValue("username")
			code := r.FormValue("otp")

			if username == "" || code == "" {
				http.Error(w, "Username and OTP are required", http.StatusBadRequest)
				return
			}

			// Validate OTP
			valid, err := t.Validate(username, code)
			message := "Invalid OTP code. Please try again."
			if valid {
				message = "Login successful! Welcome back."
			}

			tmplParsed, err := template.New("result").Parse(resultTemplate)
			if err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			tmplParsed.Execute(w, message)
			return
		}

		// Render login form
		fmt.Fprint(w, loginForm)
	}
}

func main() {

	ts := NewTOTPService()

	mux := http.NewServeMux()
	mux.HandleFunc("/on-board", OnBoard(&ts))
	mux.HandleFunc("/", HomePage(&ts))

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Println("Listening on port 8000...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}

// HTML Templates
const onboardForm = `
	<h1>User Onboarding</h1>
	<form method="POST">
		<label for="username">Enter Username:</label>
		<input type="text" id="username" name="username" required>
		<button type="submit">Submit</button>
	</form>
`

const onboardTemplate = `
	<h1>Onboarding Complete</h1>
	<p>Scan the QR code below with your authenticator app:</p>
	<img src="data:image/png;base64,{{.}}" alt="QR Code">
	<p><a href="/">Go to Login</a></p>
`

const loginForm = `
	<h1>Login</h1>
	<form method="POST">
		<label for="username">Username:</label>
		<input type="text" id="username" name="username" required>
		<br>
		<label for="otp">One-Time Password (OTP):</label>
		<input type="text" id="otp" name="otp" required>
		<br>
		<button type="submit">Login</button>
	</form>
`

const resultTemplate = `
	<h1>Authentication Result</h1>
	<p>{{.}}</p>
	<a href="/">Go back to login</a>
`
