package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

type FFProbeResp struct {
	Streams []struct {
		Width  int `json:"width,omitempty"`
		Height int `json:"height,omitempty"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_streams",
		filePath,
	)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error while running ffprobe: %v", err)
	}

	var ffProbeResp FFProbeResp
	if err := json.Unmarshal(stdout.Bytes(), &ffProbeResp); err != nil {
		return "", fmt.Errorf("couldn't parse ffprobe output data: %v", err)
	}

	if len(ffProbeResp.Streams) == 0 {
		return "", errors.New("no video streams found")
	}

	width := ffProbeResp.Streams[0].Width
	height := ffProbeResp.Streams[0].Height

	if width == 16*height/9 {
		return "16:9", nil
	} else if height == 16*width/9 {
		return "9:16", nil
	}
	return "other", nil
}

func processVideoForFastStart(filePath string) (string, error) {
	output := filePath + ".processing"
	cmd := exec.Command(
		"ffmpeg",
		"-i", filePath,
		"-c", "copy",
		"-movflags", "faststart",
		"-f", "mp4",
		output,
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error while running ffmpeg: %v", err)
	}

	return output, nil
}
