package main

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/aleksrutins/robotfindskitten/nki"
	"github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/collision"
	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/key"
	"github.com/oakmound/oak/v3/physics"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

const (
	Item collision.Label = 1
)

func hasKey(m map[event.CID]string, item event.CID) bool {
	for k := range m {
		if k == item {
			return true
		}
	}
	return false
}

func lerp(v0, v1, t float32) float32 {
	return v0 + t*(v1-v0)
}

var messages = map[event.CID]string{}

func main() {
	oak.AddScene("rfk", scene.Scene{
		Start: func(ctx *scene.Context) {
			robot := entities.NewMoving(100, 100, 32, 32, render.NewColorBox(32, 32, color.RGBA{255, 0, 0, 255}), nil, 0, 0)
			render.Draw(robot.R)
			robot.Speed = physics.NewVector(5, 5)

			hitLastFrame := &collision.Space{}

			fg := render.DefaultFontGenerator
			fg.Color = image.NewUniform(color.White)
			font, err := fg.Generate()
			dlog.ErrorCheck(err)
			font.Unsafe = true
			text := font.NewText("Hello", 0, 0)

			render.Draw(text, 2)

			robot.Bind(event.Enter, func(id event.CID, _ interface{}) int {
				char := event.GetEntity(id).(*entities.Moving)
				if oak.IsDown(key.W) {
					char.Delta.SetY(-char.Speed.Y())
				} else if oak.IsDown(key.S) {
					char.Delta.SetY(char.Speed.Y())
				} else {
					char.Delta.SetY(0)
				}
				if oak.IsDown(key.A) {
					char.Delta.SetX(-char.Speed.X())
				} else if oak.IsDown(key.D) {
					char.Delta.SetX(char.Speed.X())
				} else {
					char.Delta.SetX(0)
				}
				char.ShiftVector(char.Delta)

				oak.SetScreen(
					int(lerp(
						float32(ctx.Window.Viewport().X()),
						float32(char.R.X()-float64(ctx.Window.Width()/2)),
						0.1,
					)),
					int(lerp(
						float32(ctx.Window.Viewport().Y()),
						float32(char.R.Y()-float64(ctx.Window.Height()/2)),
						0.1,
					)),
				)
				text.SetPos(float64(ctx.Window.Viewport().X()+20), float64(ctx.Window.Viewport().Y()+20))

				hit := collision.HitLabel(char.Space, Item)
				if hit != nil && hitLastFrame != hit && hasKey(messages, hit.CID) {
					println(messages[hit.CID])
					text.SetString(messages[hit.CID])
					hitLastFrame = hit
				} else if hit == nil {
					text.SetString("")
					hitLastFrame = nil
				}

				return 0
			})

			// Generate NKIs
			for i := 0; i < 50; i++ {
				msg := nki.Generate()
				position := nki.GetPosition()
				item := entities.NewSolid(position.X(), position.Y(), 32, 32, render.NewColorBoxR(32, 32, color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}), nil, 0)
				render.Draw(item.R)
				item.UpdateLabel(Item)
				messages[item.CID] = msg
			}
		},
	})
	oak.Init("rfk")
}
