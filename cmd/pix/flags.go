package main

import ()

type Options struct {
	Verbose bool `short:"v" long:"verbose" description:"print debugging information and verbose output"`
	Version bool `short:"V" long:"version" description:"display version info and exit"`
}

type Pixels struct {
	ParseFile string `short:"p" long:"parse" description:"file with hex colors to parse and use as a color palette. Finds any and every valid hex color in the given files text"`
	Verbose   bool   `short:"v" long:"verbose" description:"print debugging information and verbose output"`

	Args struct {
		Image string
	} `positional-args:"yes" positional-arg-name:"TORRENT"`
}

type Ascii struct {
	Input         string  `short:"i" long:"input" description:"input image file, explicit flag (also accepts a trailing positional argument)"`
	Output        string  `short:"o" long:"output" description:"save image/gif as output file"`
	Charset       string  `short:"c" long:"charset" description:"character set to use from lightest to darkest ie (.*!@#$%)"`
	CharsetPreset string  `short:"C" long:"char-preset" description:"character preset to use [limited|extended|block] "`
	FontPT        float64 `short:"p" long:"font-point" description:"font size to use"`
	Font          string  `short:"f" long:"font" description:"font to use, must be monospaced"`
	Interpolate   bool    `short:"I" long:"interpolate" description:"interpolate font so that when converting successive images (gifs) the font changes less"`
	Noise         int     `short:"n" long:"noise" description:"add random noise"`
	FFMpegArgs    string  `short:"F" long:"ffmpeg" description:"extra ffmpeg args to use when converting videos"`

	Gif   bool `short:"g" long:"gif" description:"output as gif"`
	Video bool `short:"v" long:"video" description:"process each frame of a video or gif"`

	Args struct {
		Image string
	} `positional-args:"yes" positional-arg-name:"IMAGE"`
}

type Dither struct {
	Input        string   `short:"i" long:"input" description:"input image file, explicit flag (also accepts a trailing positional argument)"`
	Output       string   `short:"o" long:"output" description:"save image/gif as output file"`
	Threshold    float64  `short:"t" long:"threshold" description:"float from 0.0 - 1.0"`
	Palette      []string `short:"p" long:"palette" description:"supply a set of hex colors to apply a color dithering effect, reduces colors to the closest supplied color for each pixel"`
	PaletteFile  string   `short:"P" long:"palette-file" description:"supply a set of colors from a file, uses regex to extract any valid hex color (can use messy files, like terminal theme files, json, etc...)"`
	ColorDepth   int      `short:"c" long:"color-depth" description:"create palette from the supplied image of N colors. Less is more aesthetic, more is more accurate to source."`
	Scale        bool     `short:"s" long:"scale" description:"rescale image down and then up to accentuate fx"`
	ScaleFactor  int      `short:"S" long:"scale-factor" description:"the amount to resize the dither effect"`
	Halftone     bool     `short:"H" long:"halftone" description:"halftone dithering"`
	Bayer        bool     `short:"b" long:"bayer" description:"bayer dithering"`
	DitherType   []string `short:"d" long:"dither" description:"dither type using error diffusion dithering"`
	ListDithers  bool     `short:"z" long:"list-dithers" description:"list ditherers"`
	ListMatrices bool     `short:"x" long:"list-maps" description:"list matrix map filters"`
	ODM          []string `short:"m" long:"ordered" description:"ordered dither matrix type dithering"`

	Args struct {
		Image string
	} `positional-args:"yes" positional-arg-name:"IMAGE"`
}

type Filters struct {
	Amount float64 `short:"t" long:"threshold" description:"glitch threshold"`
}

type Glitch struct {
	Gif         bool     `short:"g" long:"gif" description:"create a gif"`
	Verbose     bool     `short:"v" long:"verbose" description:"verbose output - show glitch steps as they occur"`
	Palette     []string `short:"p" long:"palette" description:"supply a set of hex colors to apply a color dithering effect, reduces colors to the closest supplied color for each pixel"`
	PaletteFile string   `short:"P" long:"palette-file" description:"supply a set of colors from a file, uses regex to extract any valid hex color (can use messy files, like terminal theme files, json, etc...)"`
	Seed        string   `short:"s" long:"seed" description:"random seed string"`
	Factor      float64  `short:"t" long:"threshold" description:"glitch threshold"`
	FrameDelay  int      `short:"d" long:"delay" description:"delay in between frames in milliseconds"`
	FrameCount  int      `short:"f" long:"frames" description:"amount of frames to create and glitch"`
	ColorDepth  int      `short:"c" long:"color-depth" description:"create palette from the supplied image of N colors. Less is more aesthetic, more is more accurate to source."`
	Input       string   `short:"i" long:"input" description:"input image file, explicit flag (also accepts a trailing positional argument)"`
	Output      string   `short:"o" long:"output" description:"save image/gif as output file"`

	Args struct {
		Image string
	} `positional-args:"yes" positional-arg-name:"TORRENT"`
}
