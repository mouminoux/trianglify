package tools

import (
	"errors"
	"github.com/fogleman/gg"
	"golang.org/x/image/draw"
	"image"
	imgColor "image/color"
	"math"
	"math/rand"
)

func Resize(srcImg image.Image, maxWidth int, maxHeight int) image.Image {
	bounds := srcImg.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	if width > maxWidth {
		ratio := float64(height) / float64(width)
		width, height = maxWidth, int(float64(maxWidth)*ratio)
	}
	if height > maxHeight {
		ratio := float64(width) / float64(height)
		width, height = int(float64(maxHeight)*ratio), maxHeight
	}

	dstImg := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.NearestNeighbor.Scale(dstImg, dstImg.Rect, srcImg, srcImg.Bounds(), draw.Over, nil)

	return dstImg
}

func DrawPolygon(img *gg.Context, points []image.Point, color imgColor.RGBA) {
	img.SetColor(color)

	for i, point := range points {
		x := float64(point.X)
		y := float64(point.Y)
		if i == 0 {
			img.MoveTo(x, y)
		} else {
			img.LineTo(x, y)
		}
	}
	img.Fill()
}

func Score(src *gg.Context, dst *gg.Context, triangle [3]image.Point) int64 {
	srcImg := src.Image()
	dstImg := dst.Image()

	bounds := srcImg.Bounds()

	minX := int(math.Max(float64(bounds.Min.X), math.Min(math.Min(float64(triangle[0].X), float64(triangle[1].X)), float64(triangle[2].X))))
	maxX := int(math.Min(float64(bounds.Max.X), math.Max(math.Max(float64(triangle[0].X), float64(triangle[1].X)), float64(triangle[2].X))))

	minY := int(math.Max(float64(bounds.Min.Y), math.Min(math.Min(float64(triangle[0].Y), float64(triangle[1].Y)), float64(triangle[2].Y))))
	maxY := int(math.Min(float64(bounds.Max.Y), math.Max(math.Max(float64(triangle[0].Y), float64(triangle[1].Y)), float64(triangle[2].Y))))

	var score int64 = 0

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			if pointInTriangle(image.Point{x, y}, triangle) {
				r, g, b, _ := srcImg.At(x, y).RGBA()
				r &= r >> 8
				g &= g >> 8
				b &= b >> 8
				dstR, dstG, dstB, _ := dstImg.At(x, y).RGBA()
				dstR &= dstR >> 8
				dstG &= dstG >> 8
				dstB &= dstB >> 8
				score += abs(int64(r) - int64(dstR))
				score += abs(int64(g) - int64(dstG))
				score += abs(int64(b) - int64(dstB))
				if score > 9223372036854775807/1000 {
					panic(errors.New("int64 overflow"))
				}
			}
		}
	}

	return score
}

func pointInTriangle(p image.Point, triangle [3]image.Point) bool {
	p0 := triangle[0]
	p1 := triangle[1]
	p2 := triangle[2]

	dX := p.X - p2.X
	dY := p.Y - p2.Y
	dX21 := p2.X - p1.X
	dY12 := p1.Y - p2.Y
	D := dY12*(p0.X-p2.X) + dX21*(p0.Y-p2.Y)
	s := dY12*dX + dX21*dY
	t := (p2.Y-p0.Y)*dX + (p0.X-p2.X)*dY
	if D < 0 {
		return s <= 0 && t <= 0 && s+t >= D
	}
	return s >= 0 && t >= 0 && s+t <= D
}

func GetRandomTriangle(bounds image.Rectangle, size int) ([3]image.Point, image.Point) {
	xCenter := RandomInt(bounds.Min.X+1, bounds.Max.X-1)
	yCenter := RandomInt(bounds.Min.Y+1, bounds.Max.Y-1)

	angle := random(0, math.Pi)

	p0Angle := random(0, math.Pi/3) + angle
	p1Angle := random(2*math.Pi/3, 3*math.Pi/3) + angle
	p2Angle := random(4*math.Pi/3, 5*math.Pi/3) + angle

	p0 := image.Point{
		X: int(math.Round(float64(xCenter) + float64(size)*math.Cos(p0Angle))),
		Y: int(math.Round(float64(yCenter) + float64(size)*math.Sin(p0Angle))),
	}
	p1 := image.Point{
		X: int(math.Round(float64(xCenter) + float64(size)*math.Cos(p1Angle))),
		Y: int(math.Round(float64(yCenter) + float64(size)*math.Sin(p1Angle))),
	}
	p2 := image.Point{
		X: int(math.Round(float64(xCenter) + float64(size)*math.Cos(p2Angle))),
		Y: int(math.Round(float64(yCenter) + float64(size)*math.Sin(p2Angle))),
	}

	return [3]image.Point{p0, p1, p2}, image.Point{X: xCenter, Y: yCenter}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func RandomInt(min int, max int) int {
	return int(rand.Float64()*float64(max+1-min) + float64(min))
}

func random(min float64, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
