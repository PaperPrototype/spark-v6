function convertHrefs() {
	let hrefs = document.querySelectorAll("[href]");
	console.log(hrefs);

	for (let i = 0; i < hrefs.length; i++) {
		console.log("adding click + linkability to element");

		hrefs[i].addEventListener("click", function(event) {
			this.setAttribute("target", "_blank");
			window.location = this.getAttribute("href");
			console.log("clicked");
		});
	}
}

document.addEventListener("DOMContentLoaded", function(event) {
	convertHrefs();
});

document.addEventListener("alpine:init", function(event) {
	Alpine.store("sideNav", {
		show: false,
	});
});