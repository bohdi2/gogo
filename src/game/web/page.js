function init() {
    var f = document.getElementById('foo');
    var c = document.getElementById('canvas');
    var ctx = c.getContext('2d');
    var ws = new WebSocket(wsUrl);

    ws.onopen = function(unused) {
	ws.send(JSON.stringify({Type: 'REPAINT', X: 0, Y: 0}));

	c.onclick = function(e) {

	    var x;
	    var y;

	    if (e.pageX != undefined && e.pageY != undefined) {
		x = e.pageX + f.scrollLeft
		    y = e.pageY + f.scrollTop
		    }
	    else {
		x = e.clientX + f.scrollLeft;
		y = e.clientY + f.scrollTop;
	    }

	    x -= c.offsetLeft;
	    y -= c.offsetTop;

            ws.send(JSON.stringify({Type: 'CLICK', X: x, Y: y}));
	    //ws.send(JSON.stringify({Type: 'CLICK', X: e.pageX - c.offsetLeft, Y: e.pageY - c.offsetTop}));
	}
    }

    ws.onclose = function(event) {
	writeToScreen('Closed: ' + event.data);
    }

    ws.onmessage = function(event) {
	var d = JSON.parse(event.data);
	//writeToScreen('Message: ' + event.data);

	switch (d.Type) {
	case "CLEAR":
	    c.width = c.width;
	    break;

	case "DRAW": 
	    var width = d.Width;
	    var height = d.Height;
	    var faceX = d.FacePoint.X;
	    var faceY = d.FacePoint.Y;
	    
	    var x = d.GroundPoint.X;
	    var y = d.GroundPoint.Y;
	    
	    ctx.drawImage(faces, faceX, faceY, width, height, x, y, width, height);
	    break;
	}
	//ctx.strokeRect(x, y, width, height)
    }

    ws.onerror = function(event) {
	writeToScreen('<span style="color: red;">ERROR:</span> ' + event.data);
    };
}

function writeToScreen(message) { 
    var pre = document.createElement("p"); 
    pre.style.wordWrap = "break-word";
    pre.innerHTML = message; 
    //output.appendChild(pre); 
}

