document.addEventListener("DOMContentLoaded", function() {
	setTheme(); // inital page load

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
			themeText.innerText = "Gamer Mode. On.";
		} else {
			// set to auto
			theme.classList.remove("theme-gamer");
			theme.classList.add("theme");

			// set local storage
			localStorage.setItem("sparker/theme", "auto");

			console.log("theme:auto")
			themeText.innerText = "Auto";
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
			themeText.innerText = "Gamer Mode. On.";
		} else {
			// set to auto
			theme.classList.remove("theme-gamer");
			theme.classList.add("theme");

			// set local storage
			localStorage.setItem("sparker/theme", "auto");

			console.log("theme:auto")
			themeText.innerText = "Auto";
		}
	}
}

// set the theme to whatever localStorage setting says
function setTheme() {
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
			themeText.innerText = "Auto";
		} else {
			// set to gamer
			theme.classList.remove("theme");
			theme.classList.add("theme-gamer");

			// set local storage
			localStorage.setItem("sparker/theme", "gamer");

			console.log("theme:gamer")
			themeText.innerText = "Gamer Mode. On.";
		}
	}
}