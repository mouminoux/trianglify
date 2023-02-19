package main

import (
	"fmt"
	"github.com/mouminoux/trianglify/tools"
	"image"
	imgColor "image/color"
	"image/draw"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/fogleman/gg"

	_ "net/http/pprof"
)

func main() {
	s := &http.Server{
		Addr: "127.0.0.1:8081",
	}
	go func() {
		s.ListenAndServe()
	}()

	rand.Seed(time.Now().UnixNano())

	resp, _ := http.Get("https://www.usinenouvelle.com/mediatheque/0/0/3/000734300_896x598_c.jpg")
	body := resp.Body
	imgSrc, _, err := image.Decode(body)
	if err != nil {
		panic(err)
	}

	imgSrc = tools.Resize(imgSrc, 1000, 600)

	// Create a new image with the same size as the source image
	imgDst := image.NewRGBA(imgSrc.Bounds())

	src := gg.NewContextForImage(imgSrc)
	dst := gg.NewContextForImage(imgDst)

	tools.DrawPolygon(dst, []image.Point{
		{
			X: dst.Image().Bounds().Min.X,
			Y: dst.Image().Bounds().Min.Y,
		}, {
			X: dst.Image().Bounds().Max.X,
			Y: dst.Image().Bounds().Min.Y,
		}, {
			X: dst.Image().Bounds().Max.X,
			Y: dst.Image().Bounds().Max.Y,
		}, {
			X: dst.Image().Bounds().Min.X,
			Y: dst.Image().Bounds().Max.Y,
		},
	}, imgColor.RGBA{
		R: 31,
		G: 31,
		B: 31,
		A: 255,
	})

	randomColor := false
	iteration := 5000
	triangleMaxSize := 100
	triangleMinSize := 20
	triangleCurrentSize := triangleMaxSize

	numberIterationBeforeGivingUp := 0

	for {
		triangle, triangleCenter := tools.GetRandomTriangle(src.Image().Bounds(), triangleCurrentSize)

		var color imgColor.RGBA
		if randomColor {
			color = imgColor.RGBA{
				R: uint8(tools.RandomInt(0, 255)),
				G: uint8(tools.RandomInt(0, 255)),
				B: uint8(tools.RandomInt(0, 255)),
				A: 255,
			}
		} else {
			r, g, b, _ := src.Image().At(int(triangleCenter.X), int(triangleCenter.Y)).RGBA()
			color = imgColor.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: 255,
			}
		}

		scoreBefore := tools.Score(src, dst, triangle)

		backDst := image.NewRGBA(dst.Image().Bounds())
		draw.Draw(backDst, backDst.Bounds(), dst.Image(), image.Point{0, 0}, draw.Src)

		tools.DrawPolygon(dst, []image.Point{triangle[0], triangle[1], triangle[2]}, color)

		scoreAfter := tools.Score(src, dst, triangle)

		if scoreAfter < scoreBefore && true {
			if iteration%100 == 0 {
				fmt.Printf("remaining: %d, size: %d, scoreBefore: %d, scoreAfter: %d\n", iteration, triangleCurrentSize, scoreBefore, scoreAfter)
				err = dst.SavePNG("output.png")
				if err != nil {
					panic(err)
				}
			}
			iteration--
			numberIterationBeforeGivingUp = 0
			triangleCurrentSize++
			triangleCurrentSize = int(math.Min(float64(triangleCurrentSize), float64(triangleMaxSize)))
			if iteration == 0 {
				break
			}
		} else {
			dst = gg.NewContextForImage(backDst)
			numberIterationBeforeGivingUp++
			if numberIterationBeforeGivingUp >= 100 {
				numberIterationBeforeGivingUp = 0
				triangleCurrentSize--
				triangleCurrentSize = int(math.Max(float64(triangleCurrentSize), float64(triangleMinSize)))
			}
		}
	}

	err = src.SavePNG("src.png")
	if err != nil {
		panic(err)
	}

	err = dst.SavePNG("output.png")
	if err != nil {
		panic(err)
	}

}
