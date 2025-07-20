package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	paddleWidth  = 100
	paddleHeight = 20
	ballSize     = 10
	blockWidth   = 80
	blockHeight  = 30
	blockRows    = 5
	blockCols    = 10
)

type Game struct {
	paddleX     float64
	paddleY     float64
	ballX       float64
	ballY       float64
	ballVelX    float64
	ballVelY    float64
	prevBallX   float64
	prevBallY   float64
	blocks      [][]bool
	score       int
	gameOver    bool
	gameWon     bool
	gameStarted bool
}

func NewGame() *Game {
	// ブロックの初期化
	blocks := make([][]bool, blockRows)
	for i := range blocks {
		blocks[i] = make([]bool, blockCols)
		for j := range blocks[i] {
			blocks[i][j] = true
		}
	}

	return &Game{
		paddleX:   (screenWidth - paddleWidth) / 2,
		paddleY:   screenHeight - 50,
		ballX:     screenWidth / 2,
		ballY:     screenHeight - 70,
		ballVelX:  0,
		ballVelY:  0,
		prevBallX: screenWidth / 2,
		prevBallY: screenHeight - 70,
		blocks:    blocks,
		score:     0,
	}
}

func (g *Game) Update() error {
	if g.gameOver || g.gameWon {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			*g = *NewGame()
		}
		return nil
	}

	// パドルの移動
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.paddleX > 0 {
		g.paddleX -= 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && g.paddleX < screenWidth-paddleWidth {
		g.paddleX += 5
	}

	// ゲーム開始
	if !g.gameStarted {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.gameStarted = true
			g.ballVelX = 3
			g.ballVelY = -3
		}
		return nil
	}

	// 前フレームの位置を保存
	g.prevBallX = g.ballX
	g.prevBallY = g.ballY

	// ボールの移動
	g.ballX += g.ballVelX
	g.ballY += g.ballVelY

	// 壁との衝突判定
	if g.ballX <= 0 || g.ballX >= screenWidth-ballSize {
		g.ballVelX = -g.ballVelX
		g.ballX = g.prevBallX // 位置を戻す
	}
	if g.ballY <= 0 {
		g.ballVelY = -g.ballVelY
		g.ballY = g.prevBallY // 位置を戻す
	}

	// パドルとの衝突判定
	if g.ballY+ballSize >= g.paddleY && g.ballY <= g.paddleY+paddleHeight &&
		g.ballX+ballSize >= g.paddleX && g.ballX <= g.paddleX+paddleWidth {
		// 入射角と反射角を等しくする（横方向のベクトルは変化しない）
		g.ballVelY = -g.ballVelY
		g.ballY = g.paddleY - ballSize // パドルの上に配置
		// 横方向の速度はそのまま保持（入射角と反射角が等しい）
	}

	// ブロックとの衝突判定（改善版）
	g.checkBlockCollisions()

	// ゲームオーバー判定
	if g.ballY >= screenHeight {
		g.gameOver = true
	}

	// 勝利判定
	allBlocksDestroyed := true
	for _, row := range g.blocks {
		for _, block := range row {
			if block {
				allBlocksDestroyed = false
				break
			}
		}
		if !allBlocksDestroyed {
			break
		}
	}
	if allBlocksDestroyed {
		g.gameWon = true
	}

	return nil
}

// ブロック衝突判定の改善版
func (g *Game) checkBlockCollisions() {
	// ボールの現在位置と前フレーム位置から衝突をチェック
	ballLeft := g.ballX
	ballRight := g.ballX + ballSize
	ballTop := g.ballY
	ballBottom := g.ballY + ballSize

	prevBallLeft := g.prevBallX
	prevBallRight := g.prevBallX + ballSize
	prevBallTop := g.prevBallY
	prevBallBottom := g.prevBallY + ballSize

	// ボールが通る領域のブロックをチェック
	startRow := int(math.Min(prevBallTop, ballTop)) / blockHeight
	endRow := int(math.Max(prevBallBottom, ballBottom)) / blockHeight
	startCol := int(math.Min(prevBallLeft, ballLeft)) / blockWidth
	endCol := int(math.Max(prevBallRight, ballRight)) / blockWidth

	// 範囲を制限
	startRow = int(math.Max(0, float64(startRow)))
	endRow = int(math.Min(float64(blockRows-1), float64(endRow)))
	startCol = int(math.Max(0, float64(startCol)))
	endCol = int(math.Min(float64(blockCols-1), float64(endCol)))

	collision := false
	for row := startRow; row <= endRow; row++ {
		for col := startCol; col <= endCol; col++ {
			if g.blocks[row][col] {
				// ブロックの境界
				blockLeft := float64(col * blockWidth)
				blockRight := blockLeft + blockWidth
				blockTop := float64(row * blockHeight)
				blockBottom := blockTop + blockHeight

				// 衝突判定
				if ballRight > blockLeft && ballLeft < blockRight &&
					ballBottom > blockTop && ballTop < blockBottom {

					g.blocks[row][col] = false
					g.score += 10
					collision = true

					// 衝突方向を判定してボールの方向を変更
					// 左右からの衝突
					if (prevBallRight <= blockLeft && ballRight > blockLeft) ||
						(prevBallLeft >= blockRight && ballLeft < blockRight) {
						g.ballVelX = -g.ballVelX
						if g.ballVelX > 0 {
							g.ballX = blockRight
						} else {
							g.ballX = blockLeft - ballSize
						}
					} else {
						// 上下からの衝突
						g.ballVelY = -g.ballVelY
						if g.ballVelY > 0 {
							g.ballY = blockBottom
						} else {
							g.ballY = blockTop - ballSize
						}
					}
					break
				}
			}
		}
		if collision {
			break
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 背景を黒で塗りつぶし
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.Black)

	// パドルを描画
	ebitenutil.DrawRect(screen, g.paddleX, g.paddleY, paddleWidth, paddleHeight, color.RGBA{0, 255, 0, 255})

	// ボールを描画
	ebitenutil.DrawRect(screen, g.ballX, g.ballY, ballSize, ballSize, color.White)

	// ブロックを描画
	for i, row := range g.blocks {
		for j, block := range row {
			if block {
				x := float64(j * blockWidth)
				y := float64(i * blockHeight)
				// 行ごとに異なる色を使用
				r := uint8(float64(i) / float64(blockRows) * 255)
				g := uint8((1.0 - float64(i)/float64(blockRows)) * 255)
				b := uint8(float64(i) / float64(blockRows) * 255)
				ebitenutil.DrawRect(screen, x, y, blockWidth-2, blockHeight-2, color.RGBA{r, g, b, 255})
			}
		}
	}

	// スコアを表示
	ebitenutil.DebugPrint(screen, "スコア: "+fmt.Sprintf("%d", g.score))

	if !g.gameStarted {
		ebitenutil.DebugPrintAt(screen, "スペースキーを押してゲーム開始", 300, 250)
	}

	if g.gameOver {
		ebitenutil.DebugPrintAt(screen, "ゲームオーバー！", 350, 250)
		ebitenutil.DebugPrintAt(screen, "Rキーでリスタート", 350, 270)
	}

	if g.gameWon {
		ebitenutil.DebugPrintAt(screen, "勝利！", 370, 250)
		ebitenutil.DebugPrintAt(screen, "Rキーでリスタート", 350, 270)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ブロック崩しゲーム")

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
