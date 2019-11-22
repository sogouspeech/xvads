package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/sogouspeech/xvads/pkg/webrtcvad"
)

var (
	FrameDuration    = time.Millisecond * 20
	SilenceThreshold = time.Millisecond * 200
	ActiveThreshold  = time.Second * 30
)

func init() {
	if d, _ := time.ParseDuration(os.Getenv("FRAME_DURATION")); d > 0 {
		FrameDuration = d
	}
	if d, _ := time.ParseDuration(os.Getenv("SILENCE_THRESHOLD")); d > 0 {
		SilenceThreshold = d
	}
	if d, _ := time.ParseDuration(os.Getenv("ACTIVE_THRESHOLD")); d > 0 {
		ActiveThreshold = d
	}
}

func main() {
	chreader := split(os.Stdin)
	for r := range chreader {
		forward(r)
	}
}

func split(source io.Reader) <-chan io.Reader {
	ch := make(chan io.Reader, 2)
	go func() {
		defer close(ch)

		vad, err := webrtcvad.New()
		if err != nil {
			log.Fatalln(err)
		}

		_ = vad.SetMode(1)

		// 16000 采样率下，一帧的数据量
		buf := make([]byte, FrameDuration/time.Millisecond*32)

		var activeWriter io.WriteCloser
		silentDuration := time.Duration(0)
		activeDuration := time.Duration(0)

		defer func() {
			if activeWriter != nil {
				_ = activeWriter.Close()
			}
		}()

		for {

			_, err := io.ReadFull(source, buf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}

			if err != nil {
				log.Fatalln(err)
				return
			}

			act, err := vad.Process(16000, buf)

			if err != nil {
				log.Fatalln(err)
				return
			}

			if activeWriter != nil {
				_, _ = activeWriter.Write(buf)

				if activeDuration += FrameDuration; activeDuration > ActiveThreshold {
					_ = activeWriter.Close()
					activeWriter = nil
					silentDuration = 0
					continue
				}

				if !act {
					if silentDuration += FrameDuration; silentDuration >= SilenceThreshold {
						_ = activeWriter.Close()
						activeWriter = nil
						silentDuration = 0
					}
				} else {
					silentDuration = 0
				}
			} else {
				if act {
					r, w := io.Pipe()
					ch <- r
					_, _ = w.Write(buf)
					activeWriter = w
					activeDuration = FrameDuration
				}
			}
		}

	}()

	return ch
}

func forward(r io.Reader) {

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}
