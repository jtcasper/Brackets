const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

const createGenreListWithClickEvent = (genreList, genres, callback) => {
    genreList.append(...genres.map((genre) => {
        const name = genre["name"];
        const li = document.createElement("li");
        li.setAttribute("data-name", name);
        li.addEventListener("click", callback);
        li.appendChild(document.createTextNode(name));
        return li;
    }));
}

const relativeToCtx = (x1, x2, y1, y2, ctx) => {
    const newX = x1 + x2;
    const newY = y1 + y2;
    ctx.lineTo(newX, newY);
    return [newX, newY];
}

const drawBracket = (canvas, ...artists) => {
    const ctx = canvas.getContext("2d");
    const height = canvas.height;
    const width = canvas.width;
    const mid_x = width / 2;
    const mid_y = height / 2;
    const rect_width = width * .5;
    const rect_height = height * .1;
    ctx.strokeRect(mid_x - rect_width / 2, mid_y - rect_height / 2, rect_width, rect_height);
    ctx.clearRect(mid_x - rect_width / 2, mid_y - rect_height / 2, rect_width, rect_height);

    relativeTo = (x1, x2, y1, y2) => relativeToCtx(x1, x2, y1, y2, ctx);
    let [x, y] = [mid_x, mid_y - rect_height / 2];
    ctx.beginPath();
    ctx.moveTo(mid_x, mid_y - rect_height / 2);
    [x, y] = relativeTo(x, 0, y, -height * .2);
    [x, y] = relativeTo(x, -width * .05, y, 0);
    const [branch1x, branch1y] = [x, y];
    [x, y] = relativeTo(x, 0, y, -height * .05);
    [x, y] = relativeTo(x, -width * .10, y, 0);
    const [branch2x, branch2y] = [x, y];
    [x, y] = relativeTo(x, 0, y, -height * .05);
    [x, y] = relativeTo(x, -width * .15, y, 0);
    const [branch3x, branch3y] = [x, y];
    [x, y] = relativeTo(x, 0, y, -height * .05);
    [x, y] = relativeTo(x, -width * .15, y, 0);
    ctx.moveTo(branch3x, branch3y);
    [x, y] = relativeTo(branch3x, 0, branch3y, height * .05);
    [x, y] = relativeTo(x, -width * .15, y, 0);
    ctx.stroke();
}

window.onload = () => {
    const lStorage = window.localStorage;
    const genreList = document.getElementById("genre-list");
    const genreInput = document.getElementById("genre-input");
    const genreForm = document.getElementById("genre-form");
    const canvas = document.getElementById("bracket");

    const createGenreList = (genreList, genres) => {
        return createGenreListWithClickEvent(genreList, genres, (e) => {
            genreInput.value = e.target.innerText;
            genreForm.requestSubmit();
        })
    }
    let genres = JSON.parse(lStorage.getItem("genres"))
    if (genres === null) {
        fetch("http://localhost:8080/genre")
            .then((response) => response.text())
            .then((text) => {
                window.localStorage.setItem("genres", text);
                genres = JSON.parse(text);
                createGenreList(genreList, genres);
            });
    } else {
        createGenreList(genreList, genres);
    }

    genreForm.addEventListener("submit", (e) => {
        e.preventDefault();
        fetch(encodeURI(`http://localhost:8080/artist/genre?genre_name=${genreInput.value}`))
        .then((response) => response.json())
        .then((data) => console.log(data));
        drawBracket(canvas, ...[])
    });

    genreInput.addEventListener("input", (e) => {
        const input = e.target;
        for (const item of genreList.children) {
            item.style.display = item.dataset.name.includes(input.value.toLowerCase()) ? "block" : "none";
        }
    });
    genreInput.addEventListener("focus", (e) => {
        genreList.style.display = "block";
    });
    genreInput.addEventListener("blur", (e) => {
        sleep(150).then(() => genreList.style.display = "none");
    });


    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    drawBracket(canvas);

    const bgImg = new Image();
    bgImg.onload = () => {
        const ctx = canvas.getContext("2d");
        //ctx.mozImageSmoothingEnabled = false;
        //ctx.webkitImageSmoothingEnabled = false;
        //ctx.msImageSmoothingEnabled = false;
        //ctx.imageSmoothingEnabled = false;
        ctx.drawImage(
            bgImg, 
            0, 
            0, 
            (canvas.width / bgImg.width) * bgImg.width, 
            (canvas.height / bgImg.height) * bgImg.height
        );
    }

    const upload = document.getElementById("image-upload");
    upload.addEventListener("change", (e) => {
        bgImg.src = URL.createObjectURL(e.target.files[0]);
    });
}
