package tools

import (
	"fmt"
	"testing"
)

//func TestValidAllowlist(t *testing.T) {
//	p0 := image.Point{
//		X: 100,
//		Y: 100,
//	}
//	p1 := image.Point{
//		X: 400,
//		Y: 100,
//	}
//	p2 := image.Point{
//		X: 100,
//		Y: 400,
//	}
//
//	require.False(t, PointInTriangle(image.Point{0, 0}, p0, p1, p2))
//	require.True(t, PointInTriangle(image.Point{100, 100}, p0, p1, p2))
//	require.False(t, PointInTriangle(image.Point{300, 300}, p0, p1, p2))
//
//	r := uint32(255)
//	fmt.Println(r)
//	r |= r << 8
//	fmt.Println(r)
//	r &= r >> 8
//	fmt.Println(r)
//}

func TestRandomInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(RandomInt(255, 256))
	}
}
