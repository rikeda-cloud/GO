const socket = new WebSocket("ws://localhost:8000/ws/streaming");

socket.onmessage = function(event) {
	const imageBlob = event.data;
	const imgElement = document.getElementById("streaming-image");
	imgElement.src = URL.createObjectURL(imageBlob);
};

const buttonConfigs = [
	{ label: "Hough", value: "1" },
	{ label: "Gray", value: "2" },
	{ label: "Canny", value: "3" },
]

const frameHandlerDiv = document.getElementById("frame-handler-buttons");

buttonConfigs.forEach(config => {
	const button = document.createElement('button');
	button.setAttribute('data-value', config.value);
	button.innerText = config.label;

	button.addEventListener("click", () => {
		const value = button.getAttribute('data-value');
		socket.send(value);
	});
	frameHandlerDiv.appendChild(button);
});
