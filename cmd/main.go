package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/imdevinc/mockery/handlers"
	"github.com/imdevinc/mockery/services"
	"github.com/tmc/langchaingo/llms/openai"
)

func main() {
	rawPort := os.Getenv("PORT")
	port, err := strconv.Atoi(rawPort)
	if err != nil {
		slog.Error("failed to parse PORT environment variable. Using default of 8080", "error", err)
		port = 8080
	}
	userPrompt, err := os.ReadFile("prompts/user.txt")
	if err != nil {
		slog.Error("failed to read user prompt", "error", err)
		os.Exit(1)
	}
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		slog.Error("OPENAI_API_KEY environment variable is not set")
		os.Exit(1)
	}
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		slog.Error("OPENAI_MODEL environment variable is not set")
		os.Exit(1)
	}
	rawTemp := os.Getenv("LLM_TEMPERATURE")
	temp, err := strconv.ParseFloat(rawTemp, 64)
	if err != nil {
		slog.Error("failed to parse LLM_TEMPERATURE environment variable. Using default of 0.7", "error", err)
		temp = 0.7 // Default temperature if not set or invalid
	}
	llm, err := openai.New(openai.WithModel(model), openai.WithToken(openaiKey))
	if err != nil {
		slog.Error("failed to create OpenAI LLM", "error", err)
		os.Exit(1)
	}

	mockeryService := services.NewMockeryService(llm, string(userPrompt), temp)
	mockeryHandler := handlers.NewMockeryHandler(mockeryService)

	mux := http.NewServeMux()
	mockeryHandler.RegisterRoutes(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	slog.Info("starting server", "port", port)
	if err := http.ListenAndServe(":"+strconv.Itoa(port), mux); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
