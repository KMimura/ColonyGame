package systems

import (
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Player プレーヤーを表す構造体
type Player struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	direction        int     // 向き (0 => 移動中でない, 1 => 上, 2 => 右, 3 => 下 4 => 左)
	remainingHearts  int     // ライフ
	immunityTime     int     // ダメージを受けない状態の残り時間
	velocity         float32 // 移動の速度
	cellX            int     // セルのX座標
	cellY            int     // セルのY座標
	destinationPoint float32 // 移動の目標地点の座標
	facingDirection  int     // どの方向を向いているか (1 => 上, 2 => 右, 3 => 下 4 => 左)
	movingPic        bool    //移動中の画像を表示するかどうか
}

// PlayerSystem プレーヤーシステム
type PlayerSystem struct {
	world        *ecs.World
	playerEntity *Player
	texture      *common.Texture
}

// playerInstance プレーヤーのエンティティのインスタンス
var playerInstance *Player

// 画像の半径
var playerRadius float32 = 12.5

var transparentPic *common.Texture

// それぞれの向きのプレーヤーの画像
var topPicOne *common.Texture
var topPicTwo *common.Texture
var topPicThree *common.Texture
var rightPicOne *common.Texture
var rightPicTwo *common.Texture
var rightPicThree *common.Texture
var bottomPicOne *common.Texture
var bottomPicTwo *common.Texture
var bottomPicThree *common.Texture
var leftPicOne *common.Texture
var leftPicTwo *common.Texture
var leftPicThree *common.Texture

// New 新規作成時に呼び出される
func (ps *PlayerSystem) New(w *ecs.World) {
	ps.Init(w)
}

// Remove 削除する
func (ps *PlayerSystem) Remove(entity ecs.BasicEntity) {
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Remove(entity)
		}
	}
}

// Update アップデートする
func (ps *PlayerSystem) Update(dt float32) {
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

// Init 初期化
func (ps *PlayerSystem) Init(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())
	ps.world = w
	// プレーヤーの作成
	player := Player{BasicEntity: ecs.NewBasic()}

	// ライフを与える
	player.remainingHearts = 5
	// 移動はしていない
	player.direction = 0
	player.facingDirection = 1
	player.movingPic = false

	playerInstance = &player

	// 初期の配置
	ifKeepSearching := true
	if ifKeepSearching {
		tmpX := rand.Intn(screenLength)
		tmpY := rand.Intn(screenLength)
		if checkIfPassable(tmpX, tmpY) {
			ifKeepSearching = false
			player.cellX = tmpX
			player.cellY = tmpY
		}
	}
	positionX := cellLength * player.cellX
	positionY := cellLength * player.cellY
	player.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: float32(positionX), Y: float32(positionY)},
		Width:    30,
		Height:   30,
	}
	// 速度
	player.velocity = 4
	// 画像の読み込み
	loadTxt := "pics/characters.png"
	Spritesheet = common.NewSpritesheetWithBorderFromFile(loadTxt, 32, 32, 0, 0)

	topPicTmpOne := Spritesheet.Cell(2)
	topPicOne = &topPicTmpOne

	player.RenderComponent = common.RenderComponent{
		Drawable: topPicOne,
		Scale:    engo.Point{X: 1.5, Y: 1.5},
	}
	player.RenderComponent.SetZIndex(1)
	ps.playerEntity = &player
	ps.texture = topPicOne
	for _, system := range ps.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&player.BasicEntity, &player.RenderComponent, &player.SpaceComponent)
		}
	}
}
