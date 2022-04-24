document.addEventListener("DOMContentLoaded", function() {
	let toggle = document.getElementById("themeToggle");
	toggle.addEventListener("click", function() {
		let theme = document.getElementById("theme");
		let themeText = document.getElementById("themeToggleText");

		let savedTheme = sessionStorage.getItem("sparker/theme");

		if (savedTheme === "")
		{
			if (theme.classList.contains("theme")) {
				// set to dark
				theme.classList.remove("theme");
				theme.classList.add("theme-dark");

				// set session storage
				sessionStorage.setItem("sparker/theme", "dark");

				console.log("theme:dark")
				themeText.innerText = "Dark";
			} else {
				// set to auto
				theme.classList.remove("theme-dark");
				theme.classList.add("theme");

				// set session storage
				sessionStorage.setItem("sparker/theme", "auto");

				console.log("theme:auto")
				themeText.innerText = "Auto";
			}
		} else {
			if (savedTheme === "auto") {
				// set to dark
				theme.classList.remove("theme");
				theme.classList.add("theme-dark");

				// set session storage
				sessionStorage.setItem("sparker/theme", "dark");

				console.log("theme:dark")
				themeText.innerText = "Dark";
			} else {
				// set to auto
				theme.classList.remove("theme-dark");
				theme.classList.add("theme");

				// set session storage
				sessionStorage.setItem("sparker/theme", "auto");

				console.log("theme:auto")
				themeText.innerText = "Auto";
			}
		}

	});
});