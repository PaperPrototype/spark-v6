fetch("/api/github/user/repos", {
	method: "GET",
})
.then(function(resp) {
	if (!resp.ok) {
		throw new Error("Error getting /api/user/repos")
	}

	return resp.json()
})
.then(function(json) {
	console.log(`${json.length} repos found through sparker api`);

	console.log(json);

	let courseSettingsSelect = document.getElementById("courseSettingsGithubConnectionSelect");

	for (let i = 0; i < json.length; i++) {
		let option = document.createElement("option");
		option.setAttribute("value", json[i].id);
		option.innerText = json[i].name;

		courseSettingsSelect.appendChild(option);
	}
})
.catch(function(err) {
	console.error(err);
});
