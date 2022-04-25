async function loadCourses() {
	await fetch("/api/courses", {
		method: "GET"
	})
	.then(function(resp) {
		if (!resp.ok) {
			sendMessage("Error getting courses.");
			throw new Error('HTPP error status = ' + resp.status);
		}

		return resp.json();
	})
	.then(function(json) {
		let cards = document.getElementById("cards");
		cards.innerHTML = "";

		if (json.length === 0) {
			cards.innerText = "No courses found.";
		}

		for (let i = 0; i < json.length; i++) {
			let card = document.createElement("div");
			card.setAttribute("class", "course-card-wrapper");

			card.innerHTML = 
			`<div class="course-card" href="/` + json[i].User.Username + "/" + json[i].Name + `">` +
				`<div class="course-card-img-wrapper">` +
					`<img class="course-card-img" style='background-image:url(/resources/images/homepage.png);'>` + 
				`</div>` +
				`<div class="course-card-content">` + 
					`<h2 class="c-bold course-card-title">` + json[i].Title + `</h2>` +
				`</div>` +
			`</div>` + 
			`<p>` + `by <a href="/` + json[i].User.Username + `">@` + json[i].User.Username + `</a></p>`;

			cards.appendChild(card);
		}
	})
	.catch(function(err) {
		console.log(err)
	});
}

document.addEventListener("DOMContentLoaded", function() {
	loadCourses().then(function() {
		console.log("finished getting courses");
		convertHrefs(document);
	});
});