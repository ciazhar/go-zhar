package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
)

type TOTPProvider interface {
	StoreSecret(username, secret string) error
	Validate(username, code string) (bool, error)
}

type SecretVault map[string]string

// TOTPService struct
type TOTPService struct {
	userCache SecretVault
}

func NewTOTPService() TOTPService {
	return TOTPService{
		userCache: SecretVault{},
	}
}

// StoreSecret interface method
func (t *TOTPService) StoreSecret(username, secret string) error {
	t.userCache[username] = secret
	return nil
}

// Validate interface method
func (t *TOTPService) Validate(username, code string) (bool, error) {

	secret, exists := t.userCache[username]

	if !exists {
		return false, fmt.Errorf("unknown user")
	}

	totp := gotp.NewDefaultTOTP(secret)

	if totp.Now() != code {
		return false, fmt.Errorf("Invalid code")
	}

	return true, nil
}

// OnBoard handler
func OnBoard(t TOTPProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			username := r.FormValue("username")
			if username == "" {
				http.Error(w, "Username is required", http.StatusBadRequest)
				return
			}

			// Generate a new secret
			secretLength := 16
			secret := gotp.RandomSecret(secretLength)

			// Store the secret using the TOTPProvider
			if err := t.StoreSecret(username, secret); err != nil {
				http.Error(w, "Failed to store secret", http.StatusInternalServerError)
				return
			}

			// Generate a QR code in base64 format
			totp := gotp.NewDefaultTOTP(secret)
			provUri := totp.ProvisioningUri(username, "myOTPApp")

			qrBytes, err := qrcode.Encode(provUri, qrcode.Medium, 256)
			if err != nil {
				http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
				return
			}

			qrBase64 := base64.StdEncoding.EncodeToString(qrBytes)

			// Render the QR code and confirmation page
			tmpl := `
				<h1>Onboarding Complete</h1>
				<p>Scan the QR code below with your authenticator app:</p>
				<img src="data:image/png;base64,{{.}}" alt="QR Code">
			`
			tmplParsed, err := template.New("confirmation").Parse(tmpl)
			if err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			tmplParsed.Execute(w, qrBase64)
			return
		}

		// Render the default user form
		tmpl := `
			<h1>User Onboarding</h1>
			<form method="POST">
				<label for="username">Enter Username:</label>
				<input type="text" id="username" name="username" required>
				<button type="submit">Submit</button>
			</form>
		`
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, tmpl)
	}
}

// HomePage handler
func HomePage(t TOTPProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Failed to parse form", http.StatusBadRequest)
				return
			}

			// Get username and OTP code from the form
			username := r.FormValue("username")
			code := r.FormValue("otp")

			if username == "" || code == "" {
				http.Error(w, "Username and OTP code are required", http.StatusBadRequest)
				return
			}

			// Validate the OTP using the TOTPProvider
			valid, err := t.Validate(username, code)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error validating OTP: %v", err), http.StatusInternalServerError)
				return
			}

			// Render a success or failure message
			var message string
			if valid {
				message = "Login successful! Welcome back."
			} else {
				message = "Invalid OTP code. Please try again."
			}

			tmpl := `
				<h1>Authentication Result</h1>
				<p>{{.}}</p>
				<a href="/">Go back to login</a>
			`
			tmplParsed, err := template.New("result").Parse(tmpl)
			if err != nil {
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			tmplParsed.Execute(w, message)
			return
		}

		// Render the login form for GET requests
		tmpl := `
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
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, tmpl)
	}
}

func main() {

	ts := NewTOTPService()

	http.HandleFunc(
		"/on-board",
		OnBoard(&ts),
	)

	http.HandleFunc(
		"/",
		HomePage(&ts),
	)

	log.Println("Listenning on port 8000")
	http.ListenAndServe(":8000", nil)
}
