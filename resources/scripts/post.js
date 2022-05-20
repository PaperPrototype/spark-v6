document.addEventListener("alpine:init", function(event) {
	Alpine.store("post", {
		visible: false,
		state: 'view', // 'new' 'edit'
		commentsVisible: false,
		postID: 0,
		newestCommentDate: "",
		params: {},
	});

	const params = new URLSearchParams(window.location.search)

	if (params.has("post_id"))
	{
		loadPost(params.get("post_id"));
	}
});

let postCommentForm;

// aborting comments long polling
let abortCommentsController = new AbortController();

document.addEventListener("DOMContentLoaded", function(event) {
	let postCommentForm = document.getElementById("postCommentForm");

	postCommentForm.addEventListener("submit", sendPostComment, false);
});

async function postNewPost() {
	let params = Alpine.store("post").params;

	let postRawMarkdown = document.getElementById("postRawMarkdown").value;
	
	if (postRawMarkdown.value === "")
	{
		SendMessage("Can't make an empty post!");
		return;
	}

	var esc = encodeURIComponent; // esc now can be used in place of the encodeURIComponent function
	var query = Object.keys(params)
		.map(function(k) {return esc(k) + '=' + esc(params[k]);})
		.join('&');

	console.log("query params for creating new post are:", query);

	let formData = new FormData();
	formData.append("markdown", postRawMarkdown);

	let response = await fetch("/api/posts/new?"+query, {
		method: "POST",
		body: formData,
	});

	if (!response.ok)
	{
		SendMessage("Error creating post.");
		return;
	}

	// response will give back the postID
	response.json()
	.then(function(json) {
		let post = json;

		console.log("post json is:" + post);
	
		loadPost(post.ID);
	});
}

function newPostView(params={}) {
	console.log("newPostView");

	Alpine.store('post').visible = true;
	Alpine.store('post').state = 'new';

	Alpine.store("post").params = params;

	document.getElementById("postRawMarkdown").value = "";
}

function editPostView() {
	console.log("editPostView");

	Alpine.store('post').visible = true;
	Alpine.store('post').state = 'edit';
}

function openPostView(postID) {
	Alpine.store("post").postID = postID;

	Alpine.store('post').state = 'view';

	console.log("loading post with ID of " + postID + "...")

	Alpine.store("post").visible = true;

	const params = new URLSearchParams(window.location.search)
	params.set("post_id", postID);
	window.history.replaceState(null, "", window.location.pathname + "?" + params.toString());
}

function loadPost(postID) {
	openPostView(postID);

	setPostToLoading();

	fetch("/api/posts/"+postID, {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			throw new Error("Response for loadPost was not ok");
		}

		return resp.json();
	})
	.then(function(json) {

		setPost(json);

		// reset this so that when we click a different post the loadCommentsLongPolling function will know to
		// load initial comments and not comments after a specific date
		Alpine.store("post").newestCommentDate = "";

		resetComments();
	})
	.catch(function(err) {
		console.error(err);
	});
}

// load the plaintext markdown of a post into the editing textarea
function loadPostPlaintext() {
	let postID = Alpine.store("post").postID;

	let postRawMarkdown = document.getElementById("postRawMarkdown");

	postRawMarkdown.value = "Loading...";

	fetch("/api/posts/"+postID +"/plaintext", {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			SendMessage("Error getting post plaintext");
			return
		}

		return resp.json();
	})
	.then(function(json) {
		console.log("post plaintext json is:" + json);

		postRawMarkdown.value = json.Markdown;
	})
	.catch(function(err) {
		console.error(err);
	})
}

// update a post
function postUpdatePost() {
	let postID = Alpine.store("post").postID;

	let postRawMarkdown = document.getElementById("postRawMarkdown");

	if (postRawMarkdown.value === "")
	{
		SendMessage("Can't make an empty post!");
		return;
	}

	let formData = new FormData();
	formData.append("markdown", postRawMarkdown.value);

	fetch("/api/posts/" + postID + "/update", {
		method: "POST",
		body: formData,
	})
	.then(function(resp) {
		if (!resp.ok) {
			SendMessage("Error updateing post.")
			return;
		}

		SendMessage("Successfully updated post.")

		// load post again
		loadPost(postID);
	})
	.catch(function(err) {
		console.error(err);
	})
}

function closePostView() {
	Alpine.store("post").visible = false;
	Alpine.store("post").commentsVisible = false;

	window.history.replaceState(null, "", window.location.pathname);
}

function openPostComments() {
	console.log("opening comments...");
	Alpine.store("post").commentsVisible = true;

	let postID = Alpine.store("post").postID;

	loadPostComments(postID);
}

function closePostComments() {
	console.log("closing comments...");
	Alpine.store("post").commentsVisible = false;
}

function sendPostComment(event) {
	if (event) {
		event.preventDefault();
	}

	let postCommentInput = document.getElementById("postCommentInput");
	if (postCommentInput.value === "") {
		SendMessage("Can't send an empty message");
		return;
	}

	let postID = Alpine.store("post").postID;

	let formData = new FormData();
	formData.append("markdown", postCommentInput.value);

	fetch("/api/posts/" + postID + "/comment/", {
		method: "POST",
		body: formData,
	})
	.then(function(resp) {
		if (!resp.ok) {
			SendMessage("Error sending comment: " + resp.statusText);
			return;
		}

		postCommentInput.value = "";
	})
	.catch(function(err) {
		console.error(err);
	})
}

// abort any previous polling requests for comments
function resetComments() {
	abortCommentsController.abort();
	abortCommentsController = new AbortController();
}

// long polling taken from here
// https://javascript.info/long-polling
async function loadPostComments(postID) {
	resetComments();
	
	if (Alpine.store("post").visible === false) {
		// stop long polling
		return;
	}

	console.log("loading comments...");

	let lastCommentDate = Alpine.store("post").newestCommentDate;

	let resp;
	try {
		resp = await fetch("/api/posts/"+postID+"/comments?newest="+lastCommentDate, {
			method: "GET",
			signal: abortCommentsController.signal,
		});
	} catch (err) {
		if (err.name === "AbortError") {
			console.log("fetch comments was aborted...");
			return
		}

		console.log("error getting comments...:", err);

		// wait 3 seconds
		await new Promise(resolve => setTimeout(resolve, 3000));

		// then long poll comments again
		await loadPostComments(postID);

		return
	}

	let postCommentsMount = document.getElementById("postComments")

	// if first time loading, then clear inner html
	if (lastCommentDate === "") {
		postCommentsMount.innerHTML = "";
	}

	let json = await resp.json();

	let lastComment;
	for (let i = 0; i < json.length; i++) {
		let comment = document.createElement("div");
		comment.innerHTML =
		`<div class="c-bold" style="margin-top:0.5rem; font-weight:600;" href="/` + json[i].User.Username + `" external> @` + json[i].User.Username  + ` <i class="fa-solid fa-arrow-up-right-from-square"></i></div>` +
		`<p>` + json[i].Markdown + `</p>`;
		comment.style = "border-top:0.08rem solid var(--c-light); margin-left:1rem; margin-right:1rem;";

		postCommentsMount.append(comment);

		if (i === json.length - 1) {
			Alpine.store("post").newestCommentDate = json[i].CreatedAt;

			console.log("setting newestCommentDate");

			// set last comment
			lastComment = comment;
		}
	}

	// if there is a last comment
	if (lastComment !== undefined) {
		lastComment.scrollIntoView();

		// make username links work
		convertHrefs(postCommentsMount);
	}

	// wait 2 seconds
	await new Promise(resolve => setTimeout(resolve, 1000));
	await loadPostComments(postID);
}

/* HELPERS */
// takes in a posts json (with all properties preloaded)
function setPost(post) {
	let postMarkdown = document.getElementById("postMarkdown");

	let postAuthor = document.getElementById("postAuthor");
	postAuthor.innerHTML =
	`<div class="c-bold" style="margin-top:0.5rem; font-weight:600;" href="/` + post.User.Username + `" external> @` 
		+ post.User.Username  + 
		`<i class="fa-solid fa-arrow-up-right-from-square"></i>` + 
	`</div>`;

	convertHrefs(postAuthor);

	postMarkdown.innerHTML = post.Markdown;

	convertHrefs(postMarkdown);
}

function setPostToLoading() {
	let postMarkdown = document.getElementById("postMarkdown");

	let postAuthor = document.getElementById("postAuthor");
	postAuthor.innerHTML = "";

	postMarkdown.innerHTML = "Loading...";
}