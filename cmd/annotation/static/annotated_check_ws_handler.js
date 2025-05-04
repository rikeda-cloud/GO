// checkモードでのサーバとのデータのやり取りを管理するクラス
class AnnotatedCheckWsHandler {
	constructor() {
		this.webSocketUrl = `ws://${window.location.host}${window.location.pathname}ws/check`;
		this.confirmSwitch = document.getElementById('confirmSwitch');
		this.deleteButton = document.getElementById('deleteButton');
		this.nextButton = document.getElementById('nextButton');
		this.prevButton = document.getElementById('prevButton');
		this.canvas = CanvasImageManager.getCanvas();
		this.tags = "";
		this.fileName = "";
		this.control = "";
		this.ws = null;
		this.deleteButtonHandlerBound = this.deleteButtonHandler.bind(this);
		this.imageClickHandlerBound = this.imageClickHandler.bind(this);
		this.nextButtonHandlerBound = this.nextButtonHandler.bind(this);
		this.prevButtonHandlerBound = this.prevButtonHandler.bind(this);
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

		// Nextボタンを押された時
		this.nextButton.addEventListener("click", this.nextButtonHandlerBound);
		// Prevボタンを押された時
		this.prevButton.addEventListener("click", this.prevButtonHandlerBound);
	}

	close() {
		this.ws.close();
		this.ws = null;
		this.deleteButton.removeEventListener("click", this.deleteButtonHandlerBound);
		this.canvas.removeEventListener("click", this.imageClickHandlerBound);
		this.nextButton.removeEventListener("click", this.nextButtonHandlerBound);
		this.prevButton.removeEventListener("click", this.prevButtonHandlerBound);
		console.log("WebSocket connection manually closed.");
	}

	handleMessage(event) {
		const sentData = JSON.parse(event.data);
		this.control = sentData.control;
		this.tags = sentData.tags;
		if (this.control === "FINISH") {
			CanvasImageManager.drawString('全てのデータをチェックしました');
			document.getElementById("userName").textContent = "";
		} else {
			this.fileName = sentData.file_name;
			this.fetchImageAndDraw(this.fileName, sentData.act_point, sentData.annotated_point);
			document.getElementById("userName").textContent = "Annotated by: " + sentData.user_name;
		}
	}

	fetchImageAndDraw(fileName, actPoint, annotatedPoint) {
		const loc = window.location;
		const imageURL = `http://${loc.host + loc.pathname}images/${fileName}`;
		fetchImage(imageURL)
			.then(blob => {
				const imageObjectURL = URL.createObjectURL(blob);
				CanvasImageManager.loadImage(imageObjectURL)
					.then((img) => {
						CanvasImageManager.drawImageToCanvas(img);
						CanvasImageManager.drawMark(img.width / 2, img.height, 'white');
						CanvasImageManager.drawMark(actPoint.x, actPoint.y, 'red');
						CanvasImageManager.drawMark(annotatedPoint.x, annotatedPoint.y, 'yellow');
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
			act_point: { X: -1, Y: -1 },
			annotated_point: { X: -1, Y: -1 },
			control: "DELETE",
			tags: this.tags,
			user_name: sessionStorage.getItem("username"),
		};
		this.ws.send(JSON.stringify(deleteData));
	}

	nextButtonHandler(_) {
		const toNextData = {
			file_name: this.fileName,
			act_point: { X: -1, Y: -1 },
			annotated_point: { X: -1, Y: -1 },
			control: "NEXT",
			tags: this.tags,
			user_name: sessionStorage.getItem("username"),
		};
		this.ws.send(JSON.stringify(toNextData));
	}

	prevButtonHandler(_) {
		const toPrevData = {
			file_name: this.fileName,
			act_point: { X: -1, Y: -1 },
			annotated_point: { X: -1, Y: -1 },
			control: "PREV",
			tags: this.tags,
			user_name: sessionStorage.getItem("username"),
		};
		this.ws.send(JSON.stringify(toPrevData));
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
			act_point: { x: 0, y: 0 },
			annotated_point: { x: clickX, y: clickX },
			control: "NORMAL",
			tags: this.tags,
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
