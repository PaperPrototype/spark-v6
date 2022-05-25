let hierarchySearchForm;

document.addEventListener("DOMContentLoaded", function(event) {
	hierarchySearchForm = document.getElementById("hierarchySearchForm");
	hierarchySearchForm.addEventListener("submit", loadSearchResults, false);
});

document.addEventListener("alpine:init", function(event) {

});

function loadSearchResults(event) {
	if (event) {
		event.preventDefault();
	}

	let hierarchySearchInput = document.getElementById("hierarchySearchInput");

	if (hierarchySearchInput.value.trim() === "") {
		SendMessage("Add a search input.");
		return;
	}

	fetch("/api/courses?search="+hierarchySearchInput.value, {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			SendMessage("Error fetching results.");
			throw resp;
		}

		return resp.json();
	})
	.then(function(json) {
		let hierarchySearchCourses = document.getElementById("hierarchySearchCourses");
		hierarchySearchCourses.innerHTML = "";

		let courseBaseURL = document.getElementById("courseBaseURL").innerText;

		console.log(courseBaseURL);

		for (let i = 0; i < json.length; i++) {
			let courseCard = document.createElement("div");
			courseCard.setAttribute("class", "pad-05 thin-light bd");
			courseCard.setAttribute("style", "display:flex; margin-top:1rem; flex-direction:row;");

			courseCard.innerHTML = 
			`<div class="pad-05" href="/` + json[i].User.Username + `/` + json[i].Name + `" external>` + 
				json[i].Title + ` <i class="fa-solid fa-up-right-from-square"></i>` + 
			`</div>` + 
			`<div class="pad-05">` + 
				`by` + 
				`<a class="pad-05" href="/` + json[i].User.Username + `">` + 
					`@` + json[i].User.Username + 
				`</a>` + 
			`</div>` + 
			`<form style="margin-left:auto;" action="` + courseBaseURL + `/settings/prerequisites/new" method="post">` + 
				`<input name="preqCourseID" value="` + json[i].ID + `" hidden>` + 
				`<button type="submit" class="pad-05 bd-none hoverable-bg-code bd">` + 
					"Add pre-requisite" + 
				`</button>` +
			`</form>`;

			hierarchySearchCourses.append(courseCard);
			convertHrefs(courseCard);
		}
	})
	.catch(function(err) {
		console.log("error fetching results. Response is:" + err);
	})
}