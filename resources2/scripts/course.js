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

function course_loadSection(sectionID) {
    let releaseID = Alpine.store("course").releaseID;
    Alpine.store("course").sectionID = sectionID;

    let html = sessionStorage.getItem("sections/"+sectionID);

    if (html !== null)
    {
        Alpine.store("course").sectionHTML = html;
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
            course_loadSection(Alpine.store("course").sectionID);
        } else if (json.Payload.length > 0) {
            course_loadSection(json.Payload[0].ID);
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