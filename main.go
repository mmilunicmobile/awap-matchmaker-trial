package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func fetchURL(wg *sync.WaitGroup, url string) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Fetched %s with status code: %d\n", url, resp.StatusCode)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	db, err := sqlx.Connect("sqlite", "file:example.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	urls := []string{
		"https://www.example.com",
		"https://www.google.com",
		"https://www.githubbbbb.com",
	}

	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go fetchURL(&wg, u)
	}
	wg.Wait()

	http.HandleFunc("/", helloHandler)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
