/*
 Copyright (c) 2020 AT&T. All Rights Reserved.
 
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

// add the footer and header when the page loads
if (document.addEventListener) {
    document.addEventListener("DOMContentLoaded", function () {
        window.onscroll = function () {
            let header = document.getElementById("HeaderDiv");
            let sticky = header.offsetTop;

            if (window.pageYOffset > sticky) {
                header.classList.add("sticky");
            } else {
                header.classList.remove("sticky");
            }
        };
    }, false);
}

// add dashboard links to Plugins if present in $HOME/.airshipui/plugins.json
function addPlugins(json) { // eslint-disable-line no-unused-vars
    let dropdown = document.getElementById("PluginDropdown");
    for (let i = 0; i < json.length; i++) {
        let dash = json[i];

        let a = document.createElement("a");
        a.innerText = dash["name"];

        // created as a lambda in order to prevent auto firing the onclick event
        a.onclick = () => {
            let view = document.getElementById("DashView");
            view.src = dash["url"];

            document.getElementById("MainDiv").style.display = "none";
            document.getElementById("DashView").style.display = "";
        }

        dropdown.appendChild(a);
    }
}

function authenticate(json) { // eslint-disable-line no-unused-vars
    // use webview to display the auth page
    let view = document.getElementById("DashView");
    view.src = json["url"];

    document.getElementById("MainDiv").style.display = "none";
    document.getElementById("DashView").style.display = "";
}

function removeElement(id) { // eslint-disable-line no-unused-vars
    if (document.contains(document.getElementById(id))) {
        document.getElementById(id).remove();
    }
}