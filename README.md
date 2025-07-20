# ブロック崩しゲーム

Go言語とEbitenライブラリを使用して作成されたブロック崩しゲームです。

## 機能

- パドルでボールを操作してブロックを破壊
- スコアシステム
- 美しいカラフルなブロック
- ゲームオーバー・勝利判定
- リスタート機能

## 操作方法

- **左右矢印キー**: パドルを左右に移動
- **スペースキー**: ゲーム開始
- **Rキー**: ゲームオーバーまたは勝利後にリスタート

## 実行方法

### 前提条件

- Go 1.16以上がインストールされていること

### インストールと実行

1. 依存関係をインストール:
```bash
go mod tidy
```

2. ゲームを実行:
```bash
go run main.go
```

### ビルド

実行可能ファイルを作成する場合:
```bash
go build -o breakout-game.exe
```

## ゲームルール

1. スペースキーを押してゲームを開始
2. 左右矢印キーでパドルを操作
3. ボールがパドルに当たると跳ね返る
4. ブロックにボールが当たるとブロックが破壊され、スコアが加算される
5. すべてのブロックを破壊すると勝利
6. ボールを落とすとゲームオーバー

## 技術仕様

- **言語**: Go
- **ゲームエンジン**: Ebiten v2
- **画面サイズ**: 800x600ピクセル
- **フレームレート**: 60FPS
- **対応プラットフォーム**: Windows, macOS, Linux, WebAssembly

## WebAssembly版

このゲームはWebAssembly（WASM）でも動作します。

### オンライン版
GitHub Pagesでホストされているオンライン版は以下のURLでアクセスできます：
```
https://maniax-jp.github.io/sampleProgram2/

```

### ローカルでのWASM実行
1. WASMファイルをビルド:
```bash
GOOS=js GOARCH=wasm go build -o main.wasm main.go
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

2. ローカルサーバーを起動:
```bash
python -m http.server 8080
# または
npx serve .
```

3. ブラウザで `http://localhost:8080` にアクセス

## 自動デプロイ

このプロジェクトはGitHub Actionsを使用して自動的にビルド・デプロイされます：

- **mainブランチにプッシュ**すると自動的にWASM版がビルドされます
- **GitHub Pages**に自動デプロイされます
- **WebAssembly対応ブラウザ**でゲームをプレイできます

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
