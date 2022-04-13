// async load section contents
function loadUploadSection(sectionID) {
	console.log("loading section...");

	let content = document.getElementById("courseContent");
	content.innerHTML = "<p>Loading...</p>";

	fetch("/api/section/"+sectionID, {
		method: "GET",
	})
	.then(function(resp) {

		if (!resp.ok) {
			throw new Error("Response for loadSection was not ok");
		}

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

	// get current course URL
	let courseURL = document.getElementById("courseURL").innerText;

	// change location of window
	window.history.replaceState("", "", courseURL + "/" + sectionID)
}


function loadEditSectionPlaintext() {
	let currentSectionID = Alpine.store("sections").current;

	let courseID = document.getElementById("courseID").innerText;

	if (Alpine.store("courseView").usingGithub) {
		SendMessage("Editing github based courses is not available yet.")
		return
	}

	fetch("/api/section/"+ currentSectionID+ "/plaintext?course_id="+courseID,{
		method: "GET",
	})
	.then(function(resp) {

		if (!resp.ok) {
			SendMessage("You must be the course author to edit this.")
			throw new Error("Response for loadSection was not ok");
		}

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

		Alpine.store("courseView").editingContent = true; 
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