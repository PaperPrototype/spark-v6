function convertHrefs(element) {
	let hrefs = element.querySelectorAll("[href]");

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
	convertHrefs(document);
});

document.addEventListener("alpine:init", function(event) {
	Alpine.store("nav", {
		show: false,
	});
});
