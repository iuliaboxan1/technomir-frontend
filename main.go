// package main

// import (
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"

// 	"github.com/NuclearLouse/tehnomir"
// )

// var tpl = template.Must(template.ParseFiles("static/index.html"))



// func main() {
// 	token := os.Getenv("TECHNOMIR_TOKEN")
// 	if token == "" {
// 		log.Fatalln("Set TECHNOMIR_TOKEN environment variable with your API token")
// 	}

// 	cfg := tehnomir.DefaultConfig()
// 	cfg.Token = token
// 	client := tehnomir.New(cfg)

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodGet {
// 			tpl.Execute(w, nil)
// 			return
// 		}

// 		// Handle POST
// 		if err := r.ParseForm(); err != nil {
// 			fmt.Fprintf(w, "ParseForm error: %v", err)
// 			return
// 		}
// 		partNum := strings.TrimSpace(r.FormValue("partNum"))
// 		if partNum == "" {
// 			fmt.Fprintf(w, "Part number is required")
// 			return
// 		}

// 		res, err := client.SearchByBrandWithoutAnalogs(partNum, 0, tehnomir.USD)
// 		if err != nil {
// 			fmt.Fprintf(w, "Error: %v", err)
// 			return
// 		}

// 		tpl.Execute(w, res.Details)
// 	})

// 	fmt.Println("Server started at http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }


// package main

// import (
// 	"fmt"
// 	"html/template"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"

// 	"github.com/NuclearLouse/tehnomir"
// 	"github.com/supabase-community/gotrue-go/types"
// 	"github.com/supabase-community/supabase-go"
// )

// var tpl = template.Must(template.ParseFiles("static/index.html"))
// var supabaseClient *supabase.Client

// func main() {
// 	// --- Technomir setup ---
// 	token := os.Getenv("TECHNOMIR_TOKEN")
// 	if token == "" {
// 		log.Fatalln("Set TECHNOMIR_TOKEN environment variable with your API token")
// 	}
// 	cfg := tehnomir.DefaultConfig()
// 	cfg.Token = token
// 	client := tehnomir.New(cfg)

// 	// --- Supabase setup ---
// 	supabaseURL := os.Getenv("SUPABASE_URL")
// 	supabaseKey := os.Getenv("SUPABASE_KEY")
// 	if supabaseURL == "" || supabaseKey == "" {
// 		log.Fatalln("Set SUPABASE_URL and SUPABASE_KEY environment variables")
// 	}

// 	var err error
// 	supabaseClient, err = supabase.NewClient(supabaseURL, supabaseKey, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to create Supabase client: %v", err)
// 	}

// 	// --- Routes ---
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method == http.MethodGet {
// 			tpl.Execute(w, nil)
// 			return
// 		}

// 		if err := r.ParseForm(); err != nil {
// 			fmt.Fprintf(w, "ParseForm error: %v", err)
// 			return
// 		}
// 		partNum := strings.TrimSpace(r.FormValue("partNum"))
// 		if partNum == "" {
// 			fmt.Fprintf(w, "Part number is required")
// 			return
// 		}

// 		res, err := client.SearchByBrandWithoutAnalogs(partNum, 0, tehnomir.USD)
// 		if err != nil {
// 			fmt.Fprintf(w, "Error: %v", err)
// 			return
// 		}

// 		tpl.Execute(w, res.Details)
// 	})

// 	http.HandleFunc("/signup", SignupHandler)
// 	http.HandleFunc("/login", LoginHandler)

// 	// --- Start server ---
// 	fmt.Println("Server started at http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// // ---------------- Supabase Handlers ----------------
// func SignupHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	email := strings.TrimSpace(r.FormValue("email"))
// 	password := strings.TrimSpace(r.FormValue("password"))
// 	if email == "" || password == "" {
// 		http.Error(w, "Email and password required", http.StatusBadRequest)
// 		return
// 	}

// 	req := types.SignupRequest{
// 		Email:    email,
// 		Password: password,
// 	}

// 	user, err := supabaseClient.Auth.Signup(req)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Signup error: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprintf(w, "User created: %v", user)
// }

// func LoginHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	email := strings.TrimSpace(r.FormValue("email"))
// 	password := strings.TrimSpace(r.FormValue("password"))
// 	if email == "" || password == "" {
// 		http.Error(w, "Email and password required", http.StatusBadRequest)
// 		return
// 	}

// 	session, err := supabaseClient.Auth.SignInWithEmailPassword(email, password)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("Login error: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Fprintf(w, "Logged in! Session: %v", session)
// }


package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/NuclearLouse/tehnomir"
	"github.com/supabase-community/supabase-go"
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
)

var tpl = template.Must(template.ParseFiles("static/index.html"))
var signupTpl = template.Must(template.ParseFiles("static/signup.html"))
var loginTpl = template.Must(template.ParseFiles("static/login.html"))
var supabaseClient *supabase.Client



func main() {
	// --- Technomir setup ---
	token := os.Getenv("TECHNOMIR_TOKEN")
	if token == "" {
		log.Fatalln("Set TECHNOMIR_TOKEN environment variable with your API token")
	}

	cfg := tehnomir.DefaultConfig()
	cfg.Token = token
	client := tehnomir.New(cfg)

	// --- Supabase setup ---
supabaseURL := os.Getenv("SUPABASE_URL")
supabaseKey := os.Getenv("SUPABASE_KEY")
if supabaseURL == "" || supabaseKey == "" {
	log.Fatalln("Set SUPABASE_URL and SUPABASE_KEY environment variables")
}

var err error
supabaseClient, err = supabase.NewClient(supabaseURL, supabaseKey, nil)
if err != nil {
	log.Fatalf("Error initializing Supabase client: %v", err)
}

// --- Gotrue Auth setup ---
projectRef := "xgrmgyusghkuogfbkkcl" // this is your Supabase project ref
authClient := gotrue.New(projectRef, supabaseKey)
fmt.Println("Auth client initialized for project:", projectRef)







	// --- Routes ---




	// Signup page
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			signupTpl.Execute(w, nil)
			return
		}

		email := strings.TrimSpace(r.FormValue("email"))
		password := strings.TrimSpace(r.FormValue("password"))

		req := types.SignupRequest{
			Email:    email,
			Password: password,
		}

		user, err := authClient.Signup(req)
		if err != nil {
			http.Error(w, fmt.Sprintf("Signup error: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User created: %v", user)
	})




	// Login page
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loginTpl.Execute(w, nil)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	session, err := authClient.SignInWithEmailPassword(email, password)
	if err != nil {
    	http.Error(w, fmt.Sprintf("Login error: %v", err), http.StatusInternalServerError)
    	return
	}


	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: session.AccessToken,
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
})



	
	// Logout
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			Path:   "/",
			MaxAge: -1, // delete cookie
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	// Main search page (requires login)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check login
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if r.Method == http.MethodGet {
			tpl.Execute(w, nil)
			return
		}

		// Handle POST
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm error: %v", err)
			return
		}
		partNum := strings.TrimSpace(r.FormValue("partNum"))
		if partNum == "" {
			fmt.Fprintf(w, "Part number is required")
			return
		}

		res, err := client.SearchByBrandWithoutAnalogs(partNum, 0, tehnomir.USD)
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		tpl.Execute(w, res.Details)
	})

	// --- Start server ---
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}



