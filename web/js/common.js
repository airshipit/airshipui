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

var graph = null;

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

function tabAction(event, element) { // eslint-disable-line no-unused-vars
    // Declare all variables
    var i, tabcontent, tablinks;

    // Get all elements with class="tabcontent" and hide them
    tabcontent = document.getElementsByClassName("tabcontent");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }

    // Get all elements with class="tablinks" and remove the class "active"
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }

    // Show the current tab, and add an "active" class to the button that opened the tab
    let id = String(element.id);
    let div = id.replace("Btn","");
    switch (id) {
        case "DocOverviewTabBtn": document.getElementById(div).style.display = "block"; break;
        case "DocPullTabBtn": document.getElementById(div).style.display = "block"; break;
        case "YamlTabBtn": document.getElementById(div).style.display = "block"; break;
    }

    event.currentTarget.className += " active";
}

function insertGraph(data) { // eslint-disable-line no-unused-vars
    if (graph !== null) { graph.destroy(); }

    // create a network
    var container = document.getElementById("DocOverviewDiv");

    // TODO: extract these to a constants file somewhere
    var options = {
        nodes: {
            shape: "box",
            scaling: {
                max: 200, min: 100
            }
        },
        physics: {
            forceAtlas2Based: {
                gravitationalConstant: -26,
                centralGravity: 0.005,
                springLength: 230,
                springConstant: 0.18,
                avoidOverlap: 1.5
            },
            maxVelocity: 146,
            solver: "forceAtlas2Based",
            timestep: 0.35,
            stabilization: {
                enabled: true,
                iterations: 1000,
                updateInterval: 25
            }
        }
    };
    graph = new vis.Network(container, data, options);
}

// add dashboard links to Dropdown if present in $HOME/.airship/airshipui.json
function addDashboards(json) { // eslint-disable-line no-unused-vars
    if (json !== undefined) {
        for (let i = 0; i < json.length; i++) {
            let dashboard = json[i];
            let url = `${dashboard.baseURL}/${dashboard.path || ""}`;
            addDashboard(dashboard.name, url, dashboard.isProxied)
        }
    }
}

function addDashboard(name, url, proxied) {
    let nav = document.getElementById("DashDropdown");
    let li = document.createElement("li");
    li.className = "c-sidebar-nav-item";
    let a = document.createElement("a");
    a.className = "c-sidebar-nav-link";
    let span = document.createElement("span");
    span.className = "c-sidebar-nav-icon";
    a.innerText = name;
    a.appendChild(span);
    li.appendChild(a);
    nav.appendChild(li);
    if (proxied) {
        a.target = "_blank";
        a.href = "javascript:window.open('" + url + "')";
    } else {
        a.onclick = () => {
            let view = document.getElementById("DashView");
            view.src = url;
            document.getElementById("ContentDiv").style.display = "none";
            document.getElementById("DashView").style.display = "";
        }
    }
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
    div.style.width = "350px";
    div.style.marginTop = "5px";

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