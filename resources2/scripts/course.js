function course_closeMenu() {
    Alpine.store('course').menuOpen = false;
}

function course_openMenu() {
    Alpine.store('course').menuOpen = true;
}

function course_viewSection(id) {
    Alpine.store("course").view = "section";
    Alpine.store("course").active = "section"+id;
    Alpine.store("course").sectionID = id;
}

function course_loadPreviousSection() {
    let sectionID = Alpine.store("course").sectionID;
    let sections = Alpine.store("course").sections;

    for (let i = 0; i < sections.length; i++) {
        if (sections[i].ID === sectionID) {
            if (i <= sections.length) {
                course_loadSection(sections[i-1].ID, sections[i-1].Name);
            }
        }
    }
}

function course_loadNextSection() {
    let sectionID = Alpine.store("course").sectionID;
    let sections = Alpine.store("course").sections;

    for (let i = 0; i < sections.length; i++) {
        if (sections[i].ID === sectionID) {
            if (i <= sections.length) {
                course_loadSection(sections[i+1].ID, sections[i+1].Name);
            }
        }
    }
}

function course_isThereANextSection(sectionID) {
    let sections = Alpine.store("course").sections;

    for (let i = 0; i < sections.length; i++) {
        if (sections[i].ID === sectionID) {
            if (i <= sections.length) {
                return true;
            }
        }
    }

    return false;
}

function course_isThereAPreviousSection(sectionID) {
    let sections = Alpine.store("course").sections;

    for (let i = 0; i < sections.length; i++) {
        if (sections[i].ID === sectionID) {
            if (i > 0) {
                return true;
            }
        }
    }

    return false;
}

function course_loadSection(sectionID, sectionName) {
    let releaseID = Alpine.store("course").releaseID;
    Alpine.store("course").sectionID = sectionID;
    Alpine.store("course").sectionName = sectionName;

    course_viewSection(sectionID);

    // is there a next section?
    if (course_isThereANextSection(sectionID)) {
        Alpine.store("course").next = true;
    } else {
        Alpine.store("course").next = false;
    }

    // is there a previous section?
    if (course_isThereAPreviousSection(sectionID)) {
        Alpine.store("course").previous = true;
    } else {
        Alpine.store("course").previous = false;
    }

    let html = sessionStorage.getItem("sections/"+sectionID);

    if (html !== null)
    {
        Alpine.store("course").sectionHTML = html;

        let course_markdown = document.getElementById("course_markdown");
        course_markdown.scrollIntoView();
    } else {
        fetch2("/v2/sections/"+sectionID+"/html", "GET", function(json) {
            if (json.Error) {
                return;
            }
            
            console.log("/v2/sections/:sectionID/html");
            console.log(json);
            
            let html = course_fixGithubImageLnks(json.Payload, releaseID);
            sessionStorage.setItem("sections/"+sectionID, html);
            Alpine.store("course").sectionHTML = html;

            let course_markdown = document.getElementById("course_markdown");
            course_markdown.scrollIntoView();

            // when a section is loaded close the menu
            Alpine.store("course").menuOpen = false;
        });
    }
}

function course_fixGithubImageLnks(markdownHTML, releaseID) {
    let markdownHTMLElem = document.createElement("div");
    markdownHTMLElem.innerHTML = markdownHTML;

    let images = markdownHTMLElem.querySelectorAll("img");

    for (let i = 0; i < images.length; i++) {
        let src = images[i].getAttribute("src");
        if (src.includes("/Assets/")) {
            images[i].setAttribute("src", "/v2/releases/"+releaseID+"/assets/"+src.slice(8, src.length));
        }
    }

    return markdownHTMLElem.innerHTML;
}

function course_loadSections() {
    // clear existing sections
    sessionStorage.clear();

    let releaseID = Alpine.store("course").releaseID;

    fetch2("/v2/releases/"+releaseID+"/sections", "GET", function(json) {
        if (json.Error !== "") {
            return;
        }

        console.log("/v2/releases/:releaseID/sections");
        console.log(json);
        Alpine.store("course").sections = json.Payload;

        if (Alpine.store("course").sectionID !== 0) {
            course_loadSection(Alpine.store("course").sectionID, Alpine.store("course").sectionName);
        } else if (json.Payload.length > 0) {
            course_loadSection(json.Payload[0].ID, json.Payload[0].Name);
        }
    });
}

function course_viewChannel() {
    Alpine.store("course").view = "channel";
    Alpine.store("course").active = "channel";
}

function course_loadReleases() {
    let courseID = Alpine.store("course").courseID;

    fetch2("/v2/course/"+courseID+"/releases", "GET", function(json) {
        if (json.Error !== "") {
            return;
        }

        console.log("/v2/course/:courseID/releases");
        console.log(json);

        Alpine.store("course").releases = json.Payload;

        // reload the sections
        course_loadSections();
    });
}

function course_getSelectedRelease() {
    let releaseID = Alpine.store("course").releaseID;
    let releases = Alpine.store("course").releases;

    for (let i = 0; i < releases.length; i++) {
        if (releaseID === releases[i].ID) {
            return releases[i];
        }
    }

    return null;
}