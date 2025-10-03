package sound

import (
	"bytes"
	"os"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type Sound struct {
	Options *oto.NewContextOptions
	Context *oto.Context
	Player  *oto.Player
}

func New() *Sound {
	sound := Sound{}

	// load alarm file
	fileBytes, err := os.ReadFile("./beep.mp3")
	if err != nil {
		panic("reading beep.mp3 failed: " + err.Error())
	}

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// context and options
	sound.Options = &oto.NewContextOptions{}
	sound.Options.SampleRate = 44100
	sound.Options.ChannelCount = 2
	sound.Options.Format = oto.FormatSignedInt16LE

	var readyChan chan struct{}
	sound.Context, readyChan, err = oto.NewContext(sound.Options)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	sound.Player = sound.Context.NewPlayer(decodedMp3)

	return &sound
}

func (s *Sound) Play() {
	// Play starts playing the sound and returns without waiting for it (Play() is async).
	s.Player.Play()

	// We can wait for the sound to finish playing using something like this
	for s.Player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// If you don't want the player/sound anymore simply close
	err := s.Player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
}
