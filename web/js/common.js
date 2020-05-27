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

// add dashboard links to Dropdown if present in $HOME/.airship/airshipui.json
function addServiceDashboards(json) { // eslint-disable-line no-unused-vars
    for (let i = 0; i < json.length; i++) {
        let cluster = json[i];
        for (let j = 0; j < cluster.namespaces.length; j++) {
            let namespace = cluster.namespaces[j];
            for (let k = 0; k < namespace.dashboards.length; k++) {
                let dash = namespace.dashboards[k];
                let fqdn = "";
                if (dash.fqdn === undefined) {
                    fqdn = `${dash.hostname}.${cluster.namespaces[j].name}.${cluster.baseFqdn}`
                } else {
                    ({ fqdn } = dash.fqdn);
                }
                let url =  `${dash.protocol}://${fqdn}:${dash.port}/${dash.path}`;
                addDashboard("DashDropdown", dash.name, url)
            }
        }
    }
}

// if any plugins (external executables) have a corresponding web dashboard defined,
// add them to the dropdown
function addPluginDashboards(json) { // eslint-disable-line no-unused-vars
    for (let i = 0; i < json.length; i++) {
        if (json[i].executable.autoStart && json[i].dashboard.fqdn !== "") {
            let dash = json[i].dashboard;
            let url = `${dash.protocol}://${dash.fqdn}:${dash.port}/${dash.path}`;
            addDashboard("PluginDropdown", json[i].name, url)
        }
    }
}

function addDashboard(navElement, name, url) {
    let nav = document.getElementById(navElement);
    let li = document.createElement("li");
    li.className = "c-sidebar-nav-item";
    let a = document.createElement("a");
    a.className = "c-sidebar-nav-link";
    let span = document.createElement("span");
    span.className = "c-sidebar-nav-icon";
    a.innerText = name;
    a.onclick = () => {
        let view = document.getElementById("DashView");
        view.src = url;
        document.getElementById("DashView").style.display = "";
    }
    a.appendChild(span);
    li.appendChild(a);
    nav.appendChild(li);
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

function showDismissableAlert(alertLevel, msg) { // eslint-disable-line no-unused-vars
    let e = document.getElementById("alert-div");
    let alertHeading = "";

    switch (alertLevel) {
        case "danger":
            alertHeading = "Error";
            break;
        case "warning":
            alertHeading = "Warning";
            break;
        default:
            alertHeading = "Info";
    }

    let div = document.createElement("div");
    div.className = `alert alert-${alertLevel} alert-dismissable fade show`;
    div.setAttribute("role", "alert");
    div.innerHTML = `<strong>${alertHeading}: </strong>${msg}`;

    // dismissable button
    let btn = document.createElement("button");
    btn.className = "close";
    btn.type = "button";
    btn.setAttribute("data-dismiss", "alert");
    btn.setAttribute("aria-label", "Close");

    let span = document.createElement("span");
    span.setAttribute("aria-hidden", "true");
    span.innerText = "×";

    btn.appendChild(span);
    div.appendChild(btn);

    e.appendChild(div);
}