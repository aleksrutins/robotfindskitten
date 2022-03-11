package nki

import (
	"path/filepath"

	"github.com/aleksrutins/robotfindskitten/util"
	"github.com/oakmound/oak/v3/entities"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/render"
)

var itemSprite = util.AssertSprite(render.LoadSprite(filepath.Join("assets", "images", "Item.png")))

type Item struct {
	*entities.Solid
	IsKitten bool
	Message  string
}

func (item *Item) Init() event.CID {
	return event.NextID(item)
}

func NewItem(isKitten bool) *Item {
	i := new(Item)
	msg := Generate()
	position := GetPosition()
	i.Solid = entities.NewSolid(position.X(), position.Y(), 32, 32, itemSprite, nil, 0)
	i.Message = msg
	i.IsKitten = isKitten
	return i
}

func (item *Item) HandleCollision(alertText *render.Text) {
	alertText.SetString(item.Message)
	if item.IsKitten {

	}
}
