function handleCredentialResponse(response){
    console.log(response.credential);
};
var output = document.getElementById("output");
var input = document.getElementById("input");
var ws;
var userId = "";
var print = function (message) {
    var d = document.createElement("div");
    d.textContent = message;
    output.appendChild(d);
    output.scroll(0, output.scrollHeight);
};
document.getElementById("login").onclick = function (evt) {
    const xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var newUser = this.getResponseHeader("User-Id");
            if (newUser) {
                console.log("New session created");
                userId = newUser;
            }
        }
    }
    const url = "http://localhost:8080/login";
    xmlHttp.open("GET", url, true);
    if (userId != "") {
        xmlHttp.setRequestHeader("User-Id", userId);
    }
    xmlHttp.send();

    if (ws) {
        ws.send("userId:" + userId);
    }

    return false;
};
document.getElementById("open").onclick = function (evt) {
    if (ws) {
        return false;
    }
    ws = new WebSocket("ws://localhost:8080/socket");
    ws.onopen = function (evt) {
        ws.send("userId:" + userId);

        print("OPEN");
    }
    ws.onclose = function (evt) {
        print("CLOSE");
        ws = null;
    }
    ws.onmessage = function (evt) {
        print("RESPONSE: " + evt.data);
    }
    ws.onerror = function (evt) {
        print("ERROR: " + evt.data);
    }
    return false;
};
document.getElementById("send").onclick = function (evt) {
    if (!ws) {
        return false;
    }
    print("SEND: " + input.value);
    ws.send(input.value);
    return false;
};
document.getElementById("close").onclick = function (evt) {
    if (!ws) {
        return false;
    }
    ws.close();
    return false;
};