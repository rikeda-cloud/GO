class RemainImageWsHandler {
	constructor() {
		this.webSocketUrl = `ws://${window.location.host}${window.location.pathname}ws/remain-count`;
		this.remainCount = document.getElementById("remainImageCount");
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
	}

	close() {
		this.ws.close();
		this.ws = null;
		console.log("WebSocket connection manually closed.");
	}

	handleMessage(event) {
		const sentData = JSON.parse(event.data);
		const remainImageCount = sentData.current_count;
		this.remainCount.textContent = `残り: ${remainImageCount}`;
	}

	changeRemainCountColor(count) {
		// 画像枚数に応じたクラスを動的に変更
		if (count > 5000) {
			this.remainCount.classList.add("high");
			this.remainCount.classList.remove("medium", "low");
		} else if (count > 500) {
			this.remainCount.classList.add("medium");
			this.remainCount.classList.remove("high", "low");
		} else {
			this.remainCount.classList.add("low");
			this.remainCount.classList.remove("high", "medium");
		}
	}
}

