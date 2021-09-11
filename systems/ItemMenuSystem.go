package systems

import (
	"fmt"
	"image"
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

type ItemMenuSystem struct {
	world          *ecs.World
	itemMenuEntity *itemMenu
	texture        *common.Texture
}

func (ims *ItemMenuSystem) SetUp(w *ecs.World) {
	fmt.Println('s')
	itemMenu := itemMenu{BasicEntity: ecs.NewBasic()}
	itemMenu.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: engo.WindowWidth() - 50, Y: engo.WindowHeight() - 50},
		Width:    200,
		Height:   200,
	}
	itemMenu.RenderComponent.SetZIndex(1)
	hudImage := image.NewUniform(color.RGBA{0, 0, 0, 0})
	hudNRGBA := common.ImageToNRGBA(hudImage, 16, 16)
	hudImageObj := common.NewImageObject(hudNRGBA)
	hudTexture := common.NewTextureSingle(hudImageObj)
	itemMenu.RenderComponent = common.RenderComponent{
		Repeat:   common.Repeat,
		Drawable: hudTexture,
		Scale:    engo.Point{X: 1, Y: 1},
	}
	itemMenu.RenderComponent.SetShader(common.HUDShader)
	itemMenu.RenderComponent.SetZIndex(1)
	for _, system := range ims.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&itemMenu.BasicEntity, &itemMenu.RenderComponent, &itemMenu.SpaceComponent)
		}
	}
}

// Remove 削除する
func (ims *ItemMenuSystem) Remove(entity ecs.BasicEntity) {
	for _, system := range ims.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

// Update アップデートする
func (ims *ItemMenuSystem) Update(dt float32) {
}

// Init 初期化
func (ims *ItemMenuSystem) Init(w *ecs.World) {
}
