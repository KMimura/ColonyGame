package main

import (
	"bytes"
	// "image"
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/KMimura/ColonyGame/systems"
	"golang.org/x/image/font/gofont/gosmallcaps"
)

type MainScene struct{}

// type itemMenu struct {
// 	ecs.BasicEntity
// 	common.RenderComponent
// 	common.SpaceComponent
// }

func run() {
	opts := engo.RunOptions{
		Title:          "ColonyGame",
		Width:          1200,
		Height:         900,
		StandardInputs: true,
		NotResizable:   true,
	}
	engo.Run(opts, &MainScene{})
}

func (*MainScene) Type() string { return "mainScene" }

func (*MainScene) Preload() {
	engo.Files.Load("pics/overworld_tileset_grass.png")
	engo.Files.Load("pics/characters.png")
	engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
	// engo.Files.LoadReaderData("go.ttf", bytes.NewReader(gosmallcaps.TTF))
	common.SetBackground(color.RGBA{255, 250, 220, 0})
}

func (*MainScene) Setup(u engo.Updater) {
	engo.Input.RegisterButton("MoveRight", engo.KeyD, engo.KeyArrowRight)
	engo.Input.RegisterButton("MoveLeft", engo.KeyA, engo.KeyArrowLeft)
	engo.Input.RegisterButton("MoveUp", engo.KeyW, engo.KeyArrowUp)
	engo.Input.RegisterButton("MoveDown", engo.KeyS, engo.KeyArrowDown)
	engo.Input.RegisterButton("Space", engo.KeySpace)
	engo.Input.RegisterButton("Escape", engo.KeyEscape)
	world, _ := u.(*ecs.World)
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&systems.SceneSystem{})
	world.AddSystem(&systems.PlayerSystem{})
	systems.ItemMenuInit(world)

	// itemMenu := itemMenu{BasicEntity: ecs.NewBasic()}
	// itemMenu.SpaceComponent = common.SpaceComponent{
	// 	Position: engo.Point{X: 20, Y: 20},
	// 	Width:    300,
	// 	Height:   900,
	// }
	// itemMenu.RenderComponent.SetZIndex(1)
	// hudImage := image.NewUniform(color.RGBA{175, 175, 175, 225})
	// hudNRGBA := common.ImageToNRGBA(hudImage, 300, 900)
	// hudImageObj := common.NewImageObject(hudNRGBA)
	// hudTexture := common.NewTextureSingle(hudImageObj)
	// itemMenu.RenderComponent = common.RenderComponent{
	// 	Repeat:   common.Repeat,
	// 	Drawable: hudTexture,
	// 	Scale:    engo.Point{X: 1, Y: 1},
	// }
	// itemMenu.RenderComponent.SetShader(common.HUDShader)
	// itemMenu.RenderComponent.SetZIndex(1)
	// for _, system := range world.Systems() {
	// 	switch sys := system.(type) {
	// 	case *common.RenderSystem:
	// 		sys.Add(&itemMenu.BasicEntity, &itemMenu.RenderComponent, &itemMenu.SpaceComponent)
	// 	}
	// }
	world.AddSystem(&systems.ItemMenuSystem{})
	// world.AddSystem(&systems.BulletSystem{})
	// world.AddSystem(&systems.IntermissionSystem{})
}

func (*MainScene) Exit() {
	engo.Exit()
}

func main() {
	run()
}
