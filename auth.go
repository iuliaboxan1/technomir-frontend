package main

import (
    "fmt"
    "net/http"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    password := r.FormValue("password")
    user, err := SupabaseClient.Auth.SignUp(email, password)
    if err != nil {
        fmt.Fprintf(w, "Signup error: %v", err)
        return
    }
    fmt.Fprintf(w, "User created: %v", user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    password := r.FormValue("password")
    session, err := SupabaseClient.Auth.SignIn(email, password)
    if err != nil {
        fmt.Fprintf(w, "Login error: %v", err)
        return
    }
    fmt.Fprintf(w, "Logged in! Session: %v", session)
}

