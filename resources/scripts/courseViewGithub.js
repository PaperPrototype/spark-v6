function loadGithubSection(tree_sha) {
	console.log("loading github section...");

	let versionID = document.getElementById("versionID").innerText;

	fetch("/api/github/version/" + versionID + "/trees/" + tree_sha, {
		method: "GET",
	})
	.then(function(resp) {
		console.log("got response");

		if (!resp.ok) {
			throw new Error("Response for loadGithubSection was not ok");
		}

		console.log("response converted to json");
		return resp.json();
	})
	.then(function(treeJson) {
		// clear the course lecture contents
		let content = document.getElementById("courseContent");
		content.innerHTML = "";

		// set the title
		let sectionTitle = treeJson.path;
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