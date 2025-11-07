package helpers

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os/exec"

	// Add time for speaker initialization if you use gopxl/beep/speaker
	// REPLACE "github.com/gopxl/beep/mp3"
	// Keep beep core if you use its speaker/streamer
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hajimehoshi/oto/v2"
)

var globalOpts struct {
	Server string
	Format bool
}

type StreamStoppedMsg struct{}

func StreamMusic(stationID string, ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		audioURL := "https://radio.garden/api/ara/content/listen/" + stationID + "/channel.mp3"

		// 1. Make the HTTP GET request.
		req, err := http.NewRequestWithContext(ctx, "GET", audioURL, nil)
		if err != nil {
			log.Printf("Error making request: %v", err)
			return StreamStoppedMsg{}

		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return StreamStoppedMsg{}
			}
			log.Printf("Error making request: %v", err)
			return StreamStoppedMsg{}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Received non-200 status code: %d %s", resp.StatusCode, resp.Status)
		}

		cmd := exec.Command("ffmpeg",
			"-i", "pipe:0", // Read input from stdin (pipe:0)
			"-f", "s16le",
			"-acodec", "pcm_s16le",
			"-ac", "2",
			"-ar", "44100",
			"pipe:1", // Write output to stdout (pipe:1)
		)

		stdinPipe, err := cmd.StdinPipe()
		if err != nil {
			log.Fatalf("Failed to create stdin pipe for FFmpeg: %v", err)
		}

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatalf("Failed to create stdout pipe for FFmpeg: %v", err)
		}

		// Start the FFmpeg process
		if err := cmd.Start(); err != nil {
			log.Fatalf("Failed to start FFmpeg: %v. Is it installed and in PATH?", err)
		}

		// --- 3. Start I/O operations (Concurrency) ---

		// Goroutine 1: Feed the HTTP response body into FFmpeg's stdin
		go func() {
			defer stdinPipe.Close()
			_, copyErr := io.Copy(stdinPipe, resp.Body)

			// io.Copy transfers data until EOF is reached (or an error occurs)
			if copyErr != nil && !errors.Is(copyErr, context.Canceled) {
				log.Printf("Stream copy ended unexpectedly: %v", copyErr)
			}
			stdinPipe.Close()
			stdoutPipe.Close()
		}()

		const (
			sampleRate   = 44100
			channelCount = 2
		)

		op := &oto.NewContextOptions{
			SampleRate:   sampleRate,
			ChannelCount: channelCount,
			Format:       oto.FormatSignedInt16LE,
		}

		otoCtx, ready, err := oto.NewContext(op.SampleRate, op.ChannelCount, op.Format)
		if err != nil {
			log.Fatalf("Failed to create audio context: %v", err)
		}
		<-ready // Wait for the audio context to be ready

		player := otoCtx.NewPlayer(stdoutPipe)
		player.SetVolume(0.1)
		player.Play()
		<-ctx.Done()

		player.Close()
		resp.Body.Close()

		cmd.Wait()

		return StreamStoppedMsg{}
	}
}

// Helper function to cast the int16 slice to a []byte slice as required by go-mp3's Read method.
type pcm16 []int16

func (p pcm16) Write(buf []byte) (int, error) {
	if len(buf)%2 != 0 {
		panic("buffer size must be even")
	}
	n := len(buf) / 2
	for i := 0; i < n; i++ {
		// Convert little-endian bytes to int16. This is standard for audio streams.
		p[i] = int16(buf[i*2]) | int16(buf[i*2+1])<<8
	}
	return len(buf), nil
}
