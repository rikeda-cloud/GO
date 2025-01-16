// WebSocketサーバへの接続
var loc = window.location;
var uri = 'ws:';
if (loc.protocol === 'https:') { uri = 'wss:'; }
uri += '//' + loc.host + loc.pathname + 'ws';

const confirmSwitch = document.getElementById('confirmSwitch');
const deleteButton = document.getElementById('deleteButton');
const canvas = document.getElementById('canvas');

const annotationWsHandler = new AnnotationWsHandler(uri, confirmSwitch, deleteButton, canvas);
annotationWsHandler.connect();
