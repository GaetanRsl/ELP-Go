package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)


type lineRange struct{
	from int
	to int
}
var ymax, xmax int

func main() {
	// Get file

	file, err := os.Open("ouibus.png") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	b := img.Bounds()
	imgGray := image.NewGray(b)

	ymax = b.Max.Y
	xmax = b.Max.X

	nbGoroutine := ymax/200

	var inputChannel chan lineRange
	var outputChannel chan string
	inputChannel = make (chan lineRange, nbGoroutine+1)
	outputChannel = make (chan string, nbGoroutine+1)

	for goroutine:=0 ; goroutine<nbGoroutine+1 ; goroutine++ {
		fmt.Println(goroutine)
		go RGBtoGray(inputChannel, outputChannel, img, imgGray)
	}

	pushnum := 0
	for mcpt:= 0; mcpt < xmax ; mcpt+= 200{
		pushnum ++
		toPush := lineRange{from: mcpt, to: mcpt+199}
		inputChannel <- toPush
		if (mcpt == nbGoroutine*200){
			toPush := lineRange{from: mcpt, to: mcpt+ymax%200}
			inputChannel <- toPush
		}
	}

	for rescpt := 0; rescpt < pushnum; rescpt ++{
		<- outputChannel
	}

	outFile, err := os.Create("changed.png")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()
	png.Encode(outFile, imgGray)

}

func RGBtoGray(inp chan lineRange, feedback chan string, img image.Image, imgGray *image.Gray ) {
	for{
		rng := <-inp
		for i:= rng.from ; i<rng.to ; i++{
			for j:=0 ; j<xmax ; j++ {
				RGBApx := img.At(i,j)
				r, g, b, _:= RGBApx.RGBA()
				gray := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
				grayPx := color.Gray{uint8(gray / 256)}
				imgGray.Set(i, j, grayPx)

			}
		}
		feedback <- "FINI"
	}
}
