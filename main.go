package main

import (
	"encoding/base64"
    "encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"math"

	"github.com/NuclearLouse/tehnomir"
	"github.com/supabase-community/supabase-go"
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
    "github.com/supabase-community/postgrest-go"
)

var tpl = template.Must(template.ParseFiles("static/index.html"))
var cartTpl = template.Must(template.ParseFiles("static/cart.html"))
var accountTpl = template.Must(template.ParseFiles("static/account.html"))
var signupTpl = template.Must(template.ParseFiles("static/signup.html"))
var loginTpl = template.Must(template.ParseFiles("static/login.html"))
var confirmTpl = template.Must(template.ParseFiles("static/confirm.html"))
var supabaseClient *supabase.Client


// Get user info from Supabase token
// Decode JWT payload (middle part) and return email
func extractEmailFromJWT(token string) (string, error) {
    parts := strings.Split(token, ".")
    if len(parts) != 3 {
        return "", fmt.Errorf("invalid JWT format")
    }

    // Decode payload
    payload, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return "", fmt.Errorf("failed to decode JWT payload: %v", err)
    }

    // Parse JSON
    var data map[string]interface{}
    if err := json.Unmarshal(payload, &data); err != nil {
        return "", fmt.Errorf("failed to parse JWT JSON: %v", err)
    }

    // Return email claim
    email, ok := data["email"].(string)
    if !ok {
        return "", fmt.Errorf("email not found in JWT")
    }

    return email, nil
}




// Extract user_id ("sub") from Supabase JWT
func extractUserIDFromJWT(token string) (string, error) {
    parts := strings.Split(token, ".")
    if len(parts) < 2 {
        return "", fmt.Errorf("invalid token")
    }

    payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
    if err != nil {
        return "", err
    }

    var payload struct {
        Sub string `json:"sub"`
    }

    if err := json.Unmarshal(payloadJSON, &payload); err != nil {
        return "", err
    }

    return payload.Sub, nil
}










func main() {

    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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
// authClient := gotrue.New(projectRef, supabaseKey)
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

	// Try to create account
	_, err := authClient.Signup(req)
	if err != nil {
		errMsg := err.Error()

		// If Supabase says user already exists
		if strings.Contains(errMsg, "User already registered") ||
			strings.Contains(errMsg, "email already registered") {
			signupTpl.Execute(w, map[string]string{
				"Error": "An account with this email already exists. Please log in instead.",
			})
			return
		}

		http.Error(w, fmt.Sprintf("Signup error: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirect new user to confirmation page
	http.Redirect(w, r, "/confirm", http.StatusSeeOther)
})





	// ✅ Confirm route — place it right here!
http.HandleFunc("/confirm", func(w http.ResponseWriter, r *http.Request) {
	confirmTpl.Execute(w, nil)
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
        errMsg := err.Error()
        if strings.Contains(errMsg, "invalid_credentials") {
            message := "Wrong password or user does not exist. Please check your email or create an account first."
            loginTpl.Execute(w, map[string]string{"Error": message})
            return
        }
        http.Error(w, fmt.Sprintf("Login error: %v", err), http.StatusInternalServerError)
        return
    }

    // ⭐ Store ACCESS TOKEN
    http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    session.AccessToken,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // set true if using https
    })

    // ⭐ Store REFRESH TOKEN (important!!)
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    session.RefreshToken,
        Path:     "/",
        HttpOnly: true,
        Secure:   false, // set true if https
    })

    http.Redirect(w, r, "/", http.StatusSeeOther)
})





	
	// Logout
http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {

    // Remove access token
    http.SetCookie(w, &http.Cookie{
        Name:   "access_token",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    })

    // Remove refresh token
    http.SetCookie(w, &http.Cookie{
        Name:   "refresh_token",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    })

    http.Redirect(w, r, "/login", http.StatusSeeOther)
})



    // Account
// Account page
http.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {

    cookie, err := r.Cookie("access_token")
    if err != nil || cookie.Value == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    token := cookie.Value

    email, err := extractEmailFromJWT(token)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    userID, err := extractUserIDFromJWT(token)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // --- Fetch cart items ---
    cartData, _, err := supabaseClient.
        From("cart_items").
        Select("*", "", false).
        Eq("user_id", userID).
        Execute()

    if err != nil {
        http.Error(w, fmt.Sprintf("Error loading cart: %v", err), http.StatusInternalServerError)
        return
    }

    var cart []map[string]interface{}
    json.Unmarshal(cartData, &cart)

    // --- Calculate total ---
    total := 0.0
    for _, item := range cart {
        if priceFloat, ok := item["price"].(float64); ok {
            total += priceFloat
        } else if priceStr, ok := item["price"].(string); ok {
            var parsed float64
            fmt.Sscanf(priceStr, "%f", &parsed)
            total += parsed
        }
    }

    // --- Fetch orders ---
    orderData, _, err := supabaseClient.
        From("orders").
        Select("*", "", false).
        Eq("user_id", userID).
        Order("created_at", &postgrest.OrderOpts{Ascending: false}).
        Execute()

    if err != nil {
        http.Error(w, fmt.Sprintf("Error loading orders: %v", err), http.StatusInternalServerError)
        return
    }

    var orders []map[string]interface{}
    json.Unmarshal(orderData, &orders)

    // --- Render template ---
    accountTpl.Execute(w, map[string]interface{}{
        "Email":  email,
        "Cart":   cart,
        "Total":  total,
        "Orders": orders,
    })
})





	
	// Cart
http.HandleFunc("/cart", func(w http.ResponseWriter, r *http.Request) {


    cookie, err := r.Cookie("access_token")
    if err != nil || cookie.Value == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    token := cookie.Value


    email, err := extractEmailFromJWT(token)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    userID, err := extractUserIDFromJWT(token)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Fetch cart items
    data, _, err := supabaseClient.
        From("cart_items").
        Select("*", "", false).
        Eq("user_id", userID).
        Execute()

    if err != nil {
        http.Error(w, fmt.Sprintf("Error loading cart: %v", err), http.StatusInternalServerError)
        return
    }

    var cart []map[string]interface{}
    if err := json.Unmarshal(data, &cart); err != nil {
        http.Error(w, fmt.Sprintf("Error decoding cart JSON: %v", err), http.StatusInternalServerError)
        return
    }


    var total float64 = 0

for _, item := range cart {
    // price may arrive as float64 (most likely)
    if priceFloat, ok := item["price"].(float64); ok {
        total += priceFloat
        continue
    }

    // or as string (fallback)
    if priceStr, ok := item["price"].(string); ok {
        var priceParsed float64
        fmt.Sscanf(priceStr, "%f", &priceParsed)
        total += priceParsed
    }
}


    cartTpl.Execute(w, map[string]interface{}{
        "Email": email,
        "Cart":  cart,
        "Total": total,
    })

})





	
	// Add to cart
http.HandleFunc("/cart/add", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    // Read access token cookie (just for checking login)
    cookie, err := r.Cookie("access_token")
    if err != nil || cookie.Value == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Extract user ID from the token
    userID, err := extractUserIDFromJWT(cookie.Value)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    brand := r.FormValue("brand")
    code := r.FormValue("code")
    supplier := r.FormValue("supplier")
    price := r.FormValue("price")
    delivery := r.FormValue("delivery")

    payload := map[string]interface{}{
        "user_id":       userID,
        "brand":         brand,
        "code":          code,
        "supplier":      supplier,
        "price":         price,
        "delivery_days": delivery,
    }

    _, _, err = supabaseClient.
        From("cart_items").
        Insert(payload, false, "", "", "").
        Execute()

    if err != nil {
        http.Error(w, "Error adding to cart: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // http.Redirect(w, r, "/account", http.StatusSeeOther)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
})




// REMOVE from cart
http.HandleFunc("/cart/remove", func(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/cart", http.StatusSeeOther)
        return
    }

    // Must be logged in
    cookie, err := r.Cookie("access_token")
    if err != nil || cookie.Value == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    token := cookie.Value
    userID, err := extractUserIDFromJWT(token)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Which item to delete?
    id := r.FormValue("id") // cart row ID

    // DELETE only user’s own item
    _, _, err = supabaseClient.
        From("cart_items").
        Delete("", "").
        Eq("id", id).
        Eq("user_id", userID).
        Execute()

    if err != nil {
        http.Error(w, "Error deleting item: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/cart", http.StatusSeeOther)
})



    http.HandleFunc("/order/place", func(w http.ResponseWriter, r *http.Request) {

    cookie, err := r.Cookie("access_token")
    if err != nil || cookie.Value == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    token := cookie.Value

    userID, _ := extractUserIDFromJWT(token)
    // email, _ := extractEmailFromJWT(token)

    // Fetch cart
    data, _, err := supabaseClient.
        From("cart_items").
        Select("*", "", false).
        Eq("user_id", userID).
        Execute()

    if err != nil {
        http.Error(w, "Failed to load cart", 500)
        return
    }

    var cart []map[string]interface{}
    json.Unmarshal(data, &cart)

    if len(cart) == 0 {
        http.Error(w, "Cart is empty", 400)
        return
    }

    // Calculate total
    var total float64 = 0
    for _, item := range cart {
        if p, ok := item["price"].(float64); ok {
            total += p
        }
    }

    // Save order
    orderPayload := map[string]interface{}{
        "user_id": userID,
        "items":   cart,
        "total":   total,
    }

    _, _, err = supabaseClient.
        From("orders").
        Insert(orderPayload, false, "", "", "").
        Execute()

    if err != nil {
        http.Error(w, "Failed to save order: "+err.Error(), 500)
        return
    }

    // Clear cart
    _, _, err = supabaseClient.
        From("cart_items").
        Delete("", "").
        Eq("user_id", userID).
        Execute()

    if err != nil {
        http.Error(w, "Failed to clear cart", 500)
        return
    }

    // Send email
    //sendOrderEmail(email, cart, total)

    http.Redirect(w, r, "/account", http.StatusSeeOther)
})





	// Main search page (requires login)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check login
		cookie, err := r.Cookie("access_token")
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

		// fmt.Println("Search query:", partNum)
		// fmt.Println("API error:", err)
		// fmt.Println("Details count:", len(res.Details))
		// fmt.Printf("DETAIL[0]: %+v\n", res.Details[0])
		// fmt.Printf("STOCKS[0]: %+v\n", res.Details[0].Stocks[0])



		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		for i := range res.Details {
    	for j := range res.Details[i].Stocks {
        stock := &res.Details[i].Stocks[j]

        newPrice := stock.Price * 1.05
		stock.Price = math.Ceil(newPrice)
        stock.DeliveryDays = stock.DeliveryDays + 5
    	}
		}
		// tpl.Execute(w, res.Details)
		err = tpl.Execute(w, res.Details)
		if err != nil {
    	fmt.Println("TEMPLATE ERROR:", err)
   		 http.Error(w, err.Error(), 500)
   		 return
		}


	})





	
	// --- Start server ---
	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}





