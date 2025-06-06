// annotationモードでのサーバとのデータのやり取りを管理するクラス
class AnnotationWsHandler {
	constructor(webSocketPath, imagePath) {
		this.webSocketUrl = `ws://${window.location.host}${window.location.pathname}${webSocketPath}`;
		this.imagePath = imagePath;

		this.confirmSwitch = document.getElementById('confirmSwitch');
		this.deleteButton = document.getElementById('deleteButton');
		this.canvas = CanvasImageManager.getCanvas();
		this.tags = new Tags();
		this.fileName = "";
		this.control = "";
		this.ws = null;
		this.imageClickHandlerBound = this.imageClickHandler.bind(this);
		this.deleteButtonHandlerBound = this.deleteButtonHandler.bind(this);
	}

	connect() {
		this.ws = new WebSocket(this.webSocketUrl);

		// WebSocketからのメッセージの受信時
		this.ws.onmessage = (event) => {
			this.handleMessage(event);
		};

		// WebSocketでのエラー発生時
		this.ws.onerror = (error) => {
			console.error("WebSocket error:", error);
		};

		// WebSocket切断時
		this.ws.onclose = () => {
			console.log("WebSocket connection closed.");
		};

		// 削除ボタンを押された時
		this.deleteButton.addEventListener("click", this.deleteButtonHandlerBound);

		// canvasに表示した画像がクリックされた時
		this.canvas.addEventListener("click", this.imageClickHandlerBound);
	}

	close() {
		this.ws.close();
		this.ws = null;
		this.deleteButton.removeEventListener("click", this.deleteButtonHandlerBound);
		this.canvas.removeEventListener("click", this.imageClickHandlerBound);
		this.tags.removeEvent();
		this.tags = null;
		console.log("WebSocket connection manually closed.");
	}

	handleMessage(event) {
		const sentData = JSON.parse(event.data);
		this.fileName = sentData.file_name;
		this.control = sentData.control;
		if (this.control === "FINISH") {
			CanvasImageManager.drawString('全てのデータがアノテーション済みです');
		} else {
			this.fetchImageAndDraw(this.fileName, sentData.point);
		}
	}

	fetchImageAndDraw(fileName, actPoint) {
		const loc = window.location;
		const imageURL = `http://${loc.host + loc.pathname}${this.imagePath}/${fileName}`;
		fetchImage(imageURL)
			.then(blob => {
				const imageObjectURL = URL.createObjectURL(blob);
				CanvasImageManager.loadImage(imageObjectURL)
					.then((img) => {
						CanvasImageManager.drawImageToCanvas(img);
						CanvasImageManager.drawMark(img.width / 2, img.height, 'white');
						CanvasImageManager.drawMark(actPoint.x, actPoint.y, 'red');
						CanvasImageManager.drawSemicircle('green');
					}).catch((error) => {
						console.err(error);
					})

			})
			.catch(error => {
				console.error('There was a problem with the fetch operation:', error);
			});
	}

	deleteButtonHandler(_) {
		if (this.control === "FINISH") {
			return
		}
		if (this.confirmSwitch.value === 'on') {
			const userConfirmed = confirm("このデータを削除しますか？");
			if (!userConfirmed) {
				return
			}
		}

		const deleteData = {
			file_name: this.fileName,
			point: { X: -1, Y: -1 },
			control: "DELETE",
			tags: this.tags.getSelectedTag(),
			user_name: sessionStorage.getItem("username"),
		};
		this.ws.send(JSON.stringify(deleteData));
	}

	imageClickHandler(event) {
		if (this.control === "FINISH") {
			return
		}

		// 画像内でのクリック位置を取得
		const rect = event.target.getBoundingClientRect();
		const clickX = event.clientX - rect.left;
		const clickY = event.clientY - rect.top;

		if (this.confirmSwitch.value === 'on') {
			const userConfirmed = confirm(`(${clickX},${clickY})送信しますか？`);
			if (!userConfirmed) {
				return
			}
		}

		const clickData = {
			file_name: this.fileName,
			point: {
				x: clickX,
				y: clickY
			},
			control: "NORMAL",
			tags: this.tags.getSelectedTag(),
			user_name: sessionStorage.getItem("username"),
		};
		this.ws.send(JSON.stringify(clickData));
	}
}


// 画像を取得する関数
async function fetchImage(url) {
	return fetch(url)
		.then(response => {
			if (!response.ok) {
				throw new Error('Network response was not ok');
			}
			return response.blob();
		});
}
