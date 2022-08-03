
let courseNavTop;
let courseMain;
let courseMenu;
let usingGithub = false;

// offset from the top of the screen
let ogTopOffset;

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

// show posts view
function viewPosts() {
	resetComments();

	let postsTitle = document.getElementById("postsTitle");
	postsTitle.innerText = "Posts";

	Alpine.store("courseView").oldestPost = "";
	Alpine.store("courseView").view = 'posts';
	Alpine.store("courseView").menuAvailable = false;
	Alpine.store("courseView").menuOpen = false;
	Alpine.store("courseView").viewingPost = false;
	Alpine.store("courseView").editingContent = false;
	Alpine.store("courseView").viewingComments = false;
}

function loadPosts(versionID) {
	console.log("loading posts...");

	resetComments();

	let coursePostsMount = document.getElementById("coursePostsMount");
	coursePostsMount.innerHTML = "";

	fetch("/api/version/"+ versionID +"/posts", {
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
		let posts = document.createElement("div");
		posts.setAttribute("class", "post-cards");
		posts.setAttribute("style", "margin-bottom:1rem;");
		posts.innerHTML = "";

		// first post that is actually a button for creating a new post
		let newPost = document.createElement("div");
		newPost.classList.add("text-center");
		newPost.classList.add("post-card");
		newPost.setAttribute("x-data", "");
		newPost.setAttribute("x-on:click", "newPostView({ version_id:" + versionID + " })");

		// plus icon
		let plusIcon = document.createElement("div");
		plusIcon.style = "margin-top:50%;"
		plusIcon.innerHTML = `<i class="fa-solid fa-plus"></i>`;
		newPost.appendChild(plusIcon);
		
		// title
		let h3 = document.createElement("h3");
		h3.style = "margin-top:auto;";
		h3.innerText = `New Post`;
		newPost.appendChild(h3);

		posts.appendChild(newPost);

		for (let i = 0; i < json.length; i++) {
			let post = document.createElement("div");
			post.classList.add("post-card");
			post.setAttribute("postID", json[i].ID)
			posts.appendChild(post);

			let elements = document.createElement("div");
			elements.innerHTML = json[i].Markdown;

			let title = json[i].Title;

			let hasCoverMedia = false;
			let images = elements.querySelectorAll("img");
			if (images.length !== 0)
			{
				let image = images[0];
				image.style = "width:100%;";

				let imageContainer = document.createElement("div");
				imageContainer.style = " border-radius:0.45rem; height:6rem; overflow-y:hidden;"

				imageContainer.append(image);
				post.appendChild(imageContainer);

				hasCoverMedia = true;
			} 
			// no images, look for embed instead
			else {
				// lets check for embedded youtube video
				let iframes = elements.querySelectorAll("iframe");
				for (let i = 0; i < iframes.length; i++)
				{
					let iframe = iframes[i];
					if (iframe.getAttribute("src") !== "")
					{
						console.log("post has iframe.")
						console.log(iframe.getAttribute("src"));

						// edit styles
						iframe.style = "width:100%; height:6rem; border-radius:0.45rem; padding:0; margin:0; overflow:hidden;";
						let src = iframe.getAttribute("src");
						src += "?modestbranding=1";
						iframe.setAttribute("src", src);

						post.appendChild(iframe);

						hasCoverMedia = true;

						break; // exit loop
					}
				}
			}

			let h4 = document.createElement("h4");
			// align to bottom
			h4.style = "margin-top: auto; margin-bottom:0.5rem;"; // parent is flexbox
			h4.setAttribute("class", "pad-05");

			let maxTitleLength = 50;
			if (hasCoverMedia === false) // no cover media
			{
				maxTitleLength = 100;

			    h4.style = "margin-top: auto; margin-bottom:0.5rem; text-align:center; justify-content:center; padding:auto;"; // parent is flexbox
			}

			// trim title if it is too large
			if (title.length > maxTitleLength)
			{
				h4.innerText = title.slice(0, maxTitleLength - 3).trim() + "...";
			} else {
				h4.innerText = title;
			}
			post.appendChild(h4);

			post.addEventListener("click", function(event) {
				console.log("clicked on post with ID of:", this.getAttribute("postID"));
				loadPost(this.getAttribute("postID"));
			});
		}

		// all posts
		let postsWrapper = document.createElement("div");
		postsWrapper.classList.add("pad-sides-5");
		postsWrapper.setAttribute("style", "margin-bottom:1rem;");
		postsWrapper.append(posts);

		// "load more posts" button
		let footer = document.createElement("div");
		footer.innerText = "Load more posts";
		footer.setAttribute("class", "pad-05 mar-sides-5 bd hoverable thin-light text-center");

		coursePostsMount.append(postsWrapper); // set posts
		coursePostsMount.append(footer); // set footer
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
		courseMain.style.paddingTop = courseNavTop.getBoundingClientRect().height + "px";
	} else {
		// top

		// stop menu from sticking to the top as a nav bar
		courseNavTop.classList.remove("course-top-nav-fixed");

		// reset courseMain's margin
		courseMain.style.paddingTop = "0";

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

		// make username links work
		convertHrefs(channelMount);
	}

	// wait 1 second
	await new Promise(resolve => setTimeout(resolve, 1000));
	await loadChannelMessages(channelID);
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