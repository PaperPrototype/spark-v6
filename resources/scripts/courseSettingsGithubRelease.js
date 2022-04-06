function loadRepos() {
	fetch("/api/github/user/repos", {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			throw new Error("Error getting /api/user/repos")
		}
	
		return resp.json()
	})
	.then(function(jsonOptions) {
		console.log(`${jsonOptions.length} repos found through sparker api`);
	
		let courseSettingsSelects = document.querySelectorAll("[courseSettingsGithubConnectionSelect]");

		for (let i = 0; i < courseSettingsSelects.length; i++) {
			courseSettingsSelects[i].innerHTML = "";
		}
	
		for (let selectIndex = 0; selectIndex < courseSettingsSelects.length; selectIndex++) {
			for (let optionIndex = 0; optionIndex < jsonOptions.length; optionIndex++) {
				// create option element
				let option = document.createElement("option");
				option.setAttribute("value", jsonOptions[optionIndex].id);
				option.innerText = jsonOptions[optionIndex].name;


		
				// append option element to select
				courseSettingsSelects[selectIndex].appendChild(option);
			}
		}
	})
	.catch(function(err) {
		console.error(err);
	});
}

document.addEventListener("alpine:init", function(event) {
	Alpine.data('next', () => ({
		next: false,
		done: false,
		repoID: "",

		testNext(elem) {
			if (elem.value === "") {
				SendMessage("You must select a repository before continuing")
				return
			}

			this.repoID = elem.value;
			this.next = true;
		},

		loadDone(elem) {
			console.log("loadDone")

			// fetch and display branches for selected repo
			fetch("/api/github/repo/" + this.repoID + "/branches", {
				method: "GET",
			})
			.then(function(resp) {
				if (!resp.ok) {
					throw new Error("Error getting branches");
				}

				return resp.json()
			})
			.then(function(json) {
				console.log("json is:", json);

				elem.innerHTML = "";

				for (let i = 0; i < json.length; i++) {
					let option = document.createElement("option");
					option.value = json[i].name;
					option.innerHTML = json[i].name;
					
					elem.append(option);
				}

				this.done = true;
			})
			.catch(function(err) {
				console.error(err);
			});
		},
	}))
});