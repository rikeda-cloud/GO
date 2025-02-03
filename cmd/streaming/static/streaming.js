const socket = new WebSocket("ws://localhost:8000/ws/streaming");

socket.onmessage = function(event) {
	const imageBlob = event.data;
	const imgElement = document.getElementById("streaming-image");
	imgElement.src = URL.createObjectURL(imageBlob);
};

const frameHandlerButtons = document.querySelectorAll('button[data-value]');
frameHandlerButtons.forEach(button => {
	button.addEventListener("click", () => {
		const value = button.getAttribute('data-value');
		socket.send(value);
	})
});
