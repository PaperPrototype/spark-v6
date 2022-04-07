function loadRepoBranchCommits(elem, repoID, branch) {
	fetch("/api/github/repo/" + repoID +"/branch/" + branch + "/commits", {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			throw new Error("Error getting commits for repo branch");
		}

		return resp.json()
	})
	.then(function(commitsJson) {
		console.log("json is:", commitsJson);

		elem.innerHTML = "";

		for (let i = 0; i < commitsJson.length; i++) {
			let option = document.createElement("option");
			option.value = commitsJson[i].sha;
			option.innerText = commitsJson[i].commit.message;

			elem.append(option);
		}
	})
	.catch(function(err) {
		console.error(err)
	});
}