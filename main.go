package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var args arguments
	err := args.Read()
	if err != nil {
		log.Fatal(err)
	}

	tts := NewTTS()
	player := NewPlayer()

	gemini, err := NewGemini(args.ApiKey, "gemini-2.0-flash")
	if err != nil {
		log.Fatal("Failed to create GenAI client:", err)
	}
	defer gemini.Close()

	recorder, err := NewRecorder()
	if err != nil {
		log.Fatal("Failed to create recorder:", err)
	}
	err = recorder.Init(44100, 1, 5000, 200*time.Millisecond, 1*time.Second, 1*time.Second)
	if err != nil {
		log.Fatal("Failed to init recorder:", err)
	}
	defer recorder.Deinit()

	recorder.OnVoiceDetected = func(t time.Time) {
		log.Println("Voice detected on", t.Format("15:04:05"))
	}
	recorder.OnVoiceStopped = func(t time.Time) {
		log.Println("Voice stopped on", t.Format("15:04:05"))
	}
	recorder.OnVoiceEvent = func(wavData []int) {
		wavCodec, err := recorder.ToWav(
			wavData,
			recorder.GetSampleRate(),
			recorder.GetBitDepth(),
			recorder.GetChannels(),
		)
		if err != nil {
			log.Println("Failed to get user voice WAV data:", err)
			return
		}
		mp3Codec, err := WavToMp3(wavCodec)
		if err != nil {
			log.Println("Failed to get user voice MP3 data:", err)
			return
		}

		log.Println("Generating AI response...")
		resp, err := gemini.Generate(
			args.Prompt,
			GeminiAttachment{FileName: "uservoice.wav", Data: mp3Codec, MIMEType: "audio/mpeg"},
		)
		if err != nil {
			log.Println("Failed to generate response:", err)
			return
		}
		log.Printf(">> %s\n", resp)

		log.Println("Generating voice response...")
		mp3Data, err := tts.Speak(url.QueryEscape(resp), args.Language)
		if err != nil {
			log.Println("Failed to fetch voice response:", err)
			return
		}

		log.Println("Playing voice response...")
		err = player.PlayMp3(mp3Data)
		if err != nil {
			log.Println("Failed to play voice response:", err)
			return
		}
		log.Println("Voice response played")
	}

	log.Println("Listening... (Press Ctrl+C to stop)")
	err = recorder.Start()
	if err != nil {
		log.Fatal("Failed to start recorder:", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
}
