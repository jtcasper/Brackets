const sleep = (ms) => new Promise(resolve => setTimeout(resolve, ms));

const getGenres = (lStorage) => {
    return new Promise((resolve) => {
        const genres = JSON.parse(lStorage.getItem("genres"))
        if (genres === null) {
            fetch("http://api.brackets.jacobcasper.com/genre")
                .then((response) => response.text())
                .then((text) => {
                    lStorage.setItem("genres", text);
                    resolve(JSON.parse(text));
                })
        } else {
            resolve(genres);
        }
    });
}

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
        const option = document.createElement("option");
        option.setAttribute("data-name", name);
        option.addEventListener("click", callback);
        option.appendChild(document.createTextNode(name));
        return option;
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

const getRectangleDimensionsUnbound = (canvas, xScale, yScale) => {
    [width, height] = getDimensions(canvas);
    return [width * xScale, height * yScale];
}

const getRectangleDimensions = (canvas) => getRectangleDimensionsUnbound(canvas, .5, .1);

/**
 * Draw champion box and clear background.
 */
const drawWinner = (canvas) => {
    const ctx = canvas.getContext("2d");
    const [width, height] = getDimensions(canvas);
    const [mid_x, mid_y] = getCenter(canvas);
    const [rect_width, rect_height] = getRectangleDimensions(canvas);
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
    ctx.stroke();
    return [newX, newY];
}

const drawArtistOnCtx = (ctx, artistName, x, y) => {
    ctx.font = "20px sans serif";
    ctx.strokeText(artistName, x, y);
}

/**
 * Draws paths to the terminal nodes of a round
 * @param x, y the point representation of the start of the branch
 * @param baseCallback a callback that needs the terminal x,y context
 */
const drawMatchup = (canvas, x, y, iter, maxIter, left, artists, baseCallback) => {
    if (iter === maxIter) {
        return baseCallback(x, y);
    }
    const ctx =  canvas.getContext("2d");
    ctx.direction = left === -1 ? "ltr" : "rtl";
    const [width, height] = getDimensions(canvas);

    const drawBranchUp = (xDist, yDist) => drawBranchFrom(ctx, x, y, xDist, yDist, left, -1);
    const drawBranchDown = (xDist, yDist) => drawBranchFrom(ctx, x, y, xDist, yDist, left, 1);
    const drawArtist = (artistName, x, y) => drawArtistOnCtx(ctx, artistName, x, y);
    const drawArtist1 = (x, y) => {
        const artist1 = artists.shift();
        if (!artist1) {
            return;
        }
        drawArtist(artist1["name"], x, y);
    }

    const drawArtist2 = (x, y) => {
        const artist2 = artists.pop();
        if (!artist2) {
            return;
        }
        drawArtist(artist2["name"], x, y);
    }

    ctx.beginPath();
    ctx.moveTo(x, y);
    const branchDistances = [width / (5 * (iter + 1)), height / (9 * (iter + 1))];
    drawMatchup(
        canvas,
        ...drawBranchUp(...branchDistances),
        iter + 1,
        maxIter,
        left,
        artists,
        drawArtist1,
    );

    ctx.moveTo(x, y);
    drawMatchup(
        canvas,
        ...drawBranchDown(...branchDistances),
        iter + 1,
        maxIter,
        left,
        artists,
        drawArtist2,
    );
}

const drawBracket = (canvas, artists, genre) => {
    const context = canvas.getContext("2d");
    context.clearRect(0, 0, canvas.width, canvas.height);
    drawWinner(canvas);
    const [mid_x, mid_y] = getCenter(canvas);
    context.font = '28px sans serif'
    context.textAlign = "center";
    context.strokeText(genre.toUpperCase(), mid_x, 40);
    context.textAlign = "start";
    const [rect_width, rect_height] = getRectangleDimensions(canvas);
    const groups = 4;
    const rounds = Math.max(
        Math.floor(Math.log2(artists.length / groups)),
        1
    );
    for (let group = 1; group <= groups; group++) {
        if (artists.length === 0) {
            break;
        }
            drawMatchup(
                canvas,
                mid_x + (Math.pow(-1, group) * (rect_width / 6)),
                mid_y + (Math.pow(-1, Math.floor(group / 2)) * (rect_height * 2.5)),
                0,
                rounds,
                Math.pow(-1, group),
                artists,
                (x, y) => console.log("hello")
            );
    }
}

window.onload = () => {
    const lStorage = window.localStorage;
    const genreList = document.getElementById("genre-list");
    const genreInput = document.getElementById("genre-input");
    const genreForm = document.getElementById("genre-form");
    const canvas = document.getElementById("bracket");

    const genreFormSubmitPolyfill = () => formSubmitPolyfill(genreList, formSubmitAction)

    const createGenreList = (genreList, genres) => {
        return createGenreListWithClickEvent(genreList, genres, (e) => {
            genreInput.value = e.target.innerText;
            genreFormSubmitPolyfill()
        })
    }

    const genres = getGenres(lStorage)
        .then((genres) => createGenreList(genreList, genres));

    const formSubmitAction = () => {
        fetch(encodeURI(`http://api.brackets.jacobcasper.com/artist/genre?genre_name=${genreInput.value}`))
        .then((response) => response.json())
        .then((data) => drawBracket(canvas, data.slice(0, 33), genreInput.value));
    }

    genreForm.addEventListener("submit", (e) => {
        e.preventDefault();
        formSubmitAction()
    });

    genreInput.addEventListener("change", (e) => {
        genreFormSubmitPolyfill()
    });


    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    drawBracket(canvas, [], "");

    const bgImg = new Image();
    bgImg.onload = () => {
        const ctx = canvas.getContext("2d");
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
