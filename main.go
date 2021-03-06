package main

import (
	"image"
	"image/color"
	"math/rand"
	"path/filepath"

	"github.com/aleksrutins/robotfindskitten/nki"
	"github.com/aleksrutins/robotfindskitten/util"
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

func randBetween(from, to int) int {
	return rand.Intn(to-from) + from
}

func lerp(v0, v1, t float32) float32 {
	return v0 + t*(v1-v0)
}

var (
	messages    = map[event.CID]string{}
	foundKitten = false
)

func main() {
	oak.AddScene("rfk", scene.Scene{
		Start: func(ctx *scene.Context) {
			oak.SetColorBackground(image.NewUniform(color.White))

			robotSprite, err := render.LoadSprite(filepath.Join("assets", "images", "Robot.png"))
			dlog.ErrorCheck(err)

			robot := entities.NewMoving(100, 100, 32, 32, robotSprite, nil, 0, 0)
			render.Draw(robot.R, 3)
			robot.Speed = physics.NewVector(5, 5)

			hitLastFrame := &collision.Space{}

			fg := render.DefaultFontGenerator
			fg.Color = image.NewUniform(color.Black)
			font, err := fg.Generate()
			dlog.ErrorCheck(err)
			font.Unsafe = true
			text := font.NewText("Hello", 0, 0)

			render.Draw(text, 4)

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

					if messages[hit.CID] == "You found kitten!" {
						foundKitten = true
					}
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
				item := entities.NewSolid(position.X(), position.Y(), 32, 32,
					util.AssertSprite(render.LoadSprite(filepath.Join("assets", "images", "Item.png"))),
					nil, 0)
				itemText := font.NewText(string(randBetween(33, 47)), position.X()+12, position.Y()+8)
				render.Draw(item.R, 1)
				render.Draw(itemText, 2)
				item.UpdateLabel(Item)
				messages[item.CID] = msg
			}
			// Generate kitten
			msg := "You found kitten!"
			position := nki.GetPosition()
			item := entities.NewSolid(position.X(), position.Y(), 32, 32,
				util.AssertSprite(render.LoadSprite(filepath.Join("assets", "images", "Item.png"))),
				nil, 0)
			itemText := font.NewText(string(randBetween(33, 47)), position.X()+12, position.Y()+8)
			render.Draw(item.R, 1)
			render.Draw(itemText, 2)
			item.UpdateLabel(Item)
			messages[item.CID] = msg
		},
		Loop: func() (cont bool) { return !foundKitten },
		End:  func() (nextScene string, result *scene.Result) { return "end", nil },
	})

	oak.AddScene("end", scene.Scene{
		Start: func(ctx *scene.Context) {
			oak.SetColorBackground(image.Black)
			text := render.NewText("You found kitten! Press Space to restart.", float64(oak.Width())/2, float64(oak.Height())/2)
			w, h := text.GetDims()
			text.ShiftX(-float64(w) / 2)
			text.ShiftY(-float64(h) / 2)
			render.Draw(text)
		},
		Loop: func() (cont bool) { return !oak.IsDown(key.Spacebar) },
		End: func() (nextScene string, result *scene.Result) {
			foundKitten = false
			return "rfk", nil
		},
	})
	oak.Init("rfk")
}
