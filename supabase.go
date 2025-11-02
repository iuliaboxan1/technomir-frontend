package main

import (
	"log"
	"os"

	"github.com/supabase-community/supabase-go"
)

var SupabaseClient *supabase.Client

func InitSupabase() {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	if url == "" || key == "" {
		log.Fatalln("Set SUPABASE_URL and SUPABASE_KEY environment variables")
	}

	SupabaseClient = supabase.CreateClient(url, key)
}
