var loc = window.location;
var remain_count_uri = 'ws:';
if (loc.protocol === 'https:') { remain_count_uri = 'wss:'; }
remain_count_uri += '//' + loc.host + loc.pathname + 'ws/remain-count';
const reman_count_ws = new WebSocket(remain_count_uri);

reman_count_ws.onmessage = function(event) {
	const sentData = JSON.parse(event.data);
	const remainImageCount = sentData.current_count;
	const remainCountElement = document.getElementById("remainImageCount");
	remainCountElement.textContent = `残り: ${remainImageCount}`;

	// 画像枚数に応じたクラスを動的に変更
	if (remainImageCount > 5000) {
		remainCountElement.classList.add("high");
		remainCountElement.classList.remove("medium", "low");
	} else if (remainImageCount > 500) {
		remainCountElement.classList.add("medium");
		remainCountElement.classList.remove("high", "low");
	} else {
		remainCountElement.classList.add("low");
		remainCountElement.classList.remove("high", "medium");
	}
};

reman_count_ws.onerror = function(error) { console.error("WebSocket error:", error); };
reman_count_ws.onclose = function() { console.log("WebSocket connection closed."); };
