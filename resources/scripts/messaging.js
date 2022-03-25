async function SendMessage(messageText) {
	let messagesMount = document.getElementById("messagesMount");

	// message
	let messageNode = document.createElement("h4");
	messageNode.style = "padding: 1rem;";
	messageNode.setAttribute("hideOnClick", "");
	messageNode.innerText = messageText;

	messagesMount.appendChild(messageNode);

	// set event listenr for onclick to hide
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