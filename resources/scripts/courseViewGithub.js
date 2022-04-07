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

	console.log("desired tree is:", desiredTree);

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

	console.log("blobs are:", blobs);

	// clear course contents
	content.innerHTML = "";

	if (blobs.length === 0) {
		// fill in "this section is empty"
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
		console.log("text is:", json);

		let courseContentsLanguage = document.getElementById("courseContentsLanguage")
		if (courseContentsLanguage === null) {
			throw new Error("courseContentsLanguage was null!");
		}

		courseContentsLanguage.innerText = json.Name;

		let markdown = document.createElement("div");
		markdown.innerHTML = json.Markdown;

		// FIX IMAGE LINKS
		let images = markdown.querySelectorAll("img")
		console.log("links to fix are:", images);

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

		content.append(markdown);

		window.scroll(0, 0);
	});

	// close menu
	Alpine.store("courseView").menuOpen = false;

	// show contents view
	viewContents();

	// get current course URL
	let courseURL = document.getElementById("courseURL").innerText;

	// change location of window
	// window.history.replaceState("", "", courseURL + "/sha/" + tree_sha)
}

async function loadGithubBlob(versionID, sha, path) {
	let response = await fetch("/api/github/version/"+versionID+"/content/"+sha+"/"+path);
	return response.json();
}