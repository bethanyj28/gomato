package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bethanyj28/gomato"
	"github.com/joho/godotenv"
	"github.com/nlopes/slack"
)

func main() {
	if err := godotenv.Load("environment.env"); err != nil {
		log.Fatal("failed to load env")
	}

	server := gomato.NewServer()

	if err := server.Router.Run(); err != nil {
		log.Fatal("server quit unexpectedly")
	}

	log.Print("server shutting down...")
}

func handleStartTimer(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !s.ValidateToken(os.Getenv("SLACK_VERIFICATION_TOKEN")) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case "/gomato":
		response := fmt.Sprintf("Starting timer for 10 seconds")
		w.Write([]byte(response))
	default:
		w.WriteHeader(http.StatusNoContent)
	}
}
