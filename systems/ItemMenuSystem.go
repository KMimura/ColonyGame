package systems

import (
	"image"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type ItemMenu struct {
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
	text                         Text // 表示するテキスト
	menuButonPushed              bool // メニュー表示ボタンが押下された状態かどうか
	menuButonPushedRemainingTime int  // あと何フレームメニューボタン押下を無効化するか

}

var menuZIndex int = -10

var itemMenuInstance *ItemMenu

var buttonDisableTime = 15 // 一度押下されたボタンをどれだけ無効化するか

func ItemMenuInit(world *ecs.World) {
	itemMenu := ItemMenu{BasicEntity: ecs.NewBasic()}
	itemMenu.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 20, Y: 20},
		Width:    300,
		Height:   900,
	}
	itemMenu.RenderComponent.SetZIndex(1)
	hudImage := image.NewUniform(color.RGBA{175, 175, 175, 225})
	hudNRGBA := common.ImageToNRGBA(hudImage, 300, 900)
	hudImageObj := common.NewImageObject(hudNRGBA)
	hudTexture := common.NewTextureSingle(hudImageObj)
	itemMenu.RenderComponent = common.RenderComponent{
		Repeat:   common.Repeat,
		Drawable: hudTexture,
		Scale:    engo.Point{X: 1, Y: 1},
	}
	itemMenu.RenderComponent.SetShader(common.HUDShader)
	itemMenu.RenderComponent.SetZIndex(float32(menuZIndex))
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&itemMenu.BasicEntity, &itemMenu.RenderComponent, &itemMenu.SpaceComponent)
		}
	}
	itemMenuInstance = &itemMenu
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
	if ims.menuButonPushed {
		if ims.menuButonPushedRemainingTime > 0 {
			ims.menuButonPushedRemainingTime--
		} else {
			ims.menuButonPushedRemainingTime = 0
			ims.menuButonPushed = false
		}
	} else {
		if engo.Input.Button("Space").Down() {
			if menuZIndex != -10 {
				menuZIndex = -10
			} else {
				menuZIndex = 1
			}
			ims.menuButonPushed = true
			ims.menuButonPushedRemainingTime = buttonDisableTime
			itemMenuInstance.SetZIndex(float32(menuZIndex))
			ims.text.RenderComponent.SetZIndex(float32(menuZIndex + 1))
		}
	}
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
	ims.menuButonPushed = false
	ims.menuButonPushedRemainingTime = 0
	ims.text.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "Hello, world!",
	}
	ims.text.SetShader(common.TextHUDShader)
	ims.text.RenderComponent.SetZIndex(float32(menuZIndex))
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
