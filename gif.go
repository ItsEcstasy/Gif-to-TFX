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

// [+] ===================================================================== [+]
var ClearTFX = false // toggle clear replacer if you require one
var SleepTFX = true  // toggle sleep replacer if you require one

var clearReplacer = "<<clear>>" // replace with the clear replacer you use
var sleepReplacer = "<<85>>"    // replace with the sleep replacer you use

var showProgress = false // toggle the display the progress bar
var showPreview = false  // toggle the preview, note: this will take longer btw
// [+] ===================================================================== [+]

func main() {
	var gifName string
	fmt.Print("Enter GIF File Name Example: Eye.gif\nName: ")
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

	outputName := fmt.Sprintf("%s-output.txt", gifName)
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
			fmt.Print("\033[H\033[2J")
			fmt.Print("Converting...\n")
			if showPreview {
				fmt.Print(gifOutput)
			}
		}

		text := gifOutput + "\r\n"

		if ClearTFX {
			text += clearReplacer + "\r\n"
		}
		if SleepTFX {
			text += sleepReplacer + "\r\n"
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
		64:  "█", // very bright areas
		128: "▓", // bright areas
		192: "▒", // dark areas
		256: "░", // very dark areas
	}

	text := ""
	text += "\033[0m" // Reset text color

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
