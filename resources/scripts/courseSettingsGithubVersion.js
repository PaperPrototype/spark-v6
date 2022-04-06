
function loadRepoBranchCommits(repoID, branch) {
	fetch("/api/github/repo/" + repoID +"/branch/" + branch + "/commits", {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			throw new Error("Error getting commits for repo branch");
		}

		return resp.josn()
	})
	.then(function(json) {
		console.log("json is:", json);
	})
	.catch(function(err) {
		console.error(err)
	});
}