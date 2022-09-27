function SendMessage(message) {
    var messages = document.getElementById("messages");

    const div = document.createElement("div");
    div.innerHTML = `<div x-data="{ open: true }" x-show="open" class="thm-bg-hl" style="padding:1rem;">
        <div style="display:flex; flex-direction:row;" @click="open = ! open">
            <span>` + message + `</span> <i style="margin-left:auto;" class="fa-solid fa-xmark"></i>
        </div>
    </div>`;

    messages.append(div);
}