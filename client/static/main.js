var output = document.getElementById("output");
var input = document.getElementById("input");
var ws;
var jwt;

const protocol = "http:";
const hostname = "localhost";
const port = "8080";
const server_base_url = new URL(`${protocol}//${hostname}:${port}`)

var print = function (message) {
    var d = document.createElement("div");
    d.textContent = message;
    output.appendChild(d);
    output.scroll(0, output.scrollHeight);
};

window.handleCredentialResponse = function handleCredentialResponse(response) {
    jwt = response.credential;

    if (ws) {
        ws.send("jwt:" + jwt);
    }

    const xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function () {
        if (this.readyState == 4) {
            if (this.status == 200) {
                console.log("Logged in");
            } else if (this.status == 403) {
                console.log("Log in error");
            }
        }
    }
    const url = `${server_base_url.origin}/google-login`;
    xmlHttp.open("POST", url, true);
    xmlHttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xmlHttp.send("credential="+jwt);
};

document.getElementById("open").onclick = function (evt) {
    if (ws) {
        return false;
    }
    ws = new WebSocket(`ws://${server_base_url.host}/socket`);
    ws.onopen = function (evt) {
        ws.send("jwt:" + jwt);

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