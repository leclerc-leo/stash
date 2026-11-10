package transcoder

import (
	"reflect"
	"testing"
)

func TestScreenshotTimeDefaultUsesFastSeek(t *testing.T) {
	options := ScreenshotOptions{
		OutputPath: "out.jpg",
		OutputType: ScreenshotOutputTypeImage2,
	}

	got := ScreenshotTime("input.webm", 12.5, options)
	want := []string{
		"-v", "error",
		"-y",
		"-ss", "12.5",
		"-i", "input.webm",
		"-frames:v", "1",
		"-f", "image2",
		"out.jpg",
	}

	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("ScreenshotTime() = %#v, want %#v", []string(got), want)
	}
}

func TestScreenshotTimeSlowSeek(t *testing.T) {
	options := ScreenshotOptions{
		OutputPath: "out.jpg",
		OutputType: ScreenshotOutputTypeImage2,
		SlowSeek:   true,
	}

	got := ScreenshotTime("input.webm", 12.5, options)
	want := []string{
		"-v", "error",
		"-y",
		"-i", "input.webm",
		"-ss", "12.5",
		"-frames:v", "1",
		"-f", "image2",
		"out.jpg",
	}

	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("ScreenshotTime() = %#v, want %#v", []string(got), want)
	}
}

func TestScreenshotTimeExtraInputArgs(t *testing.T) {
	options := ScreenshotOptions{
		OutputPath:     "-",
		OutputType:     ScreenshotOutputTypeBMP,
		Width:          160,
		ExtraInputArgs: []string{"-hwaccel", "cuda"},
	}

	got := ScreenshotTime("input.mp4", 12.5, options)
	want := []string{
		"-v", "error",
		"-y",
		"-hwaccel", "cuda",
		"-ss", "12.5",
		"-i", "input.mp4",
		"-frames:v", "1",
		"-vf", "scale=160:-2",
		"-c:v", "bmp",
		"-f", "rawvideo",
		"-",
	}

	if !reflect.DeepEqual([]string(got), want) {
		t.Fatalf("ScreenshotTime() = %#v, want %#v", []string(got), want)
	}
}

func TestScreenshotTimeCudaFrames(t *testing.T) {
	cuda := []string{"-hwaccel", "cuda", "-hwaccel_output_format", "cuda"}
	tests := []struct {
		name   string
		width  int
		height int
		want   string
	}{
		{"width", 2048, 0, "scale_cuda=w='max(trunc(iw/4)*2,2048)':h=-2:interp_algo=bilinear," +
			"scale_cuda=w='max(trunc(iw/4)*2,2048)':h=-2:interp_algo=bilinear," +
			"scale_cuda=w='max(trunc(iw/4)*2,2048)':h=-2:interp_algo=bilinear," +
			"scale_cuda=w=2048:h=-2:format=yuv420p,hwdownload,format=yuv420p"},
		{"height", 0, 4096, "scale_cuda=w=-2:h='max(trunc(ih/4)*2,4096)':interp_algo=bilinear," +
			"scale_cuda=w=-2:h='max(trunc(ih/4)*2,4096)':interp_algo=bilinear," +
			"scale_cuda=w=-2:h=4096:format=yuv420p,hwdownload,format=yuv420p"},
		{"source size", 0, 0, "scale_cuda=format=yuv420p,hwdownload,format=yuv420p"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := ScreenshotOptions{
				OutputPath:     "out.jpg",
				OutputType:     ScreenshotOutputTypeImage2,
				Width:          tt.width,
				Height:         tt.height,
				ExtraInputArgs: cuda,
			}

			got := ScreenshotTime("input.mp4", 1, options)
			want := []string{
				"-v", "error",
				"-y",
				"-hwaccel", "cuda", "-hwaccel_output_format", "cuda",
				"-ss", "1",
				"-i", "input.mp4",
				"-frames:v", "1",
				"-vf", tt.want,
				"-f", "image2",
				"out.jpg",
			}

			if !reflect.DeepEqual([]string(got), want) {
				t.Fatalf("ScreenshotTime() = %#v, want %#v", []string(got), want)
			}
		})
	}
}
