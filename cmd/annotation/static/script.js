// WebSocketサーバへの接続
var loc = window.location;
var uri = 'ws:';
if (loc.protocol === 'https:') { uri = 'wss:'; }
uri += '//' + loc.host + loc.pathname + 'ws';
const ws = new WebSocket(uri);
var file_name;
const confirmSwitch = document.getElementById('confirmSwitch');

ws.onmessage = function(event) {
	const sentData = JSON.parse(event.data);
	file_name = sentData.file_name;
	control = sentData.control;
	if (control == "FINISH") {
		drawFinish();
	} else {
		const imageURL = `http://${loc.host + loc.pathname}images/${file_name}`;

		fetchImage(imageURL)
			.then(blob => {
				const imageObjectURL = URL.createObjectURL(blob);
				loadImageAndDrawMark(imageObjectURL, sentData.point.x, sentData.point.y);
			})
			.catch(error => {
				console.error('There was a problem with the fetch operation:', error);
			});
	}
};

document.getElementById("deleteButton").addEventListener("click", function(event) {
	if (control == "FINISH") {
		return
	};
	if (confirmSwitch.value === 'on') {
		const userConfirmed = confirm("このデータを削除しますか？");
		if (!userConfirmed) {
			return
		}
	}

	const deleteData = {
		file_name: file_name,
		point: { X: -1, Y: -1 },
		control: "DELETE",
	};
	ws.send(JSON.stringify(deleteData));
})

// 画像内のクリックイベント
document.getElementById("canvas").addEventListener("click", function(event) {
	// 全てがアノテーション済みなら送信しない
	if (control == "FINISH") {
		return
	};
	// 画像内でのクリック位置を取得
	const rect = event.target.getBoundingClientRect();
	const clickX = event.clientX - rect.left;
	const clickY = event.clientY - rect.top;

	if (confirmSwitch.value === 'on') {
		const userConfirmed = confirm(`(${clickX},${clickY})送信しますか？`);
		if (!userConfirmed) {
			return
		}
	}

	// サーバに送信する座標データ
	const clickData = {
		file_name: file_name,
		point: {
			x: clickX,
			y: clickY
		},
		control: "NORMAL",
	};

	// WebSocketでサーバにクリック座標を送信
	ws.send(JSON.stringify(clickData));
});

ws.onerror = function(error) { console.error("WebSocket error:", error); };
ws.onclose = function() { console.log("WebSocket connection closed."); };
