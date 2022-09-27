document.addEventListener("alpine:init", function(event) {
    Alpine.store("popups_github_release", {
        repoID: "",
        branch: "",
        SHA: "",

        // alpine template arrays
        repos: [],
        branches: [],
        commits: [],
    });
});

let startRepoID = 0;
let startBranch = "";
let startSHA = "";

/* POPUPS */

// load githubRelease popup
function popups_popupGithubRelease() {
    popups_LoadingOn();

    Alpine.store("popups").popup = "github_release";

    var releaseID = document.getElementById("releaseID").innerText;

    fetch2("/v2/releases/"+releaseID+"/github", "GET", function(json) {
        // load repositories
        // (can't load anything else since we don't know what repo the user has selected)
        popups_GR_loadRepos();
        popups_LoadingOn(); // force turn on loading again

        if (json.Error !== "") {
            // stays default since there is no data for prefilling
            return;
        }

        // load default values
        // default repo
        Alpine.store("popups_github_release").repoID = json.Payload.RepoID;
        
        // we know which repo the user selected
        // so we can proceed to load branches
        popups_GR_loadBranches();
        popups_LoadingOn(); // force turn on loading again

        // default branch
        Alpine.store("popups_github_release").branch = json.Payload.Branch;

        // we know which branch the user has selected
        // so we can proceed to load commits
        popups_GR_loadCommits();

        // default commit
        Alpine.store("popups_github_release").SHA = json.Payload.SHA;

        startRepoID = json.Payload.RepoID;
        startBranch = json.Payload.Branch;
        startSHA = json.Payload.SHA;
    });
}

// load githubRelease sections popup
function popups_popupGithubSections() {
    Alpine.store("popups").popup = "github_sections";
}

/* GITHUB RELEASE POPUP 
    "GR" stands for "Github Release" */

function popups_GR_loadRepos() {
    fetch2("/v2/user/github/repos", "GET", function(json) {
        if (json.Error !== "") {
            return;
        }

        Alpine.store("popups_github_release").repos = json.Payload;

        console.log("/v2/user/github/repos");
        // console.log(json);
    });
}

function popups_GR_loadBranches() {
    var repoID = Alpine.store("popups_github_release").repoID;

    fetch2("/v2/user/github/repo/"+repoID+"/branches", "GET", function(json) {
        if (json.Error !== "") {
            SendMessage("It's possible you don't have access to this repositories commits");
            return;
        }

        Alpine.store("popups_github_release").branches = json.Payload;

        console.log("/v2/user/github/repo/:repoID/branches");
        // console.log(json);
    });
}

function popups_GR_loadCommits() {
    var repoID = Alpine.store("popups_github_release").repoID;
    var branch = Alpine.store("popups_github_release").branch;

    fetch2("/v2/user/github/repo/"+repoID+"/branch/"+branch+"/commits", "GET", function(json) {
        if (json.Error !== "") {
            return;
        }

        Alpine.store("popups_github_release").commits = json.Payload;

        console.log("/v2/user/github/repo/:repoID/branch/:branch/commits");
        // console.log(json);
    });
}

/* utility multi-use functions */
function popups_GR_selectedNewRepository() {
    console.log("selected new repository");

    console.log(Alpine.store("popups_github_release"));

    Alpine.store("popups_github_release").branches = [];
    Alpine.store("popups_github_release").branch = "";
    Alpine.store("popups_github_release").commits = [];
    Alpine.store("popups_github_release").SHA = "";

    popups_GR_loadBranches();
}

function popups_GR_selectedNewBranch() {
    console.log("selected new branch");

    console.log(Alpine.store("popups_github_release"));

    Alpine.store("popups_github_release").commits = [];

    popups_GR_loadCommits();
}

/* SAVE AND SELECT SECTIONS */
function popups_GR_save() {
    popups_LoadingOn();

    let releaseID = document.getElementById("releaseID").innerText;

    let repoID = Alpine.store("popups_github_release").repoID;
    let branch =  Alpine.store("popups_github_release").branch;
    let SHA = Alpine.store("popups_github_release").SHA;

    if (startRepoID === repoID && startBranch === branch && startSHA === SHA) {
        // move to next without saving since nothing changed
        Alpine.store("popups").popup = "";
        popups_LoadingOff();
        return;
    }

    popups_LoadingOff();

    let formData = new FormData();
    formData.append("releaseID", releaseID)
    formData.append("repoID",  repoID);
    formData.append("branch", branch);
    formData.append("SHA",  SHA);

    fetch("/v2/releases/"+releaseID+"/github", {
        method: "POST",
        body: formData,
    })
    .then(function(resp) {
        if (!resp.ok) {

            var json = resp.json();
            if (json.Error !== "" && json.Error !== undefined)
            {
                throw json.Error;
            }

            throw "Error getting response";
        }

        return resp.json();
    })
    .then(function(json) {
        if (json.Error !== "") {
            SendMessage(json.Error);
        }

        Alpine.store("popups").popup = "";

        console.log("json:", json);
        popups_LoadingOff();
    })
    .catch(function(err) {
        SendMessage(err);
        popups_LoadingOff();
        console.error(err);
    });
}