package nki

import (
	_ "embed"
	"math"
	"math/rand"
	"strings"

	"github.com/oakmound/oak/v3/physics"
)

//go:embed vanilla.nki
var nki string

var positions = []physics.Vector{}

func eqf(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func hasPosition(x float64, y float64) bool {
	for _, pos := range positions {
		if eqf(pos.X(), x) && eqf(pos.Y(), y) {
			return true
		}
	}
	return false
}

func clamp(n, to int) int {
	return (n / to) * to // Yay integer division!
}

func Generate() string {
	splitNki := strings.Split(nki, "\n")
	return splitNki[rand.Intn(len(splitNki))]
}

func GetPosition() physics.Vector {
	x := float64(clamp(rand.Intn(1000), 40))
	y := float64(clamp(rand.Intn(1000), 40))
	for hasPosition(x, y) {
		x = float64(clamp(rand.Intn(1000), 40))
		y = float64(clamp(rand.Intn(1000), 40))
	}
	positions = append(positions, physics.NewVector(x, y))
	return physics.NewVector(x, y)
}
