document.addEventListener("alpine:init", function(event) {
	Alpine.store("notify", {
		show: false,
		newestDate: "",
		new: false,
		firstLoad: true,
		count: 0,
	});
});

let abortNotifsController = new AbortController();

function reloadNotifs() {
	Alpine.store("notify").new = false;
	Alpine.store("notify").firstLoad = true;
	Alpine.store("notify").newestDate = "";
	Alpine.store("notify").count = 0;

	const noneHTML = `<p>no new notifications</p>`;

	const navNotifications = document.getElementById("navNotifications");
	navNotifications.innerHTML = noneHTML;

	abortNotifsController.abort();
	abortNotifsController = new AbortController();

	loadNotifications();
}

async function loadNotifications() {
	console.log("loading notifications...");

	const navNotifications = document.getElementById("navNotifications");

	const noneHTML = `<p>no new notifications</p>`;
	const errorHTML = `<p>error loading notifications</p>`;

	if (Alpine.store("notify").firstLoad) {
		navNotifications.innerHTML = noneHTML;
	}

	let newest = Alpine.store("notify").newestDate;

	let response;
	try {
		response = await fetch("/api/notifications/newest?newest="+newest, {
			method: "GET",
			signal: abortNotifsController.signal,
		});
	} catch (err) {
		if (err.name === "AbortError") {
			console.log("fetch comments was aborted...");
			return;
		}

		if (Alpine.store("notify").firstLoad) {
			navNotifications.innerHTML = errorHTML;
		}
		
		console.log("error getting notifications...:", err);

		// wait 2 seconds before trying again
		await new Promise(resolve => setTimeout(resolve, 2000));
		await loadNotifications();
		return
	}

	const json = await response.json();

	if (json.Notifs.length === 0 && Alpine.store("notify").firstLoad) {
		navNotifications.innerHTML = noneHTML;
		Alpine.store("notify").new = false;
	} else {
		if (Alpine.store("notify").firstLoad) {
			// clear notifications
			navNotifications.innerHTML = "";
		}

		for (let i = 0; i < json.Notifs.length; i++) {
			let notif = document.createElement("p");
			notif.setAttribute("notifURL", json.Notifs[i].URL);
			notif.setAttribute("external", "");
			notif.setAttribute("notifID", json.Notifs[i].ID);
			notif.setAttribute("class", "thin-light pad-05 bd hoverable");
			notif.setAttribute("style", "cursor:pointer; display:flex; flex-direction:row; flex-wrap:nowrap;");

			notif.innerHTML = 
			`<div>` +
				json.Notifs[i].Message +
			`</div>` +
			`<i style="margin-left:auto;" class="fa-solid fa-arrow-up-right-from-square"></i>`;

			// append new notification
			navNotifications.append(notif);
	
			notif.addEventListener("click", function(event) {
				console.log("clicked notification!");

				// set notification as read
				notifSetRead(this.getAttribute("notifID"));

				// go to notification url
				// open in new tab
				window.open(this.getAttribute("notifURL"), '_blank')

				reloadNotifs();
			});
	
			// only runs once if there is any new notifications
			if (i === json.Notifs.length - 1) {
				Alpine.store("notify").newestDate = json.Notifs[i].CreatedAt;
				Alpine.store("notify").new = true;

				Alpine.store("notify").count += json.Count;
			}
		}
	}

	// helper for checkiing if this is the first load and if there were any new notifications
	if (Alpine.store("notify").new) {
		Alpine.store("notify").firstLoad = false;
	} else {
		navNotifications.innerHTML = noneHTML;
	}

	// wait 1 second
	await new Promise(resolve => setTimeout(resolve, 1000));
	await loadNotifications();
}

async function notifSetRead(notifID) {

	let formData = new FormData();
	formData.append("notifID", notifID);

	try {
		await fetch("/api/notifications/done", {
			method: "POST",
			body: formData,
		});
	} catch(err) {
		console.error("error setting notification as done...");
		return
	}

	console.log("SUCCESSFULLY set notification to done");
}