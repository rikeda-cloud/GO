// ボタンクリックでモード切り替え
document.getElementById("annotation").addEventListener("click", () => switchMode("annotation"));
document.getElementById("check").addEventListener("click", () => switchMode("check"));
document.getElementById("ai").addEventListener("click", () => switchMode("ai"));

var annotationWsHandler = null;
var remainImageWsHandler = null;

function switchMode(mode) {
	clearCurrentMode();

	if (mode === "annotation") {
		switchAnnotationMode();
	} else if (mode === "check") {
		switchCheckMode();
	} else if (mode === "ai") {
		switchAiMode();
	} else {
		console.error("Unknown mode:", mode);
	}
}

// 現在のモード要素をクリア
function clearCurrentMode() {
	const appElement = document.getElementById("app");
	appElement.innerHTML = "";

	const canvas = document.getElementById("canvas");
	canvas.getContext("2d").clearRect(0, 0, canvas.width, canvas.height);

	if (annotationWsHandler) {
		annotationWsHandler.close();
		annotationWsHandler = null;
	}
	if (remainImageWsHandler) {
		remainImageWsHandler.close();
		remainImageWsHandler = null;
	}
}

// 各モードのHTML構築
function switchAnnotationMode() {
	const appElement = document.getElementById("app");
	appElement.innerHTML = `
	<div>
		<p id = "remainImageCount" > 残り: 0</p >
	</div >
	<div>
		<label for="confirmSwitch">確認モード:</label>
		<select id="confirmSwitch">
			<option value="off">OFF</option>
			<option value="on">ON</option>
		</select>
	</div>
	<div>
		<button id="deleteButton">DEL</button>
	</div>
	<div id="tagContainer">
		<label for="tagSelect">タグ:</label>
		<select id="tagSelect">
			<option value="normal">normal</option>
			<option value="in">in</option>
			<option value="out">out</option>
			<option value="custom">custom</option>
		</select>
		<input type="text" id="customTagInput" placeholder="カスタムタグを入力" style="display: none;">
	</div>`;

	annotationWsHandler = new AnnotationWsHandler();
	annotationWsHandler.connect();
	remainImageWsHandler = new RemainImageWsHandler();
	remainImageWsHandler.connect();
}

function switchCheckMode() {
	const appElement = document.getElementById("app");
	appElement.innerHTML = `
        <div>
            <h2>確認モード</h2>
        </div>
    `;
}

function switchAiMode() {
	const appElement = document.getElementById("app");
	appElement.innerHTML = `
        <div>
            <h2>AIモード</h2>
        </div>
    `;
}
