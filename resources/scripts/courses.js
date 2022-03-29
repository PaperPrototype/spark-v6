async function loadCourses() {
	await fetch("/api/courses", {
		method: "GET"
	})
	.then(function(resp) {
		console.log("resp:", resp)
		if (!resp.ok) {
			sendMessage("Error getting courses.");
			throw new Error('HTPP error status = ' + resp.status);
		}

		return resp.json();
	})
	.then(function(json) {
		console.log("json:", json)

		let cards = document.getElementById("cards");
		cards.innerHTML = "";

		if (json.length === 0) {
			cards.innerText = "No courses found.";
		}

		for (let i = 0; i < json.length; i++) {
			console.log(json[i]);

			let card = document.createElement("div");
			card.classList.add("course-card");
			card.setAttribute("href", "/"+json[i].User.Username + "/" + json[i].Name);

			card.innerHTML = 
			`<h2>` + json[i].Title + `</h2>` +
			`<p>` + `by @` + json[i].User.Username + `</p>`;

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
		convertHrefs();
	});
});