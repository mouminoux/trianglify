package main

import (
	"fmt"
	"image"
	imgColor "image/color"
	"image/draw"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/mouminoux/trianglify/tools"

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

	resp, _ := http.Get("https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fi.etsystatic.com%2F9142131%2Fc%2F1310%2F1042%2F549%2F1036%2Fil%2Fd5ea90%2F5319327925%2Fil_680x540.5319327925_ndpw.jpg&f=1&nofb=1&ipt=92e6ddf0d11144019ac33501298b6a9abdf6056150df0f83a6609765df3b69e9")
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
	iteration := 9000
	triangleMaxSize := 100
	triangleMinSize := 10
	triangleCurrentSize := triangleMaxSize
	triangleOpacity := uint8(255) // 0-255, where 255 is fully opaque

	numberIterationBeforeGivingUp := 0

	for {
		_, triangleCenter := tools.GetRandomTriangle(src.Image().Bounds(), triangleCurrentSize)

		var color imgColor.RGBA
		if randomColor {
			color = imgColor.RGBA{
				R: uint8(tools.RandomInt(0, 255)),
				G: uint8(tools.RandomInt(0, 255)),
				B: uint8(tools.RandomInt(0, 255)),
				A: triangleOpacity,
			}
		} else {
			r, g, b, _ := src.Image().At(int(triangleCenter.X), int(triangleCenter.Y)).RGBA()
			color = imgColor.RGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: triangleOpacity,
			}
		}

		//shapeBound := func() image.Rectangle { return tools.TriangleBound(src.Image().Bounds(), triangle) }
		//pointInTriangle := func(point image.Point) bool { return tools.PointInTriangle(point, triangle) }
		shapeBound := func() image.Rectangle {
			return image.Rectangle{
				Min: image.Point{triangleCenter.X - triangleCurrentSize, triangleCenter.Y - triangleCurrentSize},
				Max: image.Point{triangleCenter.X + triangleCurrentSize, triangleCenter.Y + triangleCurrentSize},
			}
		}
		pointInTriangle := func(point image.Point) bool { return true }

		scoreBefore := tools.Score(src, dst, shapeBound, pointInTriangle)

		backDst := image.NewRGBA(dst.Image().Bounds())
		draw.Draw(backDst, backDst.Bounds(), dst.Image(), image.Point{0, 0}, draw.Src)

		//tools.DrawPolygon(dst, []image.Point{triangle[0], triangle[1], triangle[2]}, color)
		dst.SetColor(color)
		dst.DrawCircle(float64(triangleCenter.X), float64(triangleCenter.Y), float64(triangleCurrentSize))
		dst.Fill()

		scoreAfter := tools.Score(src, dst, shapeBound, pointInTriangle)

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
			if numberIterationBeforeGivingUp >= 10 {
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
