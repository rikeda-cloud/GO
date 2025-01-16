class AnnotationWsHandler {
	constructor() {
		this.webSocketUrl = `ws://${window.location.host}${window.location.pathname}ws`;
		this.confirmSwitch = document.getElementById('confirmSwitch');
		this.deleteButton = document.getElementById('deleteButton');
		this.canvas = CanvasImageManager.getCanvas();
		this.tags = new Tags();
		this.fileName = "";
		this.control = "";
		this.ws = null;
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
		this.deleteButton.addEventListener("click", this.deleteButtonHandler.bind(this));

		// canvasに表示した画像がクリックされた時
		this.canvas.addEventListener("click", this.imageClickHandler.bind(this));
	}

	close() {
		this.ws.close();
		this.ws = null;
		console.log("WebSocket connection manually closed.");
	}

	handleMessage(event) {
		const sentData = JSON.parse(event.data);
		this.fileName = sentData.file_name;
		this.control = sentData.control;
		if (this.control === "FINISH") {
			CanvasImageManager.drawString('全てのデータがアノテーション済みです');
		} else {
			this.fetchImageAndDraw(this.fileName, sentData.point.x, sentData.point.y);
		}
	}

	fetchImageAndDraw(fileName, x, y) {
		const loc = window.location;
		const imageURL = `http://${loc.host + loc.pathname}images/${fileName}`;
		fetchImage(imageURL)
			.then(blob => {
				const imageObjectURL = URL.createObjectURL(blob);
				CanvasImageManager.loadImageAndDrawMark(imageObjectURL, x, y);
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
