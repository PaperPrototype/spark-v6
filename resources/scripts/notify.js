document.addEventListener("alpine:init", function(event) {
	Alpine.store("notify", {
		show: false,
		newestDate: "",
		new: false,
		firstLoad: true,
		count: 0,
	});
});

async function loadNotifications() {
	console.log("loading notifications...");

	const navNotifications = document.getElementById("navNotifications");

	if (Alpine.store("notify").firstLoad) {
		navNotifications.innerHTML = "";
	}

	let newest = Alpine.store("notify").newestDate;

	let response;
	try {
		response = await fetch("/api/notifications/newest?newest="+newest, {
			method: "GET"
		});
	} catch (err) {

		if (Alpine.store("notify").firstLoad) {
			navNotifications.innerText = "error loading notifications";
		}
		
		console.log("error getting notifications...:", err);

		// wait 2 seconds before trying again
		await new Promise(resolve => setTimeout(resolve, 2000));
		await loadNotifications();
		return
	}

	const json = await response.json();

	if (json.Notifs.length === 0 && Alpine.store("notify").firstLoad) {
		navNotifications.innerText = "no new notifications";
		Alpine.store("notify").new = false;
	} else {
		for (let i = 0; i < json.Notifs.length; i++) {
			let notif = document.createElement("p");
			notif.innerText = json.Notifs[i].Message;
			notif.setAttribute("notifURL", json.Notifs[i].URL);
			notif.setAttribute("notifID", json.Notifs[i].ID);
	
			console.log("notification is:", json.Notifs[i]);
	
			navNotifications.append(notif);
	
			notif.addEventListener("click", function(event) {
				console.log("click notification!");
				console.log("url is:", this.getAttribute("notifURL"));
				console.log("id is:", this.getAttribute("notifID"));

				// set notification as read
				notifSetRead(this.getAttribute("notifID"));

				// go to notification url
				window.location = this.getAttribute("notifURL");

				// TODO allow button for opening in new page
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