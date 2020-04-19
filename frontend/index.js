const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

/**
 * Safari iOS pls
 */
const formSubmitPolyfill = (form, callback) => {
    if (form.requestSubmit) {
        form.requestSubmit();
    } else {
        callback();
    }
}

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

/**
 * Return (x,y) of canvas
 */
const getDimensions = (canvas) => [canvas.width, canvas.height];

/**
 * Return (x,y) midpoint of canvas
 */
const getCenter = (canvas) => getDimensions(canvas).map((dim) => dim / 2);

/**
 * Draw champion box and clear background.
 */
const drawWinner = (canvas) => {
    const ctx = canvas.getContext("2d");
    const [width, height] = getDimensions(canvas);
    const [mid_x, mid_y] = getCenter(canvas);
    const rect_width = width * .5;
    const rect_height = height * .1;
    ctx.strokeRect(mid_x - rect_width / 2, mid_y - rect_height / 2, rect_width, rect_height);
    ctx.clearRect(mid_x - rect_width / 2, mid_y - rect_height / 2, rect_width, rect_height);
}

/**
 * Draws a path in context from the current location to a point xDist * left, yDist * up away from it.
 * @param ctx RenderingContext
 * @param x, y float current location
 * @param x, y float distance away
 * @param up, left (-1|1) directions
 * @return [newNodeX, newNodeY]
 */
const drawBranchFrom = (ctx, x, y, xDist, yDist, left, up) => {
    const newX = x + (xDist * left)
    const newY = y + (yDist * up)
    ctx.lineTo(x, newY);
    ctx.lineTo(newX, newY);
    return [newX, newY];
}

const drawArtistOnCtx = (ctx, artistName, x, y) => ctx.strokeText(artistName, x, y);

/**
 * Draws paths to the terminal nodes of a round
 * @param x, y the point representation of the start of the branch
 */
const drawMatchup = (canvas, x, y, iter, left, artist1, artist2) => {
    const ctx =  canvas.getContext("2d");
    const [width, height] = getDimensions(canvas);

    const drawBranchUp = (xDist, yDist) => drawBranchFrom(ctx, x, y, xDist, yDist, left, 1);
    const drawBranchDown = (xDist, yDist) => drawBranchFrom(ctx, x, y, xDist, yDist, left, -1);
    const drawArtist = (artistName, x, y) => drawArtistOnCtx(ctx, artistName, x, y);

    ctx.beginPath();
    ctx.moveTo(x, y);
    drawArtist(
        artist1["name"],
        ...drawBranchUp(width / iter, height / iter)
    );
    ctx.moveTo(x, y);
    drawArtist(
        artist2["name"],
        ...drawBranchDown(width / iter, height / iter)
    );
    ctx.stroke();
}

const drawBracket = (canvas, artists) => {
    drawWinner(canvas);
    const [mid_x, mid_y] = getCenter(canvas);
    const groups = 4;
    const rounds = Math.log2(artists.length / groups);
    for (let group = 1; group <= groups; group++) {
        for (let round = 0; round < rounds; round++) {
            const matchup = [artists.shift(), artists.pop()];
            drawMatchup(canvas, mid_x, mid_y, 2 * round, Math.pow(-1, group), ...matchup);
        }
    }
}

window.onload = () => {
    const lStorage = window.localStorage;
    const genreList = document.getElementById("genre-list");
    const genreInput = document.getElementById("genre-input");
    const genreForm = document.getElementById("genre-form");
    const canvas = document.getElementById("bracket");

    const formSubmitAction = () => {
        fetch(encodeURI(`http://localhost:8080/artist/genre?genre_name=${genreInput.value}`))
        .then((response) => response.json())
        .then((data) => drawBracket(canvas, data.slice(0, 33)));
    }

    const createGenreList = (genreList, genres) => {
        return createGenreListWithClickEvent(genreList, genres, (e) => {
            genreInput.value = e.target.innerText;
            formSubmitPolyfill(genreForm, formSubmitAction);
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
        formSubmitAction()
    });

    genreInput.addEventListener("input", (e) => {
        const input = e.target;
        Array.from(genreList.children).forEach((item) => item.style.display= item.dataset.name.includes(input.value.toLowerCase()) ? "block" : "none");
    });
    genreInput.addEventListener("focus", (e) => {
        genreList.style.display = "block";
    });
    genreInput.addEventListener("blur", (e) => {
        sleep(150).then(() => genreList.style.display = "none");
    });


    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    drawBracket(canvas, []);

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
