package util

import "github.com/oakmound/oak/v3/render"

func AssertSprite(spr *render.Sprite, err error) *render.Sprite {
	if err != nil {
		panic(err)
	}
	return spr
}
