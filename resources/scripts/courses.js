async function loadSearchCourses() {
	let search = document.getElementById("search").value;
    
	await fetch("/api/courses?search="+search, {
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

		let cards = document.getElementById("courses");
        cards.classList.add("course-cards");

		if (json.length === 0) {
			cards.innerText = "No courses found.";
		}

		for (let i = 0; i < json.length; i++) {
			let card = createCourseCard(json[i]);

			cards.appendChild(card);
		}
        
        convertHrefs(cards);
	})
	.catch(function(err) {
		console.log(err)
	});
}

async function loadMoreLevelCourses() {
    if (Alpine.store("courses").done) {
        return;
    }

    let level = Alpine.store("courses").level;

    // increase for next call
    Alpine.store("courses").level += 1;

	await fetch("/api/level/"+level, {
		method: "GET",
	})
	.then(function(resp) {
		if (!resp.ok) {
			throw new Error("error getting resp for loadCoursesAtLevel with level "+level);
		}

		return resp.json();
	})
	.then(function(json) {
		console.log("level " + level + ":", json);

        if (json.length === 0) {
            Alpine.store("courses").done = true;
            return;
		}

        let courses = document.getElementById("courses");

        let header = document.createElement("h3");
        header.setAttribute("style", "margin-bottom:0;");
        header.setAttribute("class", "pad-sides-5");
        header.innerText = "Level " + level;
        courses.append(header);

        let cards = document.createElement("div");
        cards.setAttribute("class", "course-cards");

		for (let i = 0; i < json.length; i++) {
			let card = createCourseCard(json[i]);

			cards.appendChild(card);
		}

        courses.append(cards);
	})
	.catch(function(err) {
		console.error(err);
	});
}

function resetCourses() {
    Alpine.store("courses").done = false;
    Alpine.store("courses").level = 0;
    let courses = document.getElementById("courses");
    courses.innerHTML = "";
}

window.onscroll = function(){
    // console.log("scroll: " + (window.scrollY + window.innerHeight));
    // console.log("height: " + document.body.clientHeight);

    // if not searching then load hierarchy
    if (!Alpine.store("courses").isSearching) {
        let scrollAmount = (window.scrollY + window.innerHeight);
        let height = document.body.clientHeight;
    
        console.log("load more content");
        if  (scrollAmount >= (height - 100)){
            loadMoreLevelCourses().then(function() {
                console.log("finished getting courses");
                convertHrefs(document);
            });
        }
    }
}

document.addEventListener("alpine:init", async function(event) { 
    Alpine.store("courses", {
        level: 0,
        done: false,
        isSearching: false,
    })

    resetCourses();

	let search = document.getElementById("search").value;
    
    if (search !== "") {
        Alpine.store("courses").isSearching = true;

        loadSearchCourses();
    } else {
        // not searching
        await loadMoreLevelCourses().then(function() {
            convertHrefs(document);
        });
        await loadMoreLevelCourses().then(function() {
            convertHrefs(document);
        });
    }
})