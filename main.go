
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/NuclearLouse/tehnomir"
)

var tpl = template.Must(template.ParseFiles("static/index.html"))



func main() {
	token := os.Getenv("TECHNOMIR_TOKEN")
	if token == "" {
		log.Fatalln("Set TECHNOMIR_TOKEN environment variable with your API token")
	}

	cfg := tehnomir.DefaultConfig()
	cfg.Token = token
	client := tehnomir.New(cfg)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}



