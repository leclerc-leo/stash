package transcoder

import (
	"fmt"

	"github.com/stashapp/stash/pkg/ffmpeg"
)

type ScreenshotOptions struct {
	OutputPath string
	OutputType ScreenshotOutputType

	// Quality is the quality scale. See https://ffmpeg.org/ffmpeg.html#Main-options
	Quality int

	// Width is the width to scale the screenshot to. If 0, no scaling will be applied.
	Width int
	// Height is the height to scale the screenshot to. If 0, no scaling will be applied.
	// Not used if Width is set.
	Height int

	// Verbosity is the logging verbosity. Defaults to LogLevelError if not set.
	Verbosity ffmpeg.LogLevel

	UseSelectFilter bool

	// SlowSeek uses accurate seek by placing -ss after the input.
	SlowSeek bool

	// ExtraInputArgs are placed before the input. With "-hwaccel_output_format cuda",
	// frames stay on the GPU, so they are scaled there and then downloaded.
	// ScreenshotFrame ignores them: it only serves very short files.
	ExtraInputArgs []string
}

// hwOutputFormat returns the value of -hwaccel_output_format in args, or "".
func hwOutputFormat(args []string) string {
	for i := 0; i+1 < len(args); i++ {
		if args[i] == "-hwaccel_output_format" {
			return args[i+1]
		}
	}
	return ""
}

// screenshotScale scales to w x h (0 keeps the source size), on the GPU if the frames are on it.
func screenshotScale(extraInputArgs []string, w, h int) ffmpeg.VideoFilter {
	var vf ffmpeg.VideoFilter
	if hwOutputFormat(extraInputArgs) == "cuda" {
		size := ""
		if w > 0 || h > 0 {
			size = fmt.Sprintf("w=%d:h=%d:", w, h)
		}
		return vf.Append("scale_cuda=" + size + "format=yuv420p").Append("hwdownload").Append("format=yuv420p")
	}
	if w > 0 || h > 0 {
		return vf.ScaleDimensions(w, h)
	}
	return vf
}

func (o *ScreenshotOptions) setDefaults() {
	if o.Verbosity == "" {
		o.Verbosity = ffmpeg.LogLevelError
	}
}

type ScreenshotOutputType struct {
	codec  *ffmpeg.VideoCodec
	format ffmpeg.Format
}

func (t ScreenshotOutputType) Args() []string {
	var ret []string
	if t.codec != nil {
		ret = append(ret, t.codec.Args()...)
	}
	if t.format != "" {
		ret = append(ret, t.format.Args()...)
	}

	return ret
}

var (
	ScreenshotOutputTypeImage2 = ScreenshotOutputType{
		format: "image2",
	}
	ScreenshotOutputTypeBMP = ScreenshotOutputType{
		codec:  &ffmpeg.VideoCodecBMP,
		format: "rawvideo",
	}
)

func ScreenshotTime(input string, t float64, options ScreenshotOptions) ffmpeg.Args {
	options.setDefaults()

	var args ffmpeg.Args
	args = args.LogLevel(options.Verbosity)
	args = args.Overwrite()
	args = append(args, options.ExtraInputArgs...)

	if !options.SlowSeek {
		args = args.Seek(t)
	}
	args = args.Input(input)
	if options.SlowSeek {
		args = args.Seek(t)
	}
	args = args.VideoFrames(1)

	if options.Quality > 0 {
		args = args.FixedQualityScaleVideo(options.Quality)
	}

	w, h := 0, 0
	if options.Width > 0 {
		w, h = options.Width, -2
	} else if options.Height > 0 {
		w, h = -2, options.Height
	}

	if vf := screenshotScale(options.ExtraInputArgs, w, h); vf != "" {
		args = args.VideoFilter(vf)
	}

	args = args.AppendArgs(options.OutputType)
	args = args.Output(options.OutputPath)

	return args
}

// ScreenshotFrame uses the select filter to get a single frame from the video.
// It is very slow and should only be used for files with very small duration in secs / frame count.
func ScreenshotFrame(input string, frame int, options ScreenshotOptions) ffmpeg.Args {
	options.setDefaults()

	var args ffmpeg.Args
	args = args.LogLevel(options.Verbosity)
	args = args.Overwrite()

	args = args.Input(input)
	args = args.VideoFrames(1)

	args = args.VSync(ffmpeg.VSyncMethodPassthrough)

	var vf ffmpeg.VideoFilter
	// keep only frame number options.Frame)
	vf = vf.Select(frame)

	if options.Width > 0 {
		vf = vf.ScaleWidth(options.Width)
	}

	args = args.VideoFilter(vf)

	args = args.AppendArgs(options.OutputType)
	args = args.Output(options.OutputPath)

	return args
}
