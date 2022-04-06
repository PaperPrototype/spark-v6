async function SendMessage(messageText) {
	let messagesMount = document.getElementById("messagesMount");

	// message
	let messageNode = document.createElement("h4");
	messageNode.style = "padding: 1rem; display:flex; flex-wrap:nowrap;";
	messageNode.setAttribute("hideOnClick", "");
	messageNode.innerHTML = 
	`<span>` + messageText + `</span>` + ` <i style="margin-left:auto;" class="fa-solid fa-xmark"></i>`;

	messagesMount.appendChild(messageNode);

	// scroll up so user see's the message?
	window.scrollTo(0, 0);

	// set event listener for onclick to hide message
	setMessagesHideWhenClicked();
}

document.addEventListener("DOMContentLoaded", function(event) {
	setMessagesHideWhenClicked();
});

// make it so that when you click something with the message attribute it goes away
function setMessagesHideWhenClicked() {
	let messages = document.querySelectorAll("[hideOnClick]");

	for (let i = 0; i < messages.length; i++) {
		messages[i].addEventListener("click", function() {
			console.log("clicked");
			this.style = "display: none;";
		});
	}
}