// package main

// import (
// 	"bytes"
// 	"time"
// 	"unsafe"

// 	"github.com/gen2brain/malgo"
// 	"github.com/go-audio/audio"
// 	"github.com/go-audio/wav"
// )

// type Recorder struct {
// 	pCapturedSamples   []byte
// 	isRecording        bool
// 	lastVoiceTime      time.Time
// 	recordingStartTime time.Time
// 	buffer             *FifoBuffer

// 	volumeThreshold int16
// 	triggerDuration time.Duration
// 	silenceTimeout  time.Duration

// 	ctx    *malgo.AllocatedContext
// 	config malgo.DeviceConfig
// 	device *malgo.Device

// 	OnSampleData    func(data []int, samples int)
// 	OnVoiceDetected func(t time.Time)
// 	OnVoiceStopped  func(t time.Time)
// 	OnVoiceEvent    func(wavData []int)
// }

// func NewRecorder() (*Recorder, error) {
// 	obj := &Recorder{}
// 	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {})

// 	if err != nil {
// 		return nil, err
// 	}

// 	obj.ctx = ctx
// 	return obj, nil
// }

// func (r *Recorder) Init(sampleRate, channels, volumeThreshold int, triggerDuration, bufferDuration, silenceTimeout time.Duration) error {
// 	r.volumeThreshold = int16(volumeThreshold)
// 	r.triggerDuration = triggerDuration
// 	r.silenceTimeout = silenceTimeout

// 	r.config = malgo.DefaultDeviceConfig(malgo.Duplex)
// 	r.config.Capture.Format = malgo.FormatS16
// 	r.config.Capture.Channels = uint32(channels)
// 	r.config.SampleRate = uint32(sampleRate)
// 	r.config.Alsa.NoMMap = 1

// 	bitDepth := 16
// 	bytesPerFrame := (bitDepth / 8) * channels
// 	r.buffer = NewFifoBuffer(int(bufferDuration.Seconds()) * sampleRate * bytesPerFrame)

// 	device, err := malgo.InitDevice(r.ctx.Context, r.config, malgo.DeviceCallbacks{
// 		Data: func(_, pInputSamples []byte, frameCount uint32) {
// 			int16Data := unsafe.Slice((*int16)(unsafe.Pointer(&pInputSamples[0])), len(pInputSamples)/2)

// 			// Check for max amplitude
// 			var maxAmplitude int16 = 0
// 			for _, sample := range int16Data {
// 				absSample := sample
// 				if sample < 0 {
// 					absSample = -sample
// 				}
// 				if absSample > maxAmplitude {
// 					maxAmplitude = absSample
// 				}
// 			}

// 			// Write to buffer
// 			_, _ = r.buffer.Write(pInputSamples)

// 			// Detecting voice
// 			now := time.Now()
// 			if maxAmplitude > r.volumeThreshold {
// 				if !r.isRecording {
// 					if r.lastVoiceTime.IsZero() {
// 						r.lastVoiceTime = now
// 					}
// 					if now.Sub(r.lastVoiceTime) >= r.triggerDuration {
// 						if r.OnVoiceDetected != nil {
// 							r.OnVoiceDetected(now)
// 						}

// 						r.pCapturedSamples = r.buffer.Read() // Read from buffer
// 						r.isRecording = true
// 						r.recordingStartTime = now
// 					}
// 				} else {
// 					r.lastVoiceTime = now
// 				}
// 			}

// 			if r.isRecording {
// 				r.pCapturedSamples = append(r.pCapturedSamples, pInputSamples...)
// 			}

// 			// Check if voice stopped
// 			if r.isRecording && now.Sub(r.lastVoiceTime) >= silenceTimeout {
// 				if r.OnVoiceStopped != nil {
// 					int16Data := unsafe.Slice((*int16)(unsafe.Pointer(&r.pCapturedSamples[0])), len(r.pCapturedSamples)/2)
// 					intData := make([]int, len(int16Data))
// 					for i, v := range int16Data {
// 						intData[i] = int(v)
// 					}
// 					r.OnVoiceStopped(now)
// 					r.OnVoiceEvent(intData)
// 				}
// 				r.isRecording = false
// 				r.lastVoiceTime = time.Time{}
// 			}
// 		},
// 	})
// 	if err != nil {
// 		r.ctx.Uninit()
// 		return err
// 	}

// 	r.device = device
// 	return nil
// }

// func (r *Recorder) Start() error {
// 	return r.device.Start()
// }

// func (r *Recorder) Deinit() {
// 	if r.device != nil {
// 		r.device.Uninit()
// 	}
// 	r.ctx.Uninit()
// 	r.ctx.Free()
// }

// func (r *Recorder) GetSampleRate() int {
// 	return int(r.config.SampleRate)
// }

// func (r *Recorder) GetChannels() int {
// 	return int(r.config.Capture.Channels)
// }

// func (r *Recorder) GetBitDepth() int {
// 	return 16
// }

// func (r *Recorder) ToWav(wavData []int, sampleRate, bitDepth, channels int) ([]byte, error) {
// 	buf := &bytes.Buffer{}
// 	ws := &nopSeek{buf}

// 	encoder := wav.NewEncoder(ws, sampleRate, bitDepth, channels, 1)
// 	defer encoder.Close()

// 	err := encoder.Write(&audio.IntBuffer{
// 		Format:         &audio.Format{NumChannels: channels, SampleRate: sampleRate},
// 		Data:           wavData,
// 		SourceBitDepth: bitDepth,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return buf.Bytes(), nil
// }

package main

import (
	"bytes"
	"time"
	"unsafe"

	"github.com/gen2brain/malgo"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type Recorder struct {
	pCapturedSamples   []byte
	isRecording        bool
	lastVoiceTime      time.Time
	recordingStartTime time.Time
	buffer             *FifoBuffer

	volumeThreshold int16
	triggerDuration time.Duration
	silenceTimeout  time.Duration

	ctx    *malgo.AllocatedContext
	config malgo.DeviceConfig
	device *malgo.Device

	OnSampleData    func(data []int, samples int)
	OnVoiceDetected func(t time.Time)
	OnVoiceStopped  func(t time.Time)
	OnVoiceEvent    func(wavData []int)

	voiceEventQueue chan []int
}

func NewRecorder() (*Recorder, error) {
	obj := &Recorder{
		voiceEventQueue: make(chan []int, 10), // Buffer size for the queue
	}
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {})

	if err != nil {
		return nil, err
	}

	obj.ctx = ctx
	go obj.processVoiceEvents()
	return obj, nil
}

func (r *Recorder) Init(sampleRate, channels, volumeThreshold int, triggerDuration, bufferDuration, silenceTimeout time.Duration) error {
	r.volumeThreshold = int16(volumeThreshold)
	r.triggerDuration = triggerDuration
	r.silenceTimeout = silenceTimeout

	r.config = malgo.DefaultDeviceConfig(malgo.Duplex)
	r.config.Capture.Format = malgo.FormatS16
	r.config.Capture.Channels = uint32(channels)
	r.config.SampleRate = uint32(sampleRate)
	r.config.Alsa.NoMMap = 1

	bitDepth := 16
	bytesPerFrame := (bitDepth / 8) * channels
	r.buffer = NewFifoBuffer(int(bufferDuration.Seconds()) * sampleRate * bytesPerFrame)

	device, err := malgo.InitDevice(r.ctx.Context, r.config, malgo.DeviceCallbacks{
		Data: func(_, pInputSamples []byte, frameCount uint32) {
			int16Data := unsafe.Slice((*int16)(unsafe.Pointer(&pInputSamples[0])), len(pInputSamples)/2)

			// Check for max amplitude
			var maxAmplitude int16 = 0
			for _, sample := range int16Data {
				absSample := sample
				if sample < 0 {
					absSample = -sample
				}
				if absSample > maxAmplitude {
					maxAmplitude = absSample
				}
			}

			// Write to buffer
			_, _ = r.buffer.Write(pInputSamples)

			// Detecting voice
			now := time.Now()
			if maxAmplitude > r.volumeThreshold {
				if !r.isRecording {
					if r.lastVoiceTime.IsZero() {
						r.lastVoiceTime = now
					}
					if now.Sub(r.lastVoiceTime) >= r.triggerDuration {
						if r.OnVoiceDetected != nil {
							r.OnVoiceDetected(now)
						}

						r.pCapturedSamples = r.buffer.Read() // Read from buffer
						r.isRecording = true
						r.recordingStartTime = now
					}
				} else {
					r.lastVoiceTime = now
				}
			}

			if r.isRecording {
				r.pCapturedSamples = append(r.pCapturedSamples, pInputSamples...)
			}

			// Check if voice stopped
			if r.isRecording && now.Sub(r.lastVoiceTime) >= silenceTimeout {
				if r.OnVoiceStopped != nil {
					int16Data := unsafe.Slice((*int16)(unsafe.Pointer(&r.pCapturedSamples[0])), len(r.pCapturedSamples)/2)
					intData := make([]int, len(int16Data))
					for i, v := range int16Data {
						intData[i] = int(v)
					}
					r.OnVoiceStopped(now)
					r.voiceEventQueue <- intData
				}
				r.isRecording = false
				r.lastVoiceTime = time.Time{}
			}
		},
	})
	if err != nil {
		r.ctx.Uninit()
		return err
	}

	r.device = device
	return nil
}

func (r *Recorder) processVoiceEvents() {
	// Continuously process events from the queue
	for wavData := range r.voiceEventQueue {
		if r.OnVoiceEvent != nil {
			r.OnVoiceEvent(wavData)
		}
	}
}

func (r *Recorder) Start() error {
	return r.device.Start()
}

func (r *Recorder) Deinit() {
	if r.device != nil {
		r.device.Uninit()
	}
	r.ctx.Uninit()
	r.ctx.Free()
}

func (r *Recorder) GetSampleRate() int {
	return int(r.config.SampleRate)
}

func (r *Recorder) GetChannels() int {
	return int(r.config.Capture.Channels)
}

func (r *Recorder) GetBitDepth() int {
	return 16
}

func (r *Recorder) ToWav(wavData []int, sampleRate, bitDepth, channels int) ([]byte, error) {
	buf := &bytes.Buffer{}
	ws := &nopSeek{buf}

	encoder := wav.NewEncoder(ws, sampleRate, bitDepth, channels, 1)
	defer encoder.Close()

	err := encoder.Write(&audio.IntBuffer{
		Format:         &audio.Format{NumChannels: channels, SampleRate: sampleRate},
		Data:           wavData,
		SourceBitDepth: bitDepth,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
