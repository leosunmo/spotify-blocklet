package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/godbus/dbus"
	"github.com/leosunmo/go-spotify-dbus"
)

const (
	leftClick   = '1'
	rightClick  = '2'
	middleClick = '3'
	scrollUp    = '4'
	scrollDown  = '5'
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Fatal(err)
	}

	chanMeta := make(chan *spotify.Metadata)
	chanPs := make(chan spotify.PlaybackStatus)
	chanSvc := make(chan bool)

	listeners := &spotify.Listeners{
		OnMetadata: func(metadata *spotify.Metadata) {
			chanMeta <- metadata
		},
		OnPlaybackStatus: func(status spotify.PlaybackStatus) {
			chanPs <- status
		},
		OnServiceStart: func() {
			// Spotify is already running or has just been started
			chanSvc <- true
		},
		OnServiceStop: func() {
			// Spotify is not running
			chanSvc <- false
		},
	}

	// Start listening to the dbus
	go spotify.Listen(conn, listeners)

	// Start stdout drawing
	go drawStdout(chanSvc, chanPs, chanMeta)

	reader := bufio.NewReader(os.Stdin)

	for {
		button, _, err := reader.ReadRune()
		if err != nil {
			log.Printf("while reading stdin rune, err: %s", err.Error())
		}
		switch button {
		case leftClick:
			spotify.SendPlayPause(conn)
		case middleClick:
			spotify.SendNext(conn)
		case rightClick:
			spotify.SendPrevious(conn)
		}
	}
}

func drawStdout(svc chan bool, ps chan spotify.PlaybackStatus, metadata chan *spotify.Metadata) {
	var artist, song string
	var running, playing bool
	for {
		select {
		case running = <-svc:
			if !running {
				fmt.Println("<span foreground=\"#cccc00\" size=\"smaller\">not running</span>")
				continue
			}
		case cp := <-ps:
			if cp == "Playing" {
				playing = true
			} else {
				playing = false
			}
		case meta := <-metadata:
			artist = strings.Join(meta.Artist, ", ")
			song = meta.Title
		}
		time.Sleep(1)
		if !playing {
			fmt.Printf("<span foreground=\"#cccc00\" size=\"smaller\">%s - %s</span>\n", artist, song)
		} else {
			fmt.Printf("<small>%s - %s</small>\n", artist, song)
		}
	}
}