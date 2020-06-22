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

// add the footer and header when the page loads
if (document.addEventListener) {
    document.addEventListener("DOMContentLoaded", function () {
        const webview = document.querySelector("webview");
        const spinner = document.querySelector(".spinner");
        const loadStart = () => {
            spinner.style.display = "block";
        }
        const loadStop = () => {
            spinner.style.display = "none";
        }
        const loadFail = (err) => {
            showDismissableAlert("danger", `Error loading '${err.validatedURL}': ${err.errorDescription} (${err.errorCode})`, true);
        }
        webview.addEventListener("did-start-loading", loadStart);
        webview.addEventListener("did-stop-loading", loadStop);
        webview.addEventListener("did-fail-load", loadFail);
    }, false);
}

// add dashboard links to Dropdown if present in $HOME/.airship/airshipui.json
function addServiceDashboards(json) { // eslint-disable-line no-unused-vars
    if (json !== undefined) {
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
                    let url = `${dash.protocol}://${fqdn}:${dash.port}/${dash.path || ""}`;
                    addDashboard("DashDropdown", dash.name, url)
                }
            }
        }
    }
}

// if any plugins (external executables) have a corresponding web dashboard defined,
// add them to the dropdown
function addPluginDashboards(json) { // eslint-disable-line no-unused-vars
    if (json !== undefined) {
        for (let i = 0; i < json.length; i++) {
            if (json[i].executable.autoStart && json[i].dashboard !== undefined) {
                let dash = json[i].dashboard;
                let url = `${dash.protocol}://${dash.fqdn}:${dash.port}/${dash.path || ""}`;
                addDashboard("PluginDropdown", json[i].name, url);
            }
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
        document.getElementById("ContentDiv").style.display = "none";
        document.getElementById("DashView").style.display = "";
    }
    a.appendChild(span);
    li.appendChild(a);
    nav.appendChild(li);
}

function authenticate(json) { // eslint-disable-line no-unused-vars
    // use webview to display the auth page
    let view = document.getElementById("AuthView");
    view.src = json["url"];
}

// show a dismissable alert in the UI
// if 'fade' is set to true, the alert will fade out and disappear after 5s
function showDismissableAlert(alertLevel, msg, fade) { // eslint-disable-line no-unused-vars
    let e = document.getElementById("alert-div");
    let alertId = `alert-${Math.floor(Math.random() * 1000)}`;
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
    div.id = alertId;
    div.className = `alert alert-${alertLevel} alert-dismissable fade show`;
    div.setAttribute("role", "alert");
    div.innerHTML = `<strong>${alertHeading}: </strong>${msg}`;
    div.style.opacity = "1";

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

    // add auto-hide if fade is true
    if (fade === true) {
        let script = document.createElement("script");
        let inline = `alertFadeOut("${alertId}")`;
        script.innerText = inline;
        div.appendChild(script);
    }

    e.appendChild(div);
}

function alertFadeOut(id) { // eslint-disable-line no-unused-vars
    let element = document.getElementById(id);
    setTimeout(function() {
        element.style.transition = "opacity 2s ease";
        element.style.opacity = "0";
    }, 5000);
    element.addEventListener("transitionend", function() {
        element.parentNode.removeChild(element);
    });
}

function enableAccordion() { // eslint-disable-line no-unused-vars
    var acc = document.getElementsByClassName("accordion");
    var i;

    for (i = 0; i < acc.length; i++) {
        acc[i].addEventListener("click", function () {
            this.classList.toggle("active");
            var panel = this.nextElementSibling;
            if (panel.style.maxHeight) {
                panel.style.maxHeight = null;
            } else {
                panel.style.maxHeight = panel.scrollHeight + "px";
            }
        });
    }
}