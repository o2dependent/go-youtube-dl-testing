package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
)

type QualityInfo struct {
	Quality      string
	AudioQuality string
}

type Info struct {
	Author      string
	Title       string
	Duration    time.Duration
	PublishDate time.Time
	QualityInfo []QualityInfo
}

func main() {
	testingVideoUrls := [4]string{
		"https://www.youtube.com/watch?v=SpH83KzVKDc", // JPEGMAFIA - SIN MIEDO
		"https://www.youtube.com/watch?v=JtR9JkVk9aU", // I ᐸ3 Harajuku ft. Fraxiom
		"https://www.youtube.com/watch?v=IuJIbXpex_s", // 100 gecs stupid horse (Remix)
		"https://www.youtube.com/watch?v=1Bw2dTY3SsQ", // 100 gecs - mememe
	}
	for i := 0; i < len(testingVideoUrls); i++ {
		info, err := getImportantInfo(testingVideoUrls[i])
		if err != nil {
			panic(err)
		}

		fmt.Println("---- quality " + strconv.FormatInt(int64(i), 10) + " ----")
		fmt.Println(&info)
	}

}

func getImportantInfo(videoUrl string) (Info, error) {
	videoID, err := youtube.ExtractVideoID(videoUrl)
	if err != nil {
		return Info{}, err
	}

	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return Info{}, err
	}
	var qualityInfo []QualityInfo

	formats := video.Formats
	formats.Sort()

	for i := 0; i < len(video.Formats); i++ {
		f := formats[i]
		qualityInfo = append(qualityInfo, QualityInfo{
			Quality:      f.Quality,
			AudioQuality: f.AudioQuality,
		})
	}

	info := Info{
		Author:      video.Author,
		Title:       video.Title,
		Duration:    video.Duration,
		PublishDate: video.PublishDate,
		QualityInfo: qualityInfo,
	}

	return info, err
}

func download(videoUrl string, quality string, audioQuality string) {
	videoID, err := youtube.ExtractVideoID(videoUrl)
	if err != nil {
		panic(err)
	}

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

		return f.Quality == quality && strings.Contains(f.MimeType, "video/mp4")
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
	cmd := exec.Command("ffmpeg", "-i", "video.mp4", "-i", "audio.mp4", "-c:v", "copy", "-c:a", "aac", "output.mp4")
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// Delete video and audio file
	os.Remove("video.mp4")
	os.Remove("audio.mp4")
}
