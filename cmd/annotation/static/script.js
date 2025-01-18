// モード切替ボタンにイベントを登録
document.querySelectorAll("button[data-mode]").forEach((button) => {
	button.addEventListener("click", (e) => {
		const mode = e.target.dataset.mode;
		switchMode(mode);
	});
});

window.addEventListener("popstate", (event) => {
	const mode = event.state?.mode || "annotation";
	switchMode(mode, false);
});

window.addEventListener("load", () => {
	const hash = window.location.hash.replace("#", "") || "annotation";
	switchMode(hash, false);
});

var annotationWsHandler = null;
var remainImageWsHandler = null;
var annotatedCheckWsHandler = null;

function switchMode(mode, pushHistory = true) {
	if (pushHistory) {
		history.pushState({ mode }, '', `#${mode}`);
	}
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
	if (annotatedCheckWsHandler) {
		annotatedCheckWsHandler.close();
		annotatedCheckWsHandler = null;
	}
}

// 各モードのHTML構築
function switchAnnotationMode() {
	const appElement = document.getElementById("app");
	appElement.innerHTML = `
	<div id="remainContainer">
		<p id = "remainImageCount" > 残り: 0</p >
	</div >
	<div>
		<button id="deleteButton">DEL</button>
	</div>
	<div id="confirmModeContainer">
		<label for="confirmSwitch">確認モード:</label>
		<select id="confirmSwitch">
			<option value="off">OFF</option>
			<option value="on">ON</option>
		</select>
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
	<div id="nextPrevButtons">
		<button id="prevButton"><< Prev</button>
		<button id="nextButton">Next >></button>
	</div>
	<div>
		<button id="deleteButton">DEL</button>
	</div>
	<div id="confirmModeContainer">
		<label for="confirmSwitch">確認モード:</label>
		<select id="confirmSwitch">
			<option value="off">OFF</option>
			<option value="on">ON</option>
		</select>
	</div>`;

	annotatedCheckWsHandler = new AnnotatedCheckWsHandler();
	annotatedCheckWsHandler.connect();
}

function switchAiMode() {
	const appElement = document.getElementById("app");
	appElement.innerHTML = `
        <div>
            <h2>AIモード</h2>
        </div>
    `;
}
