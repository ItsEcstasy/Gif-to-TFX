package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

var (
	ClearTFX      = true
	SleepTFX      = true
	clearReplacer = "<<clear>>"
	sleepReplacer = "<<85>>"
	showProgress  = false
	showPreview   = false
)

func main() {
	var gifName string

	fmt.Print("File: ")
	_, err := fmt.Scan(&gifName)
	if err != nil {
		log.Fatal(err)
	}

	Gif, err := os.Open(gifName)
	if err != nil {
		log.Fatal(err)
	}
	defer Gif.Close()

	imageKek, err := gif.DecodeAll(Gif)
	if err != nil {
		log.Fatal(err)
	}

	bar := " "
	for i := 0; i < 20; i++ {
		bar += " "
	}

	roundProgress := 0
	cutName := strings.TrimSuffix(gifName, ".gif")
	outputName := fmt.Sprintf("%s.txt", cutName)
	outputFile, err := os.Create(outputName)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	for index, frame := range imageKek.Image {
		width := 80
		height := 25
		resized := resize.Resize(uint(width), uint(height), frame, resize.Lanczos3)

		gifOutput := imageToText(resized)
		showOutput := true
		if showOutput {
			fmt.Print("\f")
			fmt.Print("Converting...\n")
			if showPreview {
				fmt.Print(gifOutput)
			}
		}

		text := gifOutput

		if SleepTFX {
			text += sleepReplacer + "\n"
		}

		if ClearTFX {
			text += clearReplacer + " \n"
		}

		if text == sleepReplacer {
			continue
		}

		if showProgress {
			Progress := float64(index) / float64(len(imageKek.Image))
			UpdateProgress := int(float64(20) * Progress)
			if roundProgress != UpdateProgress {
				bar = "[" + strings.Repeat("#", UpdateProgress) + strings.Repeat(" ", 20-UpdateProgress) + "]\r\n"
				fmt.Print("\033[H" + bar)
				roundProgress = UpdateProgress
			}
		}

		_, err := outputFile.WriteString(text)
		if err != nil {
			log.Fatal(err)
		}
	}

	if showProgress {
		fmt.Print("\n")
	}
}

func imageToText(imageKek image.Image) string {
	palette := map[float64]string{
		64:  "█", // very bright
		128: "▓", // bright
		192: "▒", // dark
		256: "░", // very dark
	}

	text := ""
	text += "\033[0m"

	for y := 0; y < imageKek.Bounds().Dy(); y++ {
		for x := 0; x < imageKek.Bounds().Dx(); x++ {
			r, g, b, _ := imageKek.At(x, y).RGBA()
			rgba := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255}
			character := ""
			for k, v := range palette {
				if Brightness(rgba) < k {
					character = v
					break
				}
			}

			text += RGBEscape(rgba) + character
		}
		text += "\033[0m\n"
	}

	return text
}

func RGBEscape(c color.RGBA) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm", c.R, c.G, c.B, c.R, c.G, c.B)
}

func Brightness(c color.RGBA) float64 {
	return 0.2126*float64(c.R) + 0.7152*float64(c.G) + 0.0722*float64(c.B)
}
