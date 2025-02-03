const socket = new WebSocket("ws://localhost:8000/ws/streaming");

socket.onmessage = function(event) {
	const imageBlob = event.data;
	const imgElement = document.getElementById("streaming-image");
	imgElement.src = URL.createObjectURL(imageBlob);
};
