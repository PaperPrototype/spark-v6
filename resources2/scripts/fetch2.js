/* fetch utility abstracting away repetitive code
    - mainly automatic error handling
        - code can still check if there was an error, and decide whether to procede or not
*/
function fetch2(route, method, code) {
    fetch(route, {
        method: method,
    })
    .then(function(resp) {
        if (!resp.ok) {

            var json = resp.json();
            if (json.Error !== "" && json.Error !== undefined)
            {
                throw json.Error;
            }

            throw "Error with response";
        }

        return resp.json();
    })
    .then(function(json) {
        if (json.Error !== "") {
            SendMessage(json.Error);
        }

        code(json);
    })
    .catch(function(err) {
        SendMessage(err);
        console.error(err);
    });
}

function postFetch2(route, body, code) {
    fetch(route, {
        method: "POST",
        body: body,
    })
    .then(function(resp) {
        if (!resp.ok) {

            var json = resp.json();
            if (json.Error !== "" && json.Error !== undefined)
            {
                throw json.Error;
            }

            throw "Error with response";
        }

        return resp.json();
    })
    .then(function(json) {
        if (json.Error !== "") {
            SendMessage(json.Error);
        }

        code(json);
    })
    .catch(function(err) {
        SendMessage(err);
        console.error(err);
    });
}