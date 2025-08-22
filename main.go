package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/disintegration/imaging"
	"github.com/nickalie/go-webpbin"
)

var (
	quality            uint
	width              uint
	height             uint
	outDir             string
	safeFilename       bool
	withDimensions     bool
	lowercaseFilenames bool
	cwebpPath          string
)

var (
	whitespaceFilenameRegex = regexp.MustCompile(`\s+`)
	unsafeFilenameRegex     = regexp.MustCompile(`[^\w_-]+`)
)

func main() {
	flag.UintVar(&width, "width", 0, "Max width")
	flag.UintVar(&height, "height", 0, "Max height")
	flag.UintVar(&quality, "q", 90, "Quality to encode at")

	flag.StringVar(&outDir, "out", "", "Directory to put the webp file (defaults to same directory as input)")
	flag.BoolVar(&safeFilename, "safe", true, "Set to ensure the output filename is safe")
	flag.BoolVar(&withDimensions, "with-dimensions", true, "Set to include resized dimensions on output filename")
	flag.BoolVar(&lowercaseFilenames, "lowercase-filenames", false, "Set to make all filenames lowercase")

	flag.StringVar(&cwebpPath, "cwebp-path", "/opt/homebrew/bin", "Path to cwebp binary")

	flag.Parse()

	allPaths, err := getPossiblePaths(flag.Args())
	if err != nil {
		log.Panicf("get all paths: %s", err)
	}

	for _, path := range allPaths {
		out, err := processImage(path)
		if err != nil {
			log.Printf("couldn't process %s: %s", path, err.Error())
			continue
		}

		log.Printf("converted %s to %s", path, out)
	}
}

func getPossiblePaths(patterns []string) ([]string, error) {
	tbr := []string{}

	for _, pattern := range patterns {
		path, err := filepath.Abs(pattern)
		if err != nil {
			return nil, fmt.Errorf("malformed pattern: %s: %w", pattern, err)
		}

		paths, err := filepath.Glob(path)
		if err != nil {
			return nil, fmt.Errorf("glob: %w", err)
		}

		tbr = append(tbr, paths...)
	}

	return tbr, nil
}

func processImage(path string) (string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("stat: %w", err)
	}

	if stat.IsDir() {
		return "", fmt.Errorf("path is a directory")
	}

	by, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("os: read file: %w", err)
	}

	img, err := imaging.Decode(bytes.NewReader(by))
	if err != nil {
		return "", fmt.Errorf("imaging: decode: %w", err)
	}

	if width > 0 || height > 0 {
		img = fitImage(img, width, height)
	}

	outPath := outputFilename(path, outDir)

	wr, err := os.OpenFile(outPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		log.Panic("open file for writing:", err)
	}

	if err := webpbin.
		NewCWebP(webpbin.SetSkipDownload(true), webpbin.SetVendorPath(cwebpPath)).
		Quality(quality).
		InputImage(img).
		Output(wr).
		Run(); err != nil {
		return outPath, fmt.Errorf("cwebp: %w", err)
	}

	return outPath, nil
}

func outputFilename(inPath, outDir string) string {
	if outDir == "" {
		outDir = filepath.Dir(inPath)
	}

	outName := strings.Replace(filepath.Base(inPath), filepath.Ext(inPath), "", 1)

	if safeFilename {
		outName = whitespaceFilenameRegex.ReplaceAllString(outName, "-")
		outName = unsafeFilenameRegex.ReplaceAllString(outName, "")
	}

	if lowercaseFilenames {
		outName = strings.ToLower(outName)
	}

	if withDimensions && (width > 0 || height > 0) {
		outName = outName + fmt.Sprintf("-%s", dimensions(width, height))
	}

	return filepath.Join(outDir, outName+".webp")
}

func fitImage(img image.Image, width, height uint) image.Image {
	if width == 0 && height == 0 {
		return img
	}

	if width == 0 {
		width = 9999999
	} else if height == 0 {
		height = 9999999
	}

	return imaging.Fit(img, int(width), int(height), imaging.Lanczos)
}

func dimensions(width uint, height uint) string {
	if width == 0 && height == 0 {
		return ""
	}

	if width > 0 && height == 0 {
		return fmt.Sprintf("%dx", width)
	}

	if width == 0 && height > 0 {
		return fmt.Sprintf("x%d", height)
	}

	return fmt.Sprintf("%dx%d", width, height)
}
