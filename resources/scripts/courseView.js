let courseNavTop;
let courseMain;
let courseMenu;
let usingGithub = false;

// offset from the top of the screen
let ogTopOffset;

// form for sending comments on a post
let postCommentsForm;

document.addEventListener("DOMContentLoaded", function(event) {
	courseNavTop = document.getElementById("courseNavTop");
	courseMain = document.getElementById("courseMain");
	courseMenu = document.getElementById("courseMenu");
	ogTopOffset = courseNavTop.offsetTop;

	window.onscroll = function(){
		menuFollowScroll();
	}

	postCommentsForm = document.getElementById("postCommentsForm");
	postCommentsForm.addEventListener("submit", sendComment, false)
});

// show chat view
function viewChat() {
	Alpine.store("courseView").view = 'chat';
	Alpine.store("courseView").menuAvailable = false;
	Alpine.store("courseView").menuOpen = false;
	Alpine.store("courseView").editingContent = false;
}

// show posts view
function viewPosts() {
	Alpine.store("courseView").view = 'posts';
	Alpine.store("courseView").menuAvailable = false;
	Alpine.store("courseView").menuOpen = false;
	Alpine.store("courseView").viewingPost = false;
	Alpine.store("courseView").editingContent = false;
	Alpine.store("courseView").viewingComments = false;
}

// show contents/lecture/section view
function viewContents() {
	Alpine.store("courseView").view = 'contents';
	Alpine.store("courseView").menuAvailable = true;
	Alpine.store("courseView").allowNewPost = true;
	Alpine.store("courseView").editingContent = false;
	Alpine.store("courseView").viewingComments = false;
}

// toggle nav menu
function toggleMenu() {
	if (Alpine.store("courseView").menuAvailable) {
		Alpine.store("courseView").menuOpen = ! Alpine.store("courseView").menuOpen;
	} else {
		Alpine.store("courseView").menuOpen = false;
	}
}

function loadPosts(versionID) {
	console.log("loading posts...");

	let coursePostsMount = document.getElementById("coursePostsMount");
	coursePostsMount.innerHTML = "";

	fetch("/api/version/"+ versionID +"/posts/portfolio", {
		method: "GET",
	})
	.then(function(resp) {
		console.log("got response");

		if (!resp.ok) {
			SendMessage("Error loading post.");
			throw new Error("Response for loadPosts was not ok");
		}

		console.log("response converted to json");
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

		console.log("json is:", json);

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
			console.log("TODO create see all view:",);
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
		console.log("got response");

		if (!resp.ok) {
			SendMessage("Error loading post.");
			throw new Error("Response for loadPosts was not ok");
		}

		console.log("response converted to json");
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
		console.log("got response");

		if (!resp.ok) {
			throw new Error("Response for loadPost was not ok");
		}

		console.log("response converted to json");
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

		// convert any elems with a link to open in new page
		convertHrefs();
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
		console.log(courseMain.style.marginTop);
	} else {
		// top

		// stop menu from sticking to the top nav bar
		courseNavTop.classList.remove("course-top-nav-fixed");
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

	fetch("/api/posts/" + postID + "/comment", {
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

function loadComments() {
	// long polling
	// https://javascript.info/long-polling
	async function subscribe() {
		let response = await fetch("/subscribe");
	
		if (response.status == 502) {
		// Status 502 is a connection timeout error,
		// may happen when the connection was pending for too long,
		// and the remote server or a proxy closed it
		// let's reconnect
		await subscribe();
		} else if (response.status != 200) {
		// An error - let's show it
		showMessage(response.statusText);
		// Reconnect in one second
		await new Promise(resolve => setTimeout(resolve, 1000));
		await subscribe();
		} else {
		// Get and show the message
		let message = await response.text();
		showMessage(message);
		// Call subscribe() again to get the next message
		await subscribe();
		}
	}
	
	subscribe();
}