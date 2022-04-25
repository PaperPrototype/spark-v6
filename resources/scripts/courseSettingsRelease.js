function loadReleaseImageURLs(select, versionID) {
	select.innerHTML = `<option value="">Loading...</option>`

	fetch("/api/github/version/" + versionID + "/tree", {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			SendMessage("Error getting Cover/Banner images");
			throw new Error("Response for loadReleaseImageURLs was not ok!");
		}

		return resp.json();
	})
	.then(function(json) {
		console.log("tree is:", json);

		let assets = [];

		for (let i = 0; i < json.tree.length; i++) {
			if (
				// case insensitive
				json.tree[i].path.toLowerCase().includes("assets/")
			) {
				assets.push(json.tree[i].path);
			}
		}

		select.innerHTML = "";

		for (let i = 0; i < assets.length; i++) {
			let option = document.createElement("option");
			option.innerText = assets[i];
			option.value = "/media/" + versionID + "/name/" + assets[i].slice(7, assets[i].length);

			select.append(option);
		}
	})
	.catch(function(err) {
		console.error(err);
	});
}