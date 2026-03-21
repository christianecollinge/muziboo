package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	url := os.Getenv("SUPABASE_URL") + "/rest/v1/profiles?select=*"
	req, _ := http.NewRequest("GET", url, nil)
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)

	resp, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %d\nProfiles: %s\n", resp.StatusCode, string(body))
}
