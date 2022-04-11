function convertHrefs() {
	let hrefs = document.querySelectorAll("[href]");

	for (let i = 0; i < hrefs.length; i++) {
		hrefs[i].addEventListener("click", function(event) {
			if (this.hasAttribute("external")) {
				// open link in new window
				window.open(this.getAttribute("href"), '_blank')
				console.log("clicked");
				return
			}
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
