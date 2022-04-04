let courseNavTop;
let courseMain;
let courseMenu;

// offset from the top of the screen
let ogTopOffset;

document.addEventListener("DOMContentLoaded", function(event) {
	courseNavTop = document.getElementById("courseNavTop");
	courseMain = document.getElementById("courseMain");
	courseMenu = document.getElementById("courseMenu");
	ogTopOffset = courseNavTop.offsetTop;

	window.onscroll = function(){
		menuFollowScroll();
	}

	// sectionID when page first loads
	let sectionID = document.getElementById("sectionID").innerText;
	loadSection(sectionID);
});

document.addEventListener("alpine:init", function(event) {
	Alpine.store('courseView', {
		menuAvailable: true,
		menuOpen: false,
		allowNewPost: true,
		viewingPost: false,
		editingPost: false,
		editingContent: false,

		// views
		views: ['contents', 'chat', 'posts', 'menu'],
		view: 'contents',
	});
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
}

// show contents/lecture/section view
function viewContents() {
	Alpine.store("courseView").view = 'contents';
	Alpine.store("courseView").menuAvailable = true;
	Alpine.store("courseView").allowNewPost = true;
	Alpine.store("courseView").editingContent = false;
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

	fetch("/api/version/"+ versionID +"/posts", {
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
		let posts = document.getElementById("coursePostsMount");
		posts.innerHTML = "";

		console.log("json is:", json);

		for (let i = 0; i < json.length; i++) {
			let post = document.createElement("div");
			post.classList.add("course-card");
			post.setAttribute("postID", json[i].ID)
			posts.appendChild(post);

			let h4 = document.createElement("h3");
			h4.innerText = json[i].Markdown.slice(0, 30) + "...";
			post.appendChild(h4);

			post.addEventListener("click", function(event) {
				console.log("clicked on post with ID of:", this.getAttribute("postID"));
				loadPost(this.getAttribute("postID"));
			});
		}
	})
	.catch(function(err) {
		console.error(err);
	});
}

function viewEditPost() {
	let editPostTextInput = document.getElementById("editPostTextInput");
	let editPostID = document.getElementById("editPostID");

	fetch("/api/posts/" + editPostID.innerText + "/plaintext", {
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
	let editPostID = document.getElementById("editPostID");
	let editPostTextInput = document.getElementById("editPostTextInput");

	let formData = new FormData();
	formData.append("markdown", editPostTextInput.value);

	fetch("/api/posts/" + editPostID.innerText + "/update", {
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
		loadPost(editPostID.innerText);

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

		// if user decides to edit the post
		let editPostID = document.getElementById("editPostID");
		editPostID.innerText = json.ID;

		Alpine.store("courseView").editingPost = false;
		Alpine.store("courseView").viewingPost = true;

		postMount.innerHTML = json.Markdown;
		postMountUser.innerText = "@" + json.User.Username;

		// convert any elems with a link to open in new page
		convertHrefs();
	})
	.catch(function(err) {
		console.error(err);
	});
}

// async load section contents
function loadSection(sectionID) {
	console.log("loading section...");

	fetch("/api/section/"+sectionID, {
		method: "GET",
	})
	.then(function(resp) {
		console.log("got response");

		if (!resp.ok) {
			throw new Error("Response for loadSection was not ok");
		}

		console.log("response converted to json");
		return resp.json();
	})
	.then(function(sectionJson) {
		let content = document.getElementById("courseContent");
		content.innerHTML = "";

		let sectionTitle = document.getElementById("sectionTitle");
		if (sectionTitle === null) {
			throw new Error("sectionTitle was null!");
		}

		sectionTitle.innerText = sectionJson.Name;


		// TODO contents may be in english or spanish as well
		/*
			for (let i = 0; i < sectionJson.Contents.length; i++) {
				sectionJson.Contents[i].Language
				sectionJson.Contents[i].Markdown
			}
		*/
		if (sectionJson.Contents[0] === undefined) {
			content.innerHTML = `<p>This section is empty!</p>`;
			return
		}

		let courseContentsLanguage = document.getElementById("courseContentsLanguage")
		if (courseContentsLanguage === null) {
			throw new Error("courseContentsLanguage was null!");
		}

		courseContentsLanguage.innerText = sectionJson.Contents[0].Language;

		let markdown = document.createElement("div");
		markdown.innerHTML = sectionJson.Contents[0].Markdown;

		// FIX IMAGE LINKS
		let images = markdown.querySelectorAll("img")
		console.log("links to fix are:", images);

		let versionID = document.getElementById("versionID").innerText;

		for (let i = 0; i < images.length; i++) {
			let src = images[i].getAttribute("src")

			if (src.includes("/Assets/") || src.includes("/assets/")) {
				// get filename and strip away /Assets/
				let name = src.slice(8, src.length);
				let newSrc = "/media/"+versionID+"/name/"+name;
				images[i].setAttribute("src", newSrc);

				console.log("changed src to:", newSrc);
			}
		}

		markdown.setAttribute("markdown", "");

		window.scroll(0, 0);

		content.appendChild(markdown);

		Alpine.store("courseView").editingContent = false;
	})
	.catch(function(err) {
		console.error(err);
	});

	// set the current sectionID
	Alpine.store("sections").current = sectionID;

	// close menu
	Alpine.store("courseView").menuOpen = false;

	// show contents view
	viewContents();

	// get current course URL
	let courseURL = document.getElementById("courseURL").innerText;

	// change location of window
	window.history.replaceState("", "", courseURL + "/" + sectionID)

	// TODO
	/*
		- set next section button
		- set previous section button
	*/
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

function loadEditSectionPlaintext() {
	let currentSectionID = Alpine.store("sections").current;

	fetch("/api/section/"+ currentSectionID+ "/plaintext",{
		method: "get",
	})
	.then(function(resp) {
		console.log("got response");

		if (!resp.ok) {
			throw new Error("Response for loadSection was not ok");
		}

		console.log("response converted to json");
		return resp.json();
	})
	.then(function(sectionJson) {
		let courseContentEdit = document.getElementById("courseContentEdit");
		let courseContentEditSaveButton = document.getElementById("courseContentEditSaveButton");

		if (sectionJson.Contents[0] === undefined) {
			courseContentEdit.value = `This section is empty!`;
			return
		}

		courseContentEditSaveButton.setAttribute("contentID",  sectionJson.Contents[0].ID);
		courseContentEdit.value = sectionJson.Contents[0].Markdown;
	})
	.catch(function(err) {
		SendMessage("Failed to load section plaintext.")
		console.error(err);
	});
}

function saveEditSectionContent() {
	let courseContentEdit = document.getElementById("courseContentEdit");
	let courseContentEditSaveButton = document.getElementById("courseContentEditSaveButton");
	let contentID = courseContentEditSaveButton.getAttribute("contentID")
	let sectionID = Alpine.store("sections").current;

	let versionID = document.getElementById("versionID").innerText;

	let formData = new FormData();
	formData.append("content", courseContentEdit.value);
	formData.append("versionID", versionID);

	fetch("/api/section/"+sectionID+"/content/"+contentID+"/edit", {
		method: "post",
		body: formData,
	})
	.then(function(resp) {
		console.log("got response");

		if (!resp.ok) {
			throw new Error("Response for loadSection was not ok");
		}
	})
	.then(function() {
		SendMessage("Successfully updated section!");

		loadSection(sectionID);
	})
	.catch(function(err) {
		SendMessage("Error updating section");

		console.error(err);
	});
}