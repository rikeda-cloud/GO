# 自動運転ミニカー

このプロジェクトは、Go言語とOpenCVを用いて構築された、自動運転ミニカーのためのデータ収集・可視化・アノテーション支援システムです。

## 📦 プロジェクト構成

```bash
.
├── cmd
│   ├── annotation/           # アノテーションツール（Webベース）
│   ├── car-data-capture/     # データ収集アプリケーション
│   ├── streaming/            # 画像ストリーミング・特徴量抽出画像の可視化機能
│   ├── data-exporter/        # 収集したデータをクラウドへ送信
│   └── ip-notify/            # IPアドレス通知
├── internal/                 # 共通ライブラリ（DB, WebSocket, config, etc）
├── images/                   # 保存された画像データ(仮置きした画像を含む)
├── configs/
│   └── config.json           # アプリケーション設定ファイル
├── go.mod / go.sum           # Go Modules 設定
├── ip_address_notify         # IPアドレス通知用設定
└── README.md
```

## ⚙️ 依存環境

このプロジェクトのビルドおよび実行には、以下の環境が必要です：
* Go 1.23.4 以上
* OpenCV（C++バインディング付き、開発パッケージ）

Fedora でのインストール例
```bash
# Go のインストール（例: 1.23.4）
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
sudo tar -xzf go1.23.4.linux-amd64.tar.gz
rm go1.23.4.linux-amd64.tar.gz
sudo mv go /usr/local/
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# OpenCV のインストール
sudo dnf install opencv opencv-devel
```

## 🔧 設定ファイル：configs/config.json
本プロジェクトでは設定ファイル configs/config.json により、動作を柔軟に制御できます。

### device_number について
これは使用するカメラのデバイス番号（例: /dev/video14）を指定します。
実行環境に接続されたカメラによって番号は異なるため、必要に応じて変更してください。

## 🛠️ ビルドと実行
プロジェクトは独立したバイナリに分かれています。

### 1️⃣ データ収集アプリケーション
カメラ画像、タイヤ角度、車体スピードを記録します。

```bash
go build ./cmd/car-data-capture/main.go
./main
```

### 2️⃣ ストリーミングビジュアライザ
カメラ画像に特徴量抽出処理を行い、Webブラウザでリアルタイム表示します。

```bash
go build ./cmd/streaming/main.go
./main
# → ブラウザで http://localhost:8000 にアクセス
```

### 3️⃣ アノテーションツール
収集済み画像と車体データを元に、学習用アノテーションデータを作成できます。

```bash
config/config.json内のoauthのパラメータを適切な値に設定

go build ./cmd/annotation/main.go
./main
# → ブラウザで http://localhost:8000 にアクセス
```

### 4️⃣ IPアドレス通知サービス
コンピュータ起動時にIPアドレスをDiscordにWebhookを用いて通知します。

1. `cmd/ip-notify/webhook.url` ファイルにDiscordのWebHookURLを記述する。
2. `go build -o ip-notify cmd/ip-notify/main.go` でバイナリを作成。
3. システムに作成したバイナリを配置。
```
chmod +x ip-notify
sudo mv ip-notify /usr/local/bin/ip-notify

# SELinux が有効な場合、制限を緩和する（必要に応じて）
sudo chcon -t bin_t /usr/local/bin/ip-notify
```
4. `ip_address_notify/ip-notify.service` をサービスとして登録
```
sudo cp ip_address_notify/ip-notify.service /etc/systemd/system/ip-notify.service
sudo systemctl enable ip-notify.service
```

### 5️⃣ クラウドへのデータ送信
アノテーション済みの車体データや画像から抽出した特徴量等をクラウドへ送信します

```bash
config/config.json内のapp::data-exporter::cloud_urlを適切な値に設定

go build ./cmd/data-exporter/main.go
./main
```
