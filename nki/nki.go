package nki

import (
	_ "embed"
	"math/rand"
	"strings"

	"github.com/oakmound/oak/v3/physics"
)

//go:embed vanilla.nki
var nki string

var positions = []physics.Vector{}

func hasPosition(pos physics.Vector) bool {
	for _, x := range positions {
		if x == pos {
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
	for hasPosition(physics.NewVector(x, y)) {
		x = float64(clamp(rand.Intn(1000), 40))
		y = float64(clamp(rand.Intn(1000), 40))
	}
	return physics.NewVector(x, y)
}
