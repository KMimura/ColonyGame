package systems

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type itemMenu struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type Text struct {
	ecs.BasicEntity
	common.SpaceComponent
	common.RenderComponent
}

type ItemMenuSystem struct {
	text Text
}

// Remove 削除する
func (ims *ItemMenuSystem) Remove(entity ecs.BasicEntity) {
	// for _, system := range ims.world.Systems() {
	// 	switch sys := system.(type) {
	// 	case *common.RenderSystem:
	// 		sys.Remove(entity)
	// 	}
	// }
}

// Update アップデートする
func (ims *ItemMenuSystem) Update(dt float32) {
}

// Init 初期化
func (ims *ItemMenuSystem) New(w *ecs.World) {
	fnt := &common.Font{
		URL:  "go.ttf",
		FG:   color.Black,
		Size: 24,
	}
	fnt.CreatePreloaded()

	ims.text = Text{BasicEntity: ecs.NewBasic()}
	ims.text.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "Hello, world!",
	}
	ims.text.SetShader(common.TextHUDShader)
	ims.text.RenderComponent.SetZIndex(1001)
	ims.text.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 20, Y: 20},
		Width:    200,
		Height:   200,
	}
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&ims.text.BasicEntity, &ims.text.RenderComponent, &ims.text.SpaceComponent)
		}
	}
}
