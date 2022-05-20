function loadGithubSection(tree_sha) {
	console.log("loading github section...");

	let content = document.getElementById("courseContent");
	content.innerHTML = "<p>Loading...</p>";

	// set the current id to the tree's sha
	Alpine.store("sections").current = tree_sha;

	Alpine.store("courseView").editingContent = false;

	// get github tree
	let fullTree = Alpine.store("githubTree").json;
	let desiredTree = null;

	// find tree with sha in github.tree
	for (let i = 0; i < fullTree.tree.length; i++) {
		if (fullTree.tree[i].sha === tree_sha) {
			desiredTree = fullTree.tree[i];
		}
	}

	if (desiredTree === null) {
		console.log("tree for that section was null");
		return
	}

	let blobs = [];
	// use desiredTree path to check for other sub paths that lead to blobs that end in .md
	for (let i = 0; i < fullTree.tree.length; i++) {
		if (
			fullTree.tree[i].path.includes(desiredTree.path + "/english.md") && 
			fullTree.tree[i].type === "blob" &&
			fullTree.tree[i].path.includes(".md")) {
			blobs.push(fullTree.tree[i]);
		}
	}

	// clear course contents
	content.innerHTML = "";

	let sectionTitle = document.getElementById("sectionTitle");
	if (sectionTitle === null) {
		console.error("sectionTitle element was null!");
	}

	// if there is no blobs
	if (blobs.length === 0) {
		// fill in "this section is empty"

		sectionTitle.innerText = desiredTree.path;
		
		content.innerHTML = `<p>This section is empty!</p>`;

		console.log("no contents found for that section");

		// show contents view
		viewContents();

		// close menu
		Alpine.store("courseView").menuOpen = false;
		return
	}

	let versionID = document.getElementById("versionID").innerText;

	// load first blob as section
	let resp = loadGithubBlob(versionID, fullTree.sha, blobs[0].path)
	resp.then(function(json) {
		sectionTitle.innerText = desiredTree.path;

		let courseContentsLanguage = document.getElementById("courseContentsLanguage")
		if (courseContentsLanguage === null) {
			console.error("courseContentsLanguage element was null!");
		}

		courseContentsLanguage.innerText = json.Name;

		let markdown = document.createElement("div");
		markdown.innerHTML = json.Markdown;

		// FIX IMAGE LINKS
		let images = markdown.querySelectorAll("img")
		for (let i = 0; i < images.length; i++) {
			let src = images[i].getAttribute("src")

			if (src.includes("/Assets/") || src.includes("/assets/")) {
				// get filename and strip away /Assets/
				let name = src.slice(8, src.length);
				let newSrc = "/media/"+versionID+"/name/"+name;
				images[i].setAttribute("src", newSrc);
			}
		}

		// SET HEADER LINKS
		let headersWithID = markdown.querySelectorAll("[id]")
		for (let i = 0; i < headersWithID.length; i++) {
			headersWithID[i].style = "cursor:pointer;"

			headersWithID[i].addEventListener("click", function(event) {
				let id = this.getAttribute("id")

				// if there is an anchor tag then go to the anchor
				window.location = courseURL + "/" + Alpine.store("sections").current + "#" + id;
				
				// anchor is hidden by default so scroll up a bit.
				window.scrollBy(0, -80);

				console.log("scrolling down so user can see anchor");
			});
		}

		markdown.setAttribute("markdown", "");

		convertHrefs(markdown);

		content.append(markdown);

		// get the anchor tag
		let currentUrl = document.URL,
		urlParts = currentUrl.split('#');
		let headerID = (urlParts.length > 1) ? urlParts[1] : null;

		// get current course URL
		let courseURL = document.getElementById("courseURL").innerText;

		// if there is not an anchor tag
		if (headerID === null || headerID === 'undefined') {
			console.log("headerID was undefined or null");

			window.scroll(0, 0);

			// change location of window
			window.history.replaceState("", "", courseURL + "/" + tree_sha)
		} else {
			// if there is an anchor tag then go to the anchor
			window.location = courseURL + "/" + tree_sha + "#" + headerID;
			
			// anchor is hidden by default so scroll up a bit.
			window.scrollBy(0, -80);
		}
	});
}

async function loadGithubBlob(versionID, sha, path) {
	let response = await fetch("/api/github/version/"+versionID+"/content/"+sha+"/"+path);
	return response.json();
}