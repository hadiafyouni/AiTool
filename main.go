package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sashabaranov/go-openai"
)

// Request structure
type ChatRequest struct {
	Prompt string `json:"prompt"`
}

// Response structure
type ChatResponse struct {
	Answer string `json:"answer"`
}

func main() {
	http.HandleFunc("/ask", handleChat)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Backend is running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	// 1. CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 2. Decode JSON with Error Checking
	var req ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("❌ Decode Error:", err)
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Logic check: don't call AI if prompt is empty
	if req.Prompt == "" {
		http.Error(w, "Prompt cannot be empty", http.StatusBadRequest)
		return
	}

	fmt.Printf("📝 User asked: %s\n", req.Prompt)

	// 3. Groq Configuration
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		http.Error(w, "Server Misconfiguration: Missing API Key", http.StatusInternalServerError)
		return
	}
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.groq.com/openai/v1"
	client := openai.NewClientWithConfig(config)

	// 4. Call Groq (Update the model string here)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			// Change "llama3-8b-8192" to "llama-3.1-8b-instant"
			Model: "llama-3.1-8b-instant",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert Go and Next.js developer. Always provide code examples.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: req.Prompt,
				},
			},
		},
	)

	// 5. Handle Groq API Errors
	if err != nil {
		fmt.Println("❌ Groq API Error:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ChatResponse{
			Answer: "Groq is busy or your key is invalid. Error: " + err.Error(),
		})
		return
	}

	// 6. Final Success Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{
		Answer: resp.Choices[0].Message.Content,
	})
	fmt.Println("✅ Response sent successfully!")
}


