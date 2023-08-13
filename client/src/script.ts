const CANVAS_HEIGHT = 300;

const canvas = document.getElementById("canvas") as HTMLCanvasElement;
const ctx = canvas.getContext("2d")!;
const socket = new WebSocket("ws://192.168.1.9:8080/ws");

const keys = {
    isKeyADown: false,
    isKeyDDown: false
}

const player = {
    id: undefined,
    height: 40,
    width: 10,
    x: 0,
    y: 0,
    color: "red",

    draw() {
        this.y = CANVAS_HEIGHT - this.height;
        if (keys.isKeyDDown) {
            this.x += 5;
        }

        if (keys.isKeyADown) {
            this.x -= 5;
        }

        ctx.fillStyle = this.color;
        ctx.fillRect(this.x, this.y, this.width, this.height);
    },
}

async function init() {
    window.addEventListener("keydown", function(event) {
        if (event.defaultPrevented) {
            return;
        }

        if (event.code === "KeyA") {
            keys.isKeyADown = true;
        } else if (event.code === "KeyD"){
            keys.isKeyDDown = true;
        }

        event.preventDefault();
    }, true);


    window.addEventListener("keyup", function(event) {
        if (event.defaultPrevented) {
            return;
        }

        if (event.code === "KeyA") {
            keys.isKeyADown = false;
        } else if (event.code === "KeyD"){
            keys.isKeyDDown = false;
        }

        event.preventDefault();
    }, true);


    socket.onopen = () => {
        console.log("[open] Connection established");
        window.requestAnimationFrame(loop);
    };

    socket.onmessage = event => {
        console.log(`[message] Data received from server: ${event.data}`);
    };

    socket.onclose = event => {
        if (event.wasClean) {
            console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
        } else {
            // e.g. server process killed or network down
            // event.code is usually 1006 in this case
            console.log('[close] Connection died');
        }
    };

    socket.onerror = error => {
        console.log(error)
    };

}

async function loop() {

    canvas.width = window.innerWidth;
    canvas.height = 300;

    ctx.clearRect(0, 0, window.innerWidth, CANVAS_HEIGHT);

    player.draw();
    socket.send(`x:${player.x}, y:${player.y}`)

    window.requestAnimationFrame(loop);
}



init();

