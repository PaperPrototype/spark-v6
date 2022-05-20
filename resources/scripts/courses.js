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

document.addEventListener("DOMContentLoaded", function() {
	loadCourses().then(function() {
		console.log("finished getting courses");
		convertHrefs(document);
	});
});