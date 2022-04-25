document.addEventListener("alpine:init", function(event) {
	Alpine.store("course", {
		tab: "course", // defaults to the course tab
		viewingReview: false,
		editingReview: false,
		viewingComments: false,
		allowNewReview: true,
		offset: 0,
	});
});

document.addEventListener("DOMContentLoaded", function(event) {
	let versionID = document.getElementById("versionID").innerText;

	console.log("loading sections...");
	fetch("/api/github/version/" + versionID +"/tree", {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			throw new Error("Response was not ok!");
		}

		return resp.json();
	})
	.then(function(json) {
		let courseSections = document.getElementById("courseSections");

		for (let i = 0; i < json.tree.length; i++) {
			// if it is a folder
			if (
				json.tree[i].type === "tree" &&
				!json.tree[i].path.includes("Resources") && 
				!json.tree[i].path.includes("Assets") && 
				!json.tree[i].path.includes("Archive") && 
				!json.tree[i].path.includes("Ignore")
			) {
				let section = document.createElement("div");
				section.setAttribute("class", "pad-05");
				section.innerText = json.tree[i].path;
				courseSections.append(section);
			}
		}
	})
	.catch(function(err) {
		console.error(err);
	});
});

function loadReviews() {
	console.log("loading reviews...");

	let versionID = document.getElementById("versionID").innerText;
	let offset = 0;

	let reviewsMount = document.getElementById("reviewsMount");
	reviewsMount.innerHTML = "";

	let reviewsCount = document.getElementById("reviewsCount");

	// set offset
	Alpine.store("course").offset = 5;

	// only loads up to 20 review/posts at once
	fetch(`/api/version/` + versionID + `/reviews?offset=` + offset + `&limit=5`, {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			SendMessage("Error loading reviews");
			throw new Error("Response for course reviews was not ok!");
		}

		return resp.json();
	})
	.then(function(json) {
		reviewsCount.innerText = json.Count;

		for (let i = 0; i < json.Reviews.length; i++)
		{
			let review = document.createElement("div");
			review.setAttribute("class", "pad-05 bd thin-light");
			review.setAttribute("style", "margin-bottom:1rem;");

			let topbar = document.createElement("div");
			topbar.innerHTML = `by <a href="/` + json.Reviews[i].User.Username + `">@` + json.Reviews[i].User.Username + `</a> <div style="margin-left:auto;">` + json.Reviews[i].Rating + ` stars</div>`;
			topbar.setAttribute("style", "padding-bottom:0.5rem; display:flex;");
			review.append(topbar);

			let markdown = document.createElement("div");
			review.append(markdown);

			markdown.innerHTML = json.Reviews[i].Post.Markdown.slice(0, 500) + "...";
			reviewsMount.append(review);
		}
	})
	.catch(function(err) {
		console.error(err);
	});
}

function loadMoreReviews() {
	console.log("loading more reviews...");

	let versionID = document.getElementById("versionID").innerText;
	let offset = Alpine.store("course").offset;

	let reviewsMount = document.getElementById("reviewsMount");

	// only loads up to 20 review/posts at once
	fetch(`/api/version/` + versionID + `/reviews?offset=` + offset + `&limit=5`, {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			SendMessage("Error loading reviews");
			throw new Error("Response for course reviews was not ok!");
		}

		return resp.json();
	})
	.then(function(json) {
		for (let i = 0; i < json.Reviews.length; i++)
		{
			let review = document.createElement("div");
			review.setAttribute("class", "pad-05 bd thin-light");
			review.setAttribute("style", "margin-bottom:1rem;");

			let topbar = document.createElement("div");
			topbar.innerHTML = `by <a href="/` + json.Reviews[i].User.Username + `">@` + json.Reviews[i].User.Username + `</a>`;
			topbar.setAttribute("style", "padding-bottom:0.5rem;");
			review.append(topbar);

			let markdown = document.createElement("div");
			review.append(markdown);

			markdown.innerHTML = json.Reviews[i].Post.Markdown.slice(0, 500) + "...";
			reviewsMount.append(review);
		}

		Alpine.store("course").offset += 5;
	})
	.catch(function(err) {
		console.error(err);
	});
}

function postReview() {
	console.log("posting new review...");
	let postNewReviewMarkdown = document.getElementById("postNewReviewMarkdown");
	let postNewReviewRating = document.getElementById("postNewReviewRating");

	if (postNewReviewMarkdown.value === "" || postNewReviewRating === "")
	{
		SendMessage("You can't post an empty review!");
		return;
	}

	let versionID = document.getElementById("versionID").innerText;
	
	let formData = new FormData();
	formData.append("markdown", postNewReviewMarkdown.value);
	formData.append("rating", postNewReviewRating.value)

	fetch("/api/version/" + versionID + "/reviews/new", {
		method: "POST",
		body: formData,
	})
	.then(function(resp) {
		if (!resp.ok) {
			SendMessage("Error posting review (you can only post 1 review per course).");
			throw new Error("Response for postReview was not ok!");
		}

		return resp.text();
	})
	.then(function(text) {
		SendMessage("Successfully posted review");
		loadReviews(); // reload inital reviews to show the new review
		Alpine.store("course").allowNewReview = false;
	})
	.catch(function(err) {
		console.error(err);
	});
}