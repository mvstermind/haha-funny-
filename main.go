package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

var (
	otoCtx *oto.Context
	wg     sync.WaitGroup
	rng    = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func main() {
	var err error
	op := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}
	otoCtx, _, err = oto.NewContext(op)
	if err != nil {
		panic("Failed to initialize audio context: " + err.Error())
	}

	err = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.CtrlC {
			fmt.Println("C ya ltr")
			return true, nil
		}
		fmt.Println("Key pressed:", key.Code.String())

		wg.Add(1)
		go choseRandomSound()
		return false, nil
	})

	if err != nil {
		panic("Keyboard listening error: " + err.Error())
	}

	wg.Wait()
}

func choseRandomSound() {
	defer wg.Done()

	randNum := rng.Intn(12) + 1
	fileName := fmt.Sprintf("./goofy-sounds/%v.mp3", randNum)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decodedMp3, err := mp3.NewDecoder(file)
	if err != nil {
		fmt.Println("Error decoding MP3:", err)
		return
	}

	player := otoCtx.NewPlayer(decodedMp3)
	defer player.Close()

	player.Play()
	fmt.Println("Now playing d-_-b: ", fileName)

	for player.IsPlaying() {
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("Stopped playing:( : ", fileName)
}
