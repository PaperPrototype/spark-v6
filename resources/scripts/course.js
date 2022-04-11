document.addEventListener("alpine:init", function(event) {
	Alpine.store("course", {
		tab: "course",
	})
});

document.addEventListener("DOMContentLoaded", function(event) {
	let versionID = document.getElementById("versionID").innerText;

	fetch("/api/github/version/" + versionID +"/tree", {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			throw new Error("Response was not ok!");
		}

		return resp.json();
	})
	.then(function(json) {
		let courseSections = document.getElementById("courseSections");

		for (let i = 0; i < json.tree.length; i++) {
			// if it is a folder
			if (
				json.tree[i].type === "tree" &&
				!json.tree[i].path.includes("Resources") && 
				!json.tree[i].path.includes("Assets") && 
				!json.tree[i].path.includes("Archive") && 
				!json.tree[i].path.includes("Ignore")
			) {
				let section = document.createElement("div");
				section.setAttribute("class", "pad-05");
				section.innerText = json.tree[i].path;
				courseSections.append(section);
			}
		}
	})
	.catch(function(err) {
		console.error(err);
	});
});