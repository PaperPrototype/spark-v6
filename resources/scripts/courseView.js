
let courseNavTop;
let courseMain;
let courseMenu;
let usingGithub = false;

// offset from the top of the screen
let ogTopOffset;

// form for sending comments on a post
let postCommentsForm;

let abortCommentsController = new AbortController();
let abortMessagesController = new AbortController();

document.addEventListener("DOMContentLoaded", function(event) {
	courseNavTop = document.getElementById("courseNavTop");
	courseMain = document.getElementById("courseMain");
	courseMenu = document.getElementById("courseMenu");
	ogTopOffset = courseNavTop.offsetTop;

	menuFollowScroll();

	window.onscroll = function(){
		menuFollowScroll();
	}

	postCommentsForm = document.getElementById("postCommentsForm");
	postCommentsForm.addEventListener("submit", sendComment, false)
});

// toggle nav menu
function toggleMenu() {
	// if the menu is available
	if (Alpine.store("courseView").menuAvailable) {
		if (Alpine.store("courseView").menu === 'chat' && Alpine.store("courseView").menuOpen === true) {
			Alpine.store("courseView").menu = 'sections';
		} else {
			Alpine.store("courseView").menuOpen = ! Alpine.store("courseView").menuOpen;
			Alpine.store("courseView").menu = 'sections';
		}
	// if the menu is not available
	} else {
		Alpine.store("courseView").menuOpen = false;
	}
}

// toggle chat view
function toggleChat() {
	// if viewing sections and menu is already open
	if (Alpine.store("courseView").menu === 'sections' && Alpine.store("courseView").menuOpen === true) {
		Alpine.store("courseView").menu = 'chat';
	} else {
		Alpine.store("courseView").menuOpen = ! Alpine.store("courseView").menuOpen;
		Alpine.store("courseView").menu = 'chat';
	}

	console.log("scrolling to bottom");
	let channelMount = document.getElementById("channelMount");
	channelMount.scrollTo(0, 100);
}

// show posts view
function viewPosts() {
	resetComments();
	
	Alpine.store("courseView").view = 'posts';
	Alpine.store("courseView").menuAvailable = false;
	Alpine.store("courseView").menuOpen = false;
	Alpine.store("courseView").viewingPost = false;
	Alpine.store("courseView").editingContent = false;
	Alpine.store("courseView").viewingComments = false;
}

// show contents/lecture/section view
function viewContents() {
	resetComments();
	resetURL();
	
	Alpine.store('courseView').menuOpen = false;
	Alpine.store("courseView").view = 'contents';
	Alpine.store("courseView").menuAvailable = true;
	Alpine.store("courseView").allowNewPost = true;
	Alpine.store("courseView").editingContent = false;
	Alpine.store("courseView").viewingComments = false;
}

function loadPosts(versionID) {
	console.log("loading posts...");

	resetComments();

	let coursePostsMount = document.getElementById("coursePostsMount");
	coursePostsMount.innerHTML = "";

	fetch("/api/version/"+ versionID +"/posts/portfolio", {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			SendMessage("Error loading post.");
			throw new Error("Response for loadPosts was not ok");
		}

		return resp.json();
	})
	.then(function(json) {
		let postsTitle = document.createElement("h3");
		postsTitle.innerText = "Portfolio Posts";
		postsTitle.style = "margin-bottom:0.5rem;";

		let posts = document.createElement("div");
		posts.setAttribute("class", "post-cards");
		posts.setAttribute("style", "margin-bottom:1rem;");
		posts.innerHTML = "";

		for (let i = 0; i < json.length; i++) {
			let post = document.createElement("div");
			post.classList.add("post-card");
			post.setAttribute("postID", json[i].ID)
			posts.appendChild(post);

			let elements = document.createElement("div");
			elements.innerHTML = json[i].Markdown;
			let nodes = elements.querySelectorAll("*");

			let title = "";
			for (let nodeIndex = 0; nodeIndex < 1; nodeIndex++) {
				title = title + " " + nodes[nodeIndex].innerText;
			}

			let h4 = document.createElement("h4");
			h4.innerText = title;
			post.appendChild(h4);

			post.addEventListener("click", function(event) {
				console.log("clicked on post with ID of:", this.getAttribute("postID"));
				loadPost(this.getAttribute("postID"));
			});
		}

		let seeAll = document.createElement("div");
		seeAll.classList.add("text-center");
		seeAll.classList.add("post-card");
		seeAll.classList.add("bg-code");
		seeAll.classList.add("hoverable");
		posts.appendChild(seeAll);

		let h3 = document.createElement("h3");
		h3.innerText = `See all`;
		seeAll.appendChild(h3);

		let arrow = document.createElement("div");
		arrow.innerHTML = `<i class="fa-solid fa-arrow-right-long"></i>`;
		seeAll.appendChild(arrow);

		seeAll.addEventListener("click", function(event) {
			console.log("TODO create see all posts view");
		});

		let postsWrapper = document.createElement("div");
		postsWrapper.setAttribute("style", "margin-bottom:1rem;");
		postsWrapper.append(postsTitle);
		postsWrapper.append(posts);
		postsWrapper.classList.add("pad-sides-5");

		coursePostsMount.append(postsWrapper);
	})
	.catch(function(err) {
		console.error(err);
	});
}

function viewEditPost() {
	let editPostTextInput = document.getElementById("editPostTextInput");
	let postID = Alpine.store("courseView").postID;

	fetch("/api/posts/" + postID + "/plaintext", {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			SendMessage("Error loading post.");
			throw new Error("Response for loadPosts was not ok");
		}

		return resp.json();
	})
	.then(function(json) {
		editPostTextInput.value = json.Markdown;

		Alpine.store("courseView").editingPost = true;
	})
	.catch(function(err) {
		SendMessage("Failed to get post!");
		console.error(err);
	});
}

function updatePost() {
	let postID = Alpine.store("courseView").postID;
	let editPostTextInput = document.getElementById("editPostTextInput");

	let formData = new FormData();
	formData.append("markdown", editPostTextInput.value);

	fetch("/api/posts/" + postID + "/update", {
		method: "POST",
		body: formData,
	})
	.then(function(resp) {
		if (!resp.ok) {
			throw new Error("Error updating post!")
		}

		SendMessage("Updated post!");

		// close editing window
		Alpine.store("courseView").editingPost = false;

		// load post contents again
		loadPost(postID);

		return resp.text();
	})
	.then(function(message) {

	})
	.catch(function(err) {
		SendMessage("Failed to update post!");
		console.error(err)
	});
}

function loadPost(postID) {
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
		let postMount = document.getElementById("postMount");
		let postMountUser = document.getElementById("postMountUser");

		// if save for if user decides to edit the post
		Alpine.store("courseView").postID = json.ID;

		Alpine.store("courseView").editingPost = false;
		Alpine.store("courseView").viewingPost = true;

		postMount.innerHTML = json.Markdown;
		postMountUser.innerText = "@" + json.User.Username;

		Alpine.store("courseView").postID = postID;

		// reset this so that when we click a different post the loadCommentsLongPolling function will know to
		// load initial comments and not comments after a specific date
		Alpine.store("courseView").newestCommentDate = "";

		resetComments();

		loadPostComments(postID);

		// convert any elems with a link to open in new page if they have external
		convertHrefs(postMount);

		// get current course URL
		let courseURL = document.getElementById("courseURL").innerText;

		if (Alpine.store("sections").current === "") {
			// change location of window
			window.history.replaceState("", "", courseURL + "?post_id=" + postID);
		} else {
			// change location of window
			window.history.replaceState("", "", courseURL + "/" + Alpine.store("sections").current + "?post_id=" + postID);
		}
	})
	.catch(function(err) {
		console.error(err);
	});
}

function menuFollowScroll() {
	// if top offset is greater than window's Y-scroll offset
	// if scrolled past topbar
	if (window.scrollY >= ogTopOffset) {
		// scrolled down

		// make menu follow
		courseMenu.classList.add("course-menu-fixed");
		courseMenu.classList.remove("course-menu-normal");
		
		// set top nav to stick to top
		courseNavTop.classList.add("course-top-nav-fixed");

		// prevent course contents from jumping into place when top nav disapears
		courseMain.style.marginTop = courseNavTop.getBoundingClientRect().height + "px";
	} else {
		// top

		// stop menu from sticking to the top as a nav bar
		courseNavTop.classList.remove("course-top-nav-fixed");

		// reset courseMain's margin
		courseMain.style.marginTop = "0";

		// make menu still follow
		courseMenu.style.top = ((ogTopOffset + courseNavTop.getBoundingClientRect().height) - window.scrollY) + "px";

		// remove menu following
		courseMenu.classList.remove("course-menu-fixed");
		courseMenu.classList.add("course-menu-normal");
	}
}

// load a github section or an upload based section
function loadSection(id) {
	console.log("attempting to load...");
	if (Alpine.store("courseView").usingGithub) {
		loadGithubSection(id);
	} else {
		loadUploadSection(id);
	}
}

function sendComment(event) {
	if (event) {
		event.preventDefault();
	}

	let postCommentsSend = document.getElementById("postCommentsSend");
	let postID = Alpine.store("courseView").postID;

	if (postCommentsSend.value === "") {
		SendMessage("Can't send an empty message");
		return
	}

	let formData = new FormData();
	formData.append("markdown", postCommentsSend.value);

	let versionID = document.getElementById("versionID").innerText;

	fetch("/api/version/" + versionID + "/posts/" + postID + "/comment", {
		method: "POST",
		body: formData,
	})
	.then(function(resp) {
		if (!resp.ok) {
			throw new Error("Response was not ok");
		}

		postCommentsSend.value = "";
	})
	.catch(function(err) {
		SendMessage("Error sending message")
		console.error(err);
	});
}

// abort any previous polling requests in for comments
function resetComments() {
	abortCommentsController.abort();
	abortCommentsController = new AbortController();
}

// long polling taken from here
// https://javascript.info/long-polling
async function loadPostComments(postID) {
	resetComments();
	
	if (Alpine.store("courseView").viewingPost === false) {
		// stop long polling
		return
	}

	console.log("loading comments...");

	let lastCommentDate = Alpine.store("courseView").newestCommentDate;

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

	let postCommentsMount = document.getElementById("postCommentsMount")

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
			Alpine.store("courseView").newestCommentDate = json[i].CreatedAt;

			// set last comment
			lastComment = comment;
		}
	}

	// if there is a last comment
	if (lastComment !== undefined) {
		lastComment.scrollIntoView();

		// make links work?
		// not sure why this is here
		convertHrefs(document);
	}

	// wait 1 seconds
	await new Promise(resolve => setTimeout(resolve, 1000));
	await loadPostComments(postID);
}

async function sendMessage() {

	let chatTextarea = document.getElementById("chatTextarea");

	if (chatTextarea.value === "") {
		SendMessage("Can't send an empty message!");
		return
	}

	let versionID = document.getElementById("versionID").innerText;
	let channelID = Alpine.store("courseView").channelID;

	let formData = new FormData();
	formData.append("markdown", chatTextarea.value);

	try {
		console.log("sending message...");
		await fetch("/api/version/" + versionID + "/channel/" + channelID + "/message", {
			method: "POST",
			body: formData,
		});
	} catch {
		SendMessage("Error sending message.");
	}

	// minimize textarea after sending a message
	chatResetTextarea();

	chatTextarea.value = "";
}

// abort any previous polling requests for comments
function resetMessages() {
	abortMessagesController.abort();
	abortMessagesController = new AbortController();
}

function chatResetTextarea() {
	let textareaElement = document.getElementById("chatTextarea");
	// set the height to 0 in case of it has to be shrinked
	textareaElement.style.height = "1.2rem";
}

function loadChannel(channelID, channelName) {
	resetMessages();

	let channelTitle = document.getElementById("channelTitle");
	channelTitle.innerText = channelName;

	// new channel loaded, clear the old comments and set this as new channel
	Alpine.store("courseView").channelID = channelID;
	channelMount.innerHTML = "";
	Alpine.store("courseView").newestMessageDate = "";

	if (channelID === 0) {
		console.log("no channels to load messages from.");
		return
	}

	loadChannelMessages(channelID);
}

async function loadChannelMessages(channelID) {
	let newest = Alpine.store("courseView").newestMessageDate;

	console.log("loading messages...");

	let channelMount = document.getElementById("channelMount");

	// if this is the first time we load we should clear the html
	if (newest === "") {
		channelMount.innerHTML = "";
	}

	let resp;
	try {
		resp = await fetch("/api/channels/"+channelID+"?newest="+newest, {
			method: "GET",
			signal: abortMessagesController.signal,
		})
	} catch (err) {
		if (err.name === "AbortError") {
			console.log("fetch messages was aborted...");
			return
		}

		console.log("error getting messages...:", err);

		// wait 3 seconds
		await new Promise(resolve => setTimeout(resolve, 3000));

		// then long poll messages again
		await loadChannelMessages(channelID);
		return
	}

	let json = await resp.json();
	
	let lastMessage;
	for (let i = 0; i < json.Messages.length; i++) {
		let message = document.createElement("div");
		message.innerHTML = `<div class="c-bold" style="margin-top:0.5rem; font-weight:600;" href="/` + json.Messages[i].User.Username + `" external> @` + json.Messages[i].User.Username + ` <i class="fa-solid fa-arrow-up-right-from-square"></i></div>` +
		`<div markdown>` + json.Messages[i].Markdown + `</div>`;
		message.style = "border-top:0.08rem solid var(--c-light); margin-left:1rem; margin-right:1rem;";

		channelMount.append(message);

		if (i === json.Messages.length - 1) {
			Alpine.store("courseView").newestMessageDate = json.Messages[i].CreatedAt;

			console.log("messages delivered...");

			// set last comment
			lastMessage = message;
		}
	}

	// if there is a last comment
	if (lastMessage !== undefined) {
		lastMessage.scrollIntoView();
	}

	// wait 1 second
	await new Promise(resolve => setTimeout(resolve, 1000));
	await loadChannelMessages(channelID);
}

function fixImageLinks(markdown) {
	// FIX IMAGE LINKS
	let images = markdown.querySelectorAll("img")
	let versionID = document.getElementById("versionID").innerText;
	for (let i = 0; i < images.length; i++) {
		let src = images[i].getAttribute("src")

		if (src.includes("/Assets/") || src.includes("/assets/")) {
			// get filename and strip away /Assets/
			let name = src.slice(8, src.length);
			let newSrc = "/media/"+versionID+"/name/"+name;
			images[i].setAttribute("src", newSrc);
		}
	}
}

function resetURL() {
	// get current course URL
	let courseURL = document.getElementById("courseURL").innerText;

	try {
		// change location of window
		window.history.replaceState("", "", courseURL + "/" + Alpine.store("sections").current)
	} catch {
		// change location of window
		window.history.replaceState("", "", courseURL)
	}
}