package main

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/KMimura/ColonyGame/systems"
)

type MainScene struct{}

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
	// world.AddSystem(&systems.EnemySystem{})
	// world.AddSystem(&systems.BulletSystem{})
	// world.AddSystem(&systems.IntermissionSystem{})
}

func (*MainScene) Exit() {
	engo.Exit()
}

func main() {
	run()
}
