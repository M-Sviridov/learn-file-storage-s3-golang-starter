package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"os/exec"
)

type FFProbeResp struct {
	Streams []struct {
		Width  int `json:"width,omitempty"`
		Height int `json:"height,omitempty"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	output := &bytes.Buffer{}
	cmd.Stdout = output
	if err := cmd.Run(); err != nil {
		return "", err
	}

	var ffProbeResp FFProbeResp
	if err := json.Unmarshal(output.Bytes(), &ffProbeResp); err != nil {
		return "", err
	}
	if len(ffProbeResp.Streams) == 0 {
		return "", errors.New("0 streams found when running ffprobe")
	}

	width := ffProbeResp.Streams[0].Width
	height := ffProbeResp.Streams[0].Height

	if width == 0 || height == 0 {
		return "other", nil
	}
	ratio := float64(width) / float64(height)

	const (
		landscape = 16.0 / 9.0
		portrait  = 9.0 / 16.0
		tolerance = 0.01
	)

	if math.Abs(ratio-landscape) < tolerance {
		return "16:9", nil
	}

	if math.Abs(ratio-portrait) < tolerance {
		return "9:16", nil
	}

	return "other", nil
}
