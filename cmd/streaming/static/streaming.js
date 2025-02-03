const socket = new WebSocket("ws://localhost:8000/ws/streaming");

socket.onmessage = function(event) {
	const imageBlob = event.data;
	const imgElement = document.getElementById("streaming-image");
	imgElement.src = URL.createObjectURL(imageBlob);
};

btn0 = document.getElementById("btn-converter-1")
btn1 = document.getElementById("btn-converter-2")
btn2 = document.getElementById("btn-converter-3")

btn0.addEventListener("click", function() {
	socket.send("1");
})
btn1.addEventListener("click", function() {
	socket.send("2");
})
btn2.addEventListener("click", function() {
	socket.send("3");
})
