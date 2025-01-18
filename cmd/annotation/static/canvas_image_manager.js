// canvasに対して描画機能を管理するクラス
class CanvasImageManager {
	static canvas = document.getElementById("canvas");
	static ctx = CanvasImageManager.canvas.getContext("2d");

	static getCanvas() {
		return canvas;
	}

	// 画像URLから画像を取得
	static loadImage(imageURL) {
		return new Promise((resolve, reject) => {
			const img = new Image();
			img.onload = function() {
				resolve(img);
			};
			img.onerror = function() {
				reject(new Error('Image failed to load'));
			};
			img.src = imageURL;
		});
	}

	static drawImageToCanvas(img) {
		CanvasImageManager.canvas.width = img.width;
		CanvasImageManager.canvas.height = img.height;

		// 画像をCanvasに描画
		CanvasImageManager.ctx.drawImage(img, 0, 0);
	}

	static drawMark(x, y, color) {
		const markSize = 10;
		CanvasImageManager.ctx.beginPath();
		CanvasImageManager.ctx.arc(x, y, markSize, 0, 2 * Math.PI);
		CanvasImageManager.ctx.fillStyle = color;
		CanvasImageManager.ctx.fill();
	}

	static drawString(str) {
		CanvasImageManager.canvas.width = 500;
		CanvasImageManager.canvas.height = 500;

		CanvasImageManager.ctx.font = '24px Arial';
		CanvasImageManager.ctx.fillStyle = 'black';
		CanvasImageManager.ctx.fillText(str, 50, 100);
	}

	static drawSemicircle(color) {
		const centerX = CanvasImageManager.canvas.width / 2;
		const centerY = CanvasImageManager.canvas.height;
		const radius = centerX;

		CanvasImageManager.ctx.beginPath();
		CanvasImageManager.ctx.arc(centerX, centerY, radius, Math.PI, 2 * Math.PI);
		CanvasImageManager.ctx.lineTo(centerX, centerY);
		CanvasImageManager.ctx.closePath();
		CanvasImageManager.ctx.strokeStyle = color;
		CanvasImageManager.ctx.stroke();
	}
}
