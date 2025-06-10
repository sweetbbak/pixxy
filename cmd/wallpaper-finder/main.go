package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Extensions     []string `short:"e" long:"extension" description:"extensions to search for ie (-e png -e jpg) can specify multliple times"`
	Ratio          string   `short:"r" long:"ratio" description:"a ratio to search, in the format <widht>x<height> (16x9)"`
	Tolerance      float32  `short:"t" long:"tolerance" description:"percentage of tolerance for the ratio [5%]"`
	Color          bool     `short:"c" long:"color" description:"print paths with color"`
	FollowSymlinks bool     `short:"f" long:"follow" description:"follow symlinks"`
	Verbose        bool     `short:"v" long:"verbose" description:"print debugging information and verbose output"`
}

type Aspect struct {
	Width            float32
	Height           float32
	Ratio            float32
	upper_tolerance  float32
	lower_tolerance  float32
	calculated_ratio float32
}

var ASPECT_RATIO float32
var CURRENT_PATH string

func (self *Aspect) isRatio() bool {
	ratio := self.Width / self.Height
	lower := ratio * self.lower_tolerance
	upper := ratio * self.upper_tolerance
	self.calculated_ratio = ratio
	return lower <= self.Ratio && self.Ratio <= upper
}

// log function
var debug = func(string, ...any) {}

var DecodeError = errors.New("error decoding image")

// image types
// jpg(jpeg), png, gif, tif(tiff), bmp, webp

func matchExtension(filename, fileExt string) bool {
	real_ext := path.Ext(filename)
	if real_ext == "" {
		return false
	}

	if real_ext[1:] == fileExt[1:] {
		return true
	}

	return false
}

func walkFunc2(path string, info fs.DirEntry, err error) error {
	// path = filepath.Join(CURRENT_PATH, path)
	// debug("base_path: %v", path)

	if !info.Type().IsRegular() {
		// if !info.Mode().IsRegular() {
		return nil
	}

	var matched bool
	for _, ext := range opts.Extensions {
		if matchExtension(path, ext) {
			matched = true
			break
		}
	}

	if !matched {
		return nil
	}

	// debug("match %v", path)

	file, err := os.Open(path)
	if err != nil {
		// debug("open err: %v", path)
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		// debug("decode err: %v - format %v", path, format)
		return err
		// return fmt.Errorf("error decoding image: %s - %w", path, err)
	}

	aspect := &Aspect{
		Width:           float32(img.Width),
		Height:          float32(img.Height),
		Ratio:           ASPECT_RATIO,
		upper_tolerance: (100 + opts.Tolerance) / 100,
		lower_tolerance: (100 - opts.Tolerance) / 100,
	}

	if aspect.isRatio() {
		// fmt.Fprintf(os.Stdout, "%s [\x1b[32m %.2f \x1b[0m]\n", path, aspect.calculated_ratio)
		fmt.Fprintf(os.Stdout, "%s\n", path)
	}

	return nil
}

func walkFunc(path string, info os.FileInfo, err error) error {
	path = filepath.Join(CURRENT_PATH, path)
	// debug("base_path: %v", path)

	// if !info.Type().IsRegular() {
	if !info.Mode().IsRegular() {
		return nil
	}

	var matched bool
	for _, ext := range opts.Extensions {
		if matchExtension(path, ext) {
			matched = true
			break
		}
	}

	if !matched {
		return nil
	}

	debug("match %v", path)

	file, err := os.Open(path)
	if err != nil {
		// debug("open err: %v", path)
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		// debug("decode err: %v - format %v", path, format)
		// return fmt.Errorf("error decoding image: %s - %w", path, err)
		return DecodeError
	}

	aspect := &Aspect{
		Width:           float32(img.Width),
		Height:          float32(img.Height),
		Ratio:           ASPECT_RATIO,
		upper_tolerance: (100 + opts.Tolerance) / 100,
		lower_tolerance: (100 - opts.Tolerance) / 100,
	}

	if aspect.isRatio() {
		// fmt.Fprintf(os.Stdout, "%s [\x1b[32m %.2f \x1b[0m]\n", path, aspect.calculated_ratio)
		fmt.Fprintf(os.Stdout, "%s\n", path)
	}

	return nil
}

func parseRatio(r string) (float32, error) {
	rsplit := strings.Split(r, "x")
	if len(rsplit) != 2 {
		return -1.0, fmt.Errorf("couldn't parse ratio. Format must by '<int_width>x<int_height>' ex: 16x9")
	}
	x, err := strconv.Atoi(rsplit[0])
	if err != nil {
		return -1.0, fmt.Errorf("couldn't parse ratio. Format must by '<int_width>x<int_height>' ex: 16x9")
	}
	y, err := strconv.Atoi(rsplit[1])
	if err != nil {
		return -1.0, fmt.Errorf("couldn't parse ratio. Format must by '<int_width>x<int_height>' ex: 16x9")
	}
	return float32(x) / float32(y), nil
}

func Wall(args []string) error {
	// var exitErrors []error
	for _, p := range args {
		CURRENT_PATH = p
		Walk(p, walkFunc)
	}

	return nil
}

// set default supported extensions
func SetExtensions() {
	if len(opts.Extensions) < 1 {
		opts.Extensions = append(opts.Extensions, ".png", ".jpg", ".jpeg")
	}
}

// add "." to extensions if they don't have them
func fixExtensions() {
	for i, ext := range opts.Extensions {
		if ext[0] != '.' {
			opts.Extensions[i] = fmt.Sprintf(".%s", ext)
		}
	}
}

func init() {
	opts.Tolerance = 5.0
}

var msg string = `Examples:
  wallpaper-finder ~/Pictures ~/hdd/Pictures
  wallpaper-finder -r 16x9 ~/Pictures
  wallpaper-finder -e png ~/Pictures
`

func main() {
	args, err := flags.Parse(&opts)
	if flags.WroteHelp(err) {
		println(msg)
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}

	SetExtensions() // if extensions arent set, set the defaults
	fixExtensions() // if extenstions dont have a dot fix them so they work with the Go stdlib

	if opts.Ratio == "" {
		ASPECT_RATIO = 1.7777777
	} else {
		r, err := parseRatio(opts.Ratio)
		if err != nil {
			log.Fatal(err)
		}
		ASPECT_RATIO = r
	}

	if opts.Verbose {
		debug = log.Printf
	}

	debug("OPTARGS: %v - %v", opts, args)
	debug("aspect_ratio: %v", ASPECT_RATIO)

	if err := Wall(args); err != nil {
		log.Fatal("wall: ", err)
	}
}
