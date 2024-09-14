package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// RequestData holds the data for generating a cURL command
type RequestData struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"` // Added headers field
}

// buildCurlCommand builds a cURL command string from the provided data
func buildCurlCommand(method, body, url string, headers map[string]string) string {
	var sb strings.Builder

	sb.WriteString("curl -X ")
	sb.WriteString(method)
	sb.WriteString(" ")

	// Add headers to the cURL command
	for key, value := range headers {
		sb.WriteString("-H '")
		sb.WriteString(key)
		sb.WriteString(": ")
		sb.WriteString(value)
		sb.WriteString("' ")
	}

	if body != "" {
		sb.WriteString("-d '")
		sb.WriteString(body)
		sb.WriteString("' ")
	}

	sb.WriteString(fmt.Sprintf("'%s'", url))

	return sb.String()
}

// generateCurlHandler handles HTTP requests to generate cURL commands
func generateCurlHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	
	var data RequestData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	curlCommand := buildCurlCommand(data.Method, data.Body, data.URL, data.Headers)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"curlCommand": curlCommand})
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/generate-curl", generateCurlHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
