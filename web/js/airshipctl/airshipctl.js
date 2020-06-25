/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

// Splitting up scripts into sub components to keep it the functions in logical divisions
var includes = ["js/airshipctl/config.js","js/airshipctl/baremetal.js","js/airshipctl/document.js"];

// append CTL scripts to the page, this is an independent content loaded listener from the one in common.js
document.addEventListener("DOMContentLoaded", () => {
    for (let include of includes) {
        var script = document.createElement("script");
        script.src = include;
        document.head.appendChild(script);
    }
}, false)

// Displays the alerts from the backend
function handleCTLResponse(json) { // eslint-disable-line no-unused-vars
    let message = json["type"] + " " + json["component"] + " " + json["subComponent"] + " ";
    if (json.hasOwnProperty("error")) {
        showDismissableAlert("danger", message + json["error"], false);
    } else {
        showDismissableAlert("info", message + json["message"], true);
    }
}

function ctlGetDefaults(element) { // eslint-disable-line no-unused-vars
    let id = String(element.id);

    var json = { "type": "airshipctl",  "subComponent": "getDefaults" };
    switch(id) {
        case "LiBaremetal": json = Object.assign(json, { "component": "baremetal" }); break;
        case "LiConfig": json = Object.assign(json, { "component": "config" }); break;
        case "LiDocument": json = Object.assign(json, { "component": "document" }); break;
    }

    ws.send(JSON.stringify(json));
}

function displayCTLInfo(json) { // eslint-disable-line no-unused-vars
    if (json.hasOwnProperty("html")) {
        document.getElementById("DashView").style.display = "none";
        let div = document.getElementById("ContentDiv");
        div.style.display = "";
        div.innerHTML =  json["html"];
        if (!! document.getElementById("DocOverviewDiv") && json.hasOwnProperty("data")) {
            insertGraph(json["data"]);
        }
    } else {
        if (json.hasOwnProperty("error")) {
            showDismissableAlert("danger", json["error"], false);
        }
    }
}

function buttonHelper(id,text,disabled) { // eslint-disable-line no-unused-vars
    let button = document.getElementById(id);
    button.innerText = text;
    button.disabled = disabled;
}