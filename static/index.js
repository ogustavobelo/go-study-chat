let socket = null;
    let o = document.getElementById("output");
    let userField = document.getElementById("username");
    let messageField = document.getElementById("message");

    window.onbeforeunload = function () {
        console.log("Leaving");
        let jsonData = {};
        jsonData["action"] = "left";
        socket.send(JSON.stringify(jsonData))
    }

    document.addEventListener("DOMContentLoaded", function () {
        socket = new ReconnectingWebSocket("ws://160.20.161.56:8080/ws", null, {debug: true, reconectInterval: 3000});

        const offline = `<span class="badge bg-danger">Not connected</span>`
        const online = `<span class="badge bg-success">Connected</span>`
        let statusDiv = document.getElementById("status");

        socket.onopen = () => {
            console.log("Successfully connected");
            statusDiv.innerHTML = online;
        }

        socket.onclose = () => {
            console.log("connection closed");
            statusDiv.innerHTML = offline;
        }

        socket.onerror = error => {
            console.log("there was an error");
            statusDiv.innerHTML = offline;
        }

        socket.onmessage = msg => {
            let data = JSON.parse(msg.data);
            console.log("Action is", data.action);

            switch (data.action) {
                case "list_users":
                    let ul = document.getElementById("online_users");
                    while (ul.firstChild) ul.removeChild(ul.firstChild);

                    if (data.connected_users.length > 0) {
                        data.connected_users.forEach(function (item) {
                            let li = document.createElement("li");
                            li.appendChild(document.createTextNode(item));
                            ul.appendChild(li);
                        })
                    }
                    break;

                case "broadcast":
                    o.innerHTML = o.innerHTML + data.message + "<br>";
                    break;
            }

        }


        userField.addEventListener("change", function () {
            let jsonData = {};
            jsonData["action"] = "username";
            jsonData["username"] = this.value;
            socket.send(JSON.stringify(jsonData));
        })

        messageField.addEventListener("keydown", function (event) {
            if (event.code === "Enter") {
                if (!socket) {
                    console.log("no connection");
                    return false
                }

                if ((userField.value === "") || (messageField.value === "")) {
                    errorMessage("Fill out username and message!");
                    return false;
                } else {
                    sendMessage()
                }

                event.preventDefault();
                event.stopPropagation();
            }
        })


        document.getElementById("sendBtn").addEventListener("click", function () {
            if ((userField.value === "") || (messageField.value === "")) {
                errorMessage("Fill out username and message!");
                return false;
            } else {
                sendMessage()
            }
        })
    })

    function sendMessage() {
        let jsonData = {};
        jsonData["action"] = "broadcast";
        jsonData["username"] = userField.value;
        jsonData["message"] = messageField.value;
        socket.send(JSON.stringify(jsonData))
        messageField.value = "";
    }

    function errorMessage(msg) {
        notie.alert({
            type: 'error',
            text: msg,
        })
    }