package systems

import (
	"math/rand"

	// "reflect"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Spritesheet タイルの画像
var Spritesheet *common.Spritesheet

// camEntity カメラシステムのエンティティ
var camEntity *common.CameraSystem

// cameraInitialPositionX,Y カメラの初期位置
var cameraInitialPositionX int
var cameraInitialPositionY int

// cellLength セル一辺のピクセル数（必ず16の倍数にすること）
const cellLength = 48

const screenLength = 50

type tileInfo struct {
	spritesheetNum int
}

// Tile タイル一つ一つを表す構造体
type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// SceneSystem シーンシステム
type SceneSystem struct {
	world   *ecs.World
	texture *common.Texture
}

// Remove 削除する
func (ss *SceneSystem) Remove(entity ecs.BasicEntity) {
	for _, system := range ss.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

// Update アップデートする
func (ss *SceneSystem) Update(dt float32) {
	if engo.Input.Button("MoveUp").Down() {
		engo.Mailbox.Dispatch(common.CameraMessage{
			Axis:        common.YAxis,
			Value:       -16,
			Incremental: true,
		})
	} else if engo.Input.Button("MoveDown").Down() {
		engo.Mailbox.Dispatch(common.CameraMessage{
			Axis:        common.YAxis,
			Value:       16,
			Incremental: true,
		})
	} else if engo.Input.Button("MoveLeft").Down() {
		engo.Mailbox.Dispatch(common.CameraMessage{
			Axis:        common.XAxis,
			Value:       -16,
			Incremental: true,
		})
	} else if engo.Input.Button("MoveRight").Down() {
		engo.Mailbox.Dispatch(common.CameraMessage{
			Axis:        common.XAxis,
			Value:       16,
			Incremental: true,
		})
	}
}

// New 作成時に呼び出される
func (ss *SceneSystem) New(w *ecs.World) {
	ss.init(w)
}

// Init 初期化
func (ss *SceneSystem) init(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())
	ss.world = w
	// 素材シートの読み込み
	loadTxt := "pics/overworld_tileset_grass.png"
	Spritesheet = common.NewSpritesheetWithBorderFromFile(loadTxt, 16, 16, 0, 0)
	// createRiver(ss.world)
	for _, system := range ss.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			var stageTiles [screenLength][screenLength]tileInfo
			for i, s := range stageTiles {
				for j, _ := range s {
					stageTiles[i][j].spritesheetNum = rand.Intn(4)
				}
			}
			createRiver(w, &stageTiles)
			for i, s := range stageTiles {
				for j, y := range s {
					tile := &Tile{BasicEntity: ecs.NewBasic()}
					tile.SpaceComponent.Position = engo.Point{
						X: float32(j * cellLength),
						Y: float32(i * cellLength),
					}
					tile.RenderComponent = common.RenderComponent{
						Drawable: Spritesheet.Cell(y.spritesheetNum),
						Scale:    engo.Point{X: float32(cellLength / 16), Y: float32(cellLength / 16)},
					}
					tile.RenderComponent.SetZIndex(0)
					sys.Add(&tile.BasicEntity, &tile.RenderComponent, &tile.SpaceComponent)
				}
			}
		}
	}
	// カメラエンティティの取得
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.CameraSystem:
			camEntity = sys
			common.CameraBounds.Max.X = float32(screenLength * cellLength)
			common.CameraBounds.Max.Y = float32(screenLength * cellLength)
		}
	}
}

func createRiver(w *ecs.World, stageTiles *[screenLength][screenLength]tileInfo) {
	rand.Seed(time.Now().UnixNano())
	type riverInfo struct {
		X       int
		Y       int
		tilenum int
	}
	ifGoingSouth := true
	if rand.Intn(2) == 1 {
		ifGoingSouth = false
	}
	// just_curved := false
	// 川の始まりの地点の選択
	riverStartPoint := rand.Intn(screenLength/2) + screenLength/4
	var riverInfoArray []riverInfo
	// 初期位置作成
	if ifGoingSouth {
		riverInfoArray = append(riverInfoArray, riverInfo{0, riverStartPoint, 60})
		riverInfoArray = append(riverInfoArray, riverInfo{0, riverStartPoint + 1, 61})
		riverInfoArray = append(riverInfoArray, riverInfo{0, riverStartPoint + 2, 62})
	} else {
		riverInfoArray = append(riverInfoArray, riverInfo{riverStartPoint, 0, 49})
		riverInfoArray = append(riverInfoArray, riverInfo{riverStartPoint + 1, 0, 61})
		riverInfoArray = append(riverInfoArray, riverInfo{riverStartPoint + 2, 0, 73})
	}

	// 初期値以降作成
	// 川の描画を終えた座標
	var riverCursorX int
	var riverCursorY int
	if ifGoingSouth {
		riverCursorX = riverStartPoint
		riverCursorY = 1
	} else {
		riverCursorX = 1
		riverCursorY = riverStartPoint
	}
	// 川の蛇行を始めてからの時間
	curveGen := 0
	shouldContinue := true
	var yArray [3]int
	var xArray [3]int
	var tileNum [3]int
	for shouldContinue {
		if ifGoingSouth {
			switch curveGen {
			case 0:
				yArray[0] = riverCursorY
				yArray[1] = riverCursorY
				yArray[2] = riverCursorY
				xArray[0] = riverCursorX
				xArray[1] = riverCursorX + 1
				xArray[2] = riverCursorX + 2
				tileNum[0] = 60
				tileNum[1] = 61
				tileNum[2] = 62
			case 1:
				yArray[0] = riverCursorY
				yArray[1] = riverCursorY
				yArray[2] = riverCursorY
				xArray[0] = riverCursorX
				xArray[1] = riverCursorX + 1
				xArray[2] = riverCursorX + 2
				tileNum[0] = 49
				tileNum[1] = 49
				tileNum[2] = 50
				curveGen = 2
			case 2:
				yArray[0] = riverCursorY
				yArray[1] = riverCursorY
				yArray[2] = riverCursorY
				xArray[0] = riverCursorX
				xArray[1] = riverCursorX + 1
				xArray[2] = riverCursorX + 2
				tileNum[0] = 61
				tileNum[1] = 61
				tileNum[2] = 62
				curveGen = 3
			case 3:
				yArray[0] = riverCursorY
				yArray[1] = riverCursorY
				yArray[2] = riverCursorY
				xArray[0] = riverCursorX
				xArray[1] = riverCursorX + 1
				xArray[2] = riverCursorX + 2
				tileNum[0] = 85
				tileNum[1] = 61
				tileNum[2] = 62
				curveGen = 0
			}
			riverCursorY++
			if riverCursorY >= screenLength || riverCursorX+2 >= screenLength {
				shouldContinue = false
			}
			if curveGen == 0 && rand.Intn(15) == 0 {
				ifGoingSouth = !ifGoingSouth
				curveGen = 1
			}
		} else {
			switch curveGen {
			case 0:
				yArray[0] = riverCursorY
				yArray[1] = riverCursorY + 1
				yArray[2] = riverCursorY + 2
				xArray[0] = riverCursorX
				xArray[1] = riverCursorX
				xArray[2] = riverCursorX
				tileNum[0] = 49
				tileNum[1] = 61
				tileNum[2] = 73
			case 1:
				yArray[0] = riverCursorY
				yArray[1] = riverCursorY + 1
				yArray[2] = riverCursorY + 2
				xArray[0] = riverCursorX
				xArray[1] = riverCursorX
				xArray[2] = riverCursorX
				tileNum[0] = 60
				tileNum[1] = 60
				tileNum[2] = 72
				curveGen = 2
			case 2:
				yArray[0] = riverCursorY
				yArray[1] = riverCursorY + 1
				yArray[2] = riverCursorY + 2
				xArray[0] = riverCursorX
				xArray[1] = riverCursorX
				xArray[2] = riverCursorX
				tileNum[0] = 61
				tileNum[1] = 61
				tileNum[2] = 73
				curveGen = 3
			case 3:
				yArray[0] = riverCursorY
				yArray[1] = riverCursorY + 1
				yArray[2] = riverCursorY + 2
				xArray[0] = riverCursorX
				xArray[1] = riverCursorX
				xArray[2] = riverCursorX
				tileNum[0] = 96
				tileNum[1] = 61
				tileNum[2] = 73
				curveGen = 0
			}
			riverCursorX++
			if riverCursorX >= screenLength || riverCursorY+2 >= screenLength {
				shouldContinue = false
			}
			if curveGen == 0 && rand.Intn(15) == 0 {
				ifGoingSouth = !ifGoingSouth
				curveGen = 1
			}
		}
		riverInfoArray = append(riverInfoArray, riverInfo{yArray[0], xArray[0], tileNum[0]})
		riverInfoArray = append(riverInfoArray, riverInfo{yArray[1], xArray[1], tileNum[1]})
		riverInfoArray = append(riverInfoArray, riverInfo{yArray[2], xArray[2], tileNum[2]})
	}
	// 引数として受けとったステージ情報を書き換える
	for _, r := range riverInfoArray {
		stageTiles[r.X][r.Y].spritesheetNum = r.tilenum
	}
}
