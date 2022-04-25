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
		console.log("courses:", json);

		let cards = document.getElementById("cards");
		cards.innerHTML = "";

		if (json.length === 0) {
			cards.innerText = "No courses found.";
		}

		for (let i = 0; i < json.length; i++) {
			let card = document.createElement("div");
			card.setAttribute("class", "course-card-wrapper");

			let title = json[i].Title.slice(0, 100);
			if (title.length < json[i].Title.length)
			{
				title += "...";
			}

			let imageURL = json[i].Release.ImageURL;
			if (imageURL === "")
			{
				// set default
				imageURL = "/resources/images/planet.png";
			}

			card.innerHTML = 
			`<div class="course-card" href="/` + json[i].User.Username + "/" + json[i].Name + `">` +
				`<div class="course-card-img-wrapper">` +
					`<img class="course-card-img" style='background-image:url(` + imageURL + `);'>` + 
				`</div>` +
				`<div class="course-card-content">` + 
					`<h3 class="c-bold course-card-title">` + title + `</h3>` +
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