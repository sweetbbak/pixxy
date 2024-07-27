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

type Dither struct {
	Halftone    bool     `short:"H" long:"halftone" description:"halftone dithering"`
	Bayer       bool     `short:"b" long:"bayer" description:"bayer dithering"`
	Scale       bool     `short:"s" long:"scale" description:"rescale image down and then up to accentuate fx"`
	Threshold   float64  `short:"t" long:"threshold" description:"float from 0.0 - 1.0"`
	Input       string   `short:"i" long:"input" description:"input image file, explicit flag (also accepts a trailing positional argument)"`
	DitherType  []string `short:"d" long:"dither" description:"dither type using error diffusion dithering"`
	ODM         []string `short:"m" long:"ordered" description:"ordered dither matrix type dithering"`
	Output      string   `short:"o" long:"output" description:"save image/gif as output file"`
	Palette     []string `short:"p" long:"palette" description:"supply a set of hex colors to apply a color dithering effect, reduces colors to the closest supplied color for each pixel"`
	PaletteFile string   `short:"P" long:"palette-file" description:"supply a set of colors from a file, uses regex to extract any valid hex color (can use messy files, like terminal theme files, json, etc...)"`
	ColorDepth  int      `short:"c" long:"color-depth" description:"create palette from the supplied image of N colors. Less is more aesthetic, more is more accurate to source."`

	Args struct {
		Image string
	} `positional-args:"yes" positional-arg-name:"TORRENT"`
}

type Filters struct {
	Halftone bool    `short:"H" long:"halftone" description:"halftone dithering"`
	Amount   float64 `short:"t" long:"threshold" description:"glitch threshold"`
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
