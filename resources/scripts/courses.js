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
			let card = createCourseCard(json[i]);

			cards.appendChild(card);
		}
	})
	.catch(function(err) {
		console.log(err)
	});
}

function createCourseCard(course) {
	let card = document.createElement("div");
	let title = course.Title.slice(0, 60);
	if (title.length < course.Title.length)
	{
		title += "...";
	}

	let imageURL = course.Release.ImageURL;
	if (imageURL === "")
	{
		// set default
		imageURL = "/resources/images/planet.png";
	}

	card.innerHTML = 
	`<div class="course-card-wrapper">` +
		`<div class="course-card" href="/` + course.User.Username + "/" + course.Name + `">` +
			`<div class="course-card-img-wrapper">` +
				`<img class="course-card-img" src="` + imageURL + `">` + 
			`</div>` +
			`<div class="course-card-content">` + 
				`<h3 class="c-bold course-card-title">` + title + `</h3>` +
				`<div class="course-card-subtitle">` + course.Subtitle + `</div>` + 
			`</div>` +
		`</div>` +
	`</div>` + 
	`<p class="course-card-footer">` + 
		`by <a href="/` + course.User.Username + `">@` + course.User.Username + `</a>` + 
	`</p>`;

	return card;
}

document.addEventListener("DOMContentLoaded", function() {
	loadCourses().then(function() {
		console.log("finished getting courses");
		convertHrefs(document);
	});
});