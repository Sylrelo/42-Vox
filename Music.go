package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

type SoundData struct {
	channel     chan float64
	done        chan bool
	volumeValue float64
	volume      *effects.Volume
}

type SoundEngine struct {
	desertBiomeMusic SoundData
	snowBiomeMusic   SoundData
	ambiantMusic     SoundData
}

func loadMusic(path string) (beep.StreamSeekCloser, beep.Format) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return streamer, format
}

func prepareMusic(streamer beep.StreamSeekCloser, format beep.Format) *effects.Volume {
	buffer2 := beep.NewBuffer(format)
	buffer2.Append(streamer)
	streamer.Close()
	musicc := buffer2.Streamer(0, buffer2.Len())
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, musicc), Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	return volume
}

func MusicHandler(soundEngine SoundEngine) {
	soundEngine.desertBiomeMusic.done = make(chan bool)
	soundEngine.snowBiomeMusic.done = make(chan bool)
	soundEngine.ambiantMusic.done = make(chan bool)

	speaker.Play(beep.Seq(soundEngine.ambiantMusic.volume, beep.Callback(func() {
		soundEngine.ambiantMusic.done <- true
	})))
	speaker.Play(beep.Seq(soundEngine.desertBiomeMusic.volume, beep.Callback(func() {
		soundEngine.desertBiomeMusic.done <- true
	})))
	speaker.Play(beep.Seq(soundEngine.snowBiomeMusic.volume, beep.Callback(func() {
		soundEngine.snowBiomeMusic.done <- true
	})))

	for {
		speaker.Lock()
		soundEngine.ambiantMusic.volume.Volume = soundEngine.ambiantMusic.volumeValue
		soundEngine.desertBiomeMusic.volume.Volume = soundEngine.desertBiomeMusic.volumeValue
		soundEngine.snowBiomeMusic.volume.Volume = soundEngine.snowBiomeMusic.volumeValue
		speaker.Unlock()
		soundEngine.desertBiomeMusic.volumeValue = <-soundEngine.desertBiomeMusic.channel
		soundEngine.snowBiomeMusic.volumeValue = <-soundEngine.snowBiomeMusic.channel
	}
}

func InitMusics(soundEngine *SoundEngine) {
	soundEngine.ambiantMusic.volumeValue = -2
	soundEngine.desertBiomeMusic.volumeValue = -7.5
	soundEngine.snowBiomeMusic.volumeValue = -7.5
	soundEngine.desertBiomeMusic.channel = make(chan float64)
	soundEngine.snowBiomeMusic.channel = make(chan float64)

	fmt.Println("\x1b[94mLoading sounds...\x1b[0m")
	streamerAmbient, formatAmbient := loadMusic("Assets/music/ambient.mp3")
	streamerDesert, formatDesert := loadMusic("Assets/music/desert.mp3")
	streamerSnow, formatSnow := loadMusic("Assets/music/snow.mp3")

	speaker.Init(formatAmbient.SampleRate, formatAmbient.SampleRate.N(time.Second/10))

	soundEngine.ambiantMusic.volume = prepareMusic(streamerAmbient, formatAmbient)
	soundEngine.desertBiomeMusic.volume = prepareMusic(streamerDesert, formatDesert)
	soundEngine.snowBiomeMusic.volume = prepareMusic(streamerSnow, formatSnow)
	fmt.Println("\x1b[92mSounds successfully loaded.\x1b[0m")

	go MusicHandler(*soundEngine)
}

func handleMusic(vox Vox, soundEngine SoundEngine) {
	DesertSnowNoise := Noise2dSimplex(float64(vox.pos[0]), float64(vox.pos[2]), 0.0, 0.85, 0.00025, 5, 2)

	sDesertVol := (DesertSnowNoise) * 4
	sDesertVol = math.Pow(sDesertVol, 10)
	sDesertVol *= -1
	sDesertVol--
	sDesertVol *= 0.5
	soundEngine.desertBiomeMusic.channel <- sDesertVol

	sSnowVol := (1 - DesertSnowNoise) * 4
	sSnowVol = math.Pow(sSnowVol, 10)
	sSnowVol *= -1
	sSnowVol--
	sSnowVol *= 0.5
	soundEngine.snowBiomeMusic.channel <- sSnowVol
}
