class CanvasImageManager {
	static canvas = document.getElementById("canvas");
	static ctx = CanvasImageManager.canvas.getContext("2d");
	static MARK_SIZE = 10;
	static FONT = '24px Arial';
	static FILL_STYLE = 'black';

	static getCanvas() {
		return canvas;
	}

	static loadImageAndDrawMark(imageURL, x, y) {
		const img = new Image();
		img.onload = function() {
			CanvasImageManager.canvas.width = img.width;
			CanvasImageManager.canvas.height = img.height;

			// 画像をCanvasに描画
			CanvasImageManager.ctx.drawImage(img, 0, 0);
			// マークの描画
			CanvasImageManager.drawMark(x, y, 'red');
		};
		img.src = imageURL;
	}

	static drawMark(x, y, color) {
		CanvasImageManager.ctx.beginPath();
		CanvasImageManager.ctx.arc(x, y, CanvasImageManager.MARK_SIZE, 0, 2 * Math.PI);
		CanvasImageManager.ctx.fillStyle = color;
		CanvasImageManager.ctx.fill();
	}

	static drawString(str) {
		CanvasImageManager.canvas.width = 500;
		CanvasImageManager.canvas.height = 500;

		CanvasImageManager.ctx.font = CanvasImageManager.FONT;
		CanvasImageManager.ctx.fillStyle = CanvasImageManager.FILL_STYLE;
		CanvasImageManager.ctx.fillText(str, 50, 100);
	}
}
