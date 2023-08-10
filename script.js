const CANVAS_HEIGHT = 300;

/** @type {HTMLCanvasElement} */
const canvas = document.getElementById("canvas");
const ctx = canvas.getContext("2d");

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
    window.requestAnimationFrame(loop);

    window.addEventListener("keydown", function(event) {
        if (event.defaultPrevented) {
            return;
        }

        console.log(event.code);
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

    const data = {
        x: 0,
        y: 0
    }
    const response = await fetch("http://localhost:8080/registerPlayer", {
        method: "POST",
        mode: "cors", 
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });

    let responseJson = await response.json();
    player.id = responseJson.id;
}

function loop() {

    canvas.width = window.innerWidth;
    canvas.height = 300;

    ctx.clearRect(0, 0, window.innerWidth, CANVAS_HEIGHT);

    player.draw();

    window.requestAnimationFrame(loop);
}

init();
