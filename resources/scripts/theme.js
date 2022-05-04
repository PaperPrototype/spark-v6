document.addEventListener("DOMContentLoaded", function() {
	setToSavedTheme(); // inital page load

	let toggle = document.getElementById("themeToggle");
	toggle.addEventListener("click", function() {
		toggleTheme();
	});
});

function toggleTheme() {
	let theme = document.getElementById("theme");
	let themeText = document.getElementById("themeToggleText");

	let savedTheme = localStorage.getItem("sparker/theme");

	if (savedTheme === "")
	{
		console.log("saved theme localStorage was empty");
		if (theme.classList.contains("theme")) {
			// set to gamer
			theme.classList.remove("theme");
			theme.classList.add("theme-gamer");

			// set local storage
			localStorage.setItem("sparker/theme", "gamer");

			console.log("theme:gamer")
			themeText.innerText = "Theme: Gamer Mode";
		} else {
			// set to auto
			theme.classList.remove("theme-gamer");
			theme.classList.add("theme");

			// set local storage
			localStorage.setItem("sparker/theme", "auto");

			console.log("theme:auto")
			themeText.innerText = "Theme: Auto";
		}
	} else {
		console.log("saved theme was not empty");
		if (savedTheme === "auto") {
			// set to gamer
			theme.classList.remove("theme");
			theme.classList.add("theme-gamer");

			// set local storage
			localStorage.setItem("sparker/theme", "gamer");

			console.log("theme:gamer")
			themeText.innerText = "Theme: Gamer Mode";
		} else {
			// set to auto
			theme.classList.remove("theme-gamer");
			theme.classList.add("theme");

			// set local storage
			localStorage.setItem("sparker/theme", "auto");

			console.log("theme:auto")
			themeText.innerText = "Theme: Auto";
		}
	}
}

// set the theme to whatever localStorage setting says
function setToSavedTheme() {
	let theme = document.getElementById("theme");
	let themeText = document.getElementById("themeToggleText");

	let savedTheme = localStorage.getItem("sparker/theme");

	// default leave as is
	// otherwise 
	if (savedTheme !== "") {
		if (savedTheme === "auto") {
			// set to auto
			theme.classList.remove("theme-gamer");
			theme.classList.add("theme");

			// set local storage
			localStorage.setItem("sparker/theme", "auto");

			console.log("theme:auto")
			themeText.innerText = "Theme: Auto";
		} else {
			// set to gamer
			theme.classList.remove("theme");
			theme.classList.add("theme-gamer");

			// set local storage
			localStorage.setItem("sparker/theme", "gamer");

			console.log("theme:gamer")
			themeText.innerText = "Theme: Gamer Mode";
		}
	}
}