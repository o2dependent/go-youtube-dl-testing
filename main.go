package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func main() {
	videoID := "5URefVYaJrA"
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	fmt.Println(video.Author)

	if err != nil {
		panic(err)
	}

	// Video File
	fmt.Println("---- VIDEO ONLY ----")
	formats := video.Formats
	formats = formats.Select(func(f youtube.Format) bool {
		fmt.Println(f.AudioQuality)
		fmt.Println(f.AudioSampleRate)
		fmt.Println(f.Quality)
		fmt.Println(f.MimeType)

		return f.Quality == "hd1440" && strings.Contains(f.MimeType, "video/mp4")
	})
	fmt.Println(formats)

	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	file, err := os.Create("video.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}

	// Audio File
	fmt.Println("---- AUDIO ONLY ----")
	formats = video.Formats.WithAudioChannels() // only get videos with audio
	formats = formats.Select(func(f youtube.Format) bool {
		fmt.Println(f.AudioQuality)
		fmt.Println(f.AudioSampleRate)
		fmt.Println(f.Quality)
		fmt.Println(f.MimeType)

		return strings.Contains(f.MimeType, "audio/mp4")
	})

	audioStream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	defer audioStream.Close()

	audioFile, err := os.Create("audio.mp4")
	if err != nil {
		panic(err)
	}
	defer audioFile.Close()

	_, err = io.Copy(audioFile, audioStream)
	if err != nil {
		panic(err)
	}

	// Concat audio and video using ffmpeg
	cmd := exec.Command("ffmpeg", "-i video.mp4 -i audio.mp4 -map 0:v -map 1:a -codec copy -shortest out.mp4")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
