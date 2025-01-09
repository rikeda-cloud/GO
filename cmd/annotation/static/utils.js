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

// 画像をCanvasに描画し、マークを描画する関数
function loadImageAndDrawMark(imageURL, x, y) {
	const img = new Image();
	img.onload = function() {
		const canvas = document.getElementById("canvas");
		const ctx = canvas.getContext("2d");
		canvas.width = img.width;
		canvas.height = img.height;

		// 画像をCanvasに描画
		ctx.drawImage(img, 0, 0);

		// マークの描画
		drawMark(ctx, x, y, 'red');
	};
	img.src = imageURL;
}

function drawAnnotationMark(x, y) {
	const canvas = document.getElementById("canvas");
	const ctx = canvas.getContext("2d");
	drawMark(ctx, x, y, 'green');
}

// マークを描画する関数
function drawMark(ctx, x, y, color) {
	const markSize = 10; // マークのサイズ

	// マークを描画（ここでは赤い円を描画）
	ctx.beginPath();
	ctx.arc(x, y, markSize, 0, 2 * Math.PI);
	ctx.fillStyle = color;
	ctx.fill();
}

function drawFinish() {
	const canvas = document.getElementById('canvas');
	const ctx = canvas.getContext('2d');

	canvas.width = 500;
	canvas.height = 500;

	ctx.font = '24px Arial';
	ctx.fillStyle = 'black';

	ctx.fillText('全てのデータがアノテーション済みです', 50, 100);
}
