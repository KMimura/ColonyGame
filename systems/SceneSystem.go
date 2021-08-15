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
			var stage_tiles [screenLength][screenLength]tileInfo
			for i, s := range stage_tiles {
				for j, _ := range s {
					stage_tiles[i][j].spritesheetNum = rand.Intn(4)
				}
			}
			createRiver(w, &stage_tiles)
			for i, s := range stage_tiles {
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

func createRiver(w *ecs.World, stage_tiles *[screenLength][screenLength]tileInfo) {
	rand.Seed(time.Now().UnixNano())
	type river_info struct {
		X       int
		Y       int
		tilenum int
	}
	if_going_south := true
	if rand.Intn(2) == 1 {
		if_going_south = false
	}
	// just_curved := false
	// 川の始まりの地点の選択
	river_start_point := rand.Intn(screenLength/2) + screenLength/4
	var river_info_array []river_info
	// 初期位置作成
	if if_going_south {
		river_info_array = append(river_info_array, river_info{0, river_start_point, 60})
		river_info_array = append(river_info_array, river_info{0, river_start_point + 1, 61})
		river_info_array = append(river_info_array, river_info{0, river_start_point + 2, 62})
	} else {
		river_info_array = append(river_info_array, river_info{river_start_point, 0, 49})
		river_info_array = append(river_info_array, river_info{river_start_point + 1, 0, 61})
		river_info_array = append(river_info_array, river_info{river_start_point + 2, 0, 73})
	}

	// 初期値以降作成
	if if_going_south {
		for i := 1; i < screenLength; i++ {
			river_info_array = append(river_info_array, river_info{i, river_start_point, 60})
			river_info_array = append(river_info_array, river_info{i, river_start_point + 1, 61})
			river_info_array = append(river_info_array, river_info{i, river_start_point + 2, 62})
		}
	} else {
		for i := 1; i < screenLength; i++ {
			river_info_array = append(river_info_array, river_info{river_start_point, i, 49})
			river_info_array = append(river_info_array, river_info{river_start_point + 1, i, 61})
			river_info_array = append(river_info_array, river_info{river_start_point + 2, i, 73})
		}
	}

	// 引数として受けとったステージ情報を書き換える
	for _, r := range river_info_array {
		stage_tiles[r.X][r.Y].spritesheetNum = r.tilenum
	}
}
