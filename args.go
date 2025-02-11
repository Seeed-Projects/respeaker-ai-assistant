package main

import (
	"errors"
	"flag"
)

type arguments struct {
	ApiKey   string
	Language string
	Prompt   string
}

func (a *arguments) Read() error {
	flag.StringVar(&a.ApiKey, "key", "", "Your Gemini API key for authentication")
	flag.StringVar(&a.Language, "lang", "en", "Response language in ISO 639-1 format (e.g., en, ja, zh-TW)")
	flag.StringVar(&a.Prompt, "prompt", "Respond to the user's audio, keep it short and use the same language as the user, no emoji.", "Prompt for response generation")
	flag.Parse()

	if a.ApiKey == "" {
		return errors.New("API key is required, see `-help` for usage")
	}

	return nil
}
