package main

import (
	"bytes"
	"io"

	"github.com/bclswl0827/respeaker-ai/lame"
)

func WavToMp3(wavData []byte) ([]byte, error) {
	wavReader := io.NopCloser(bytes.NewReader(wavData))
	wavHdr, err := lame.ReadWavHeader(wavReader)
	if err != nil {
		return nil, err
	}

	mp3Data := &bytes.Buffer{}
	wr, _ := lame.NewWriter(mp3Data)
	wr.EncodeOptions = wavHdr.ToEncodeOptions()
	io.Copy(wr, wavReader)
	wr.Close()

	return mp3Data.Bytes(), nil
}
