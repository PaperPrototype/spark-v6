document.addEventListener("DOMContentLoaded", function(event) {
    loadTakenCourses();
    loadAuthoredCourses();
});

async function loadTakenCourses() {
    let profileUsername = document.getElementById("profileUsername").innerText;
    let takenCourses = document.getElementById("takenCourses");

    fetch("/api/user/"+profileUsername+"/courses", {
        method: "GET",
    })
    .then(function(resp) {
        if (!resp.ok) {
            throw "error getting taken courses for user";
        }

        return resp.json();
    })
    .then(function(json) {
        for (let i = 0; i < json.length; i++) {
            takenCourses.append(createCourseCard(json[i]));
        }

        convertHrefs(takenCourses);
    })
    .catch(function(err) {
        console.error(err);
    })
}

async function loadAuthoredCourses() {
    let profileUsername = document.getElementById("profileUsername").innerText;
    let authoredCourses = document.getElementById("authoredCourses");

    fetch("/api/user/"+profileUsername+"/authored", {
        method: "GET",
    })
    .then(function(resp) {
        if (!resp.ok) {
            throw "error getting authored courses for user";
        }

        return resp.json();
    })
    .then(function(json) {
        for (let i = 0; i < json.length; i++) {
            authoredCourses.append(createCourseCard(json[i]));
        }

        convertHrefs(authoredCourses);
    })
    .catch(function(err) {
        console.error(err);
    });
}