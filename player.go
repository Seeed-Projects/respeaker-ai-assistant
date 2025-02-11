package main

import (
	"bytes"
	"io"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type Player struct {
	hasInit bool
}

func NewPlayer() *Player {
	return &Player{hasInit: false}
}

func (p *Player) PlayMp3(data []byte) error {
	reader := io.NopCloser(bytes.NewReader(data))
	streamer, _format, err := mp3.Decode(reader)
	if err != nil {
		return err
	}
	defer streamer.Close()

	if !p.hasInit {
		speaker.Init(_format.SampleRate, _format.SampleRate.N(time.Second/10))
		p.hasInit = true
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() { done <- true })))

	<-done
	return nil
}
