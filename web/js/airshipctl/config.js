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

function ctlGetConfig() { // eslint-disable-line no-unused-vars
    var json = { "type": "airshipctl", "component": "config",  "subComponent": "getDefaults" };
    ws.send(JSON.stringify(json));
}

function ctlParseConfig(json) { // eslint-disable-line no-unused-vars
    switch(json["subComponent"]) {
        case "getDefaults": displayConfigInfo(json); break;
        default: handleCTLResponse(json);
    }
}

function displayConfigInfo(json) {
    displayCTLInfo(json);
    enableAccordion();
}

function saveConfig(element) { // eslint-disable-line no-unused-vars
    var json = {
        "type": "airshipctl",
        "component": "config",
    };

    tableID = getTableId(element);
    while (element = element.parentNode) {
        if (element.tagName === "TR") {
            switch (tableID) {
                case "ClusterTable": json = Object.assign(json, saveCluster(element)); break;
                case "ClusterAddTable": json = Object.assign(json, addCluster(element)); break;
                case "ContextTable": json = Object.assign(json, saveContext(element)); break;
                case "ContextAddTable": json = Object.assign(json, addContext(element)); break;
                case "CredentialTable": json = Object.assign(json, saveCredential(element)); break;
                case "CredentialAddTable": json = Object.assign(json, addCredential(element)); break;
            }
            break;
        }
    }

    console.log("Save Config Request: ", json);
    ws.send(JSON.stringify(json));
}

function addCluster(row) {
    return {
        "subComponent": "cluster",
        "clusterOptions": {
            "Name": row.cells[0].children[0].value,
            "ClusterType": row.cells[1].children[0].value,
            "Server": row.cells[2].children[0].value
        }
    };
}

function saveCluster(row) {
    sa = row.cells[1].textContent.split("_");
    return {
        "subComponent": "cluster",
        "clusterOptions": {
            "Name": sa[0],
            "ClusterType": sa[1],
            "Server": row.cells[4].textContent
        }
    };
}

function addContext(row) {
    return {
        "subComponent": "context",
        "contextOptions": {
            "Name": row.cells[0].children[0].value,
            "ClusterType": row.cells[1].children[0].value,
            "Cluster": row.cells[2].children[0].value,
            "AuthInfo": row.cells[3].children[0].value
        }
    };
}

function saveContext(row) {
    sa = row.cells[0].textContent.split("_");
    return {
        "subComponent": "context",
        "contextOptions": {
            "Name": sa[0],
            "ClusterType": sa[1],
            "Cluster": row.cells[3].textContent,
            "AuthInfo": row.cells[4].textContent
        }
    };
}

function addCredential(row) {
    return {
        "subComponent": "credential",
        "authInfoOptions": {
            "Name": row.cells[0].children[0].value,
            "Username": row.cells[1].children[0].value
        }
    };
}

function saveCredential(row) {
    return {
        "subComponent": "credential",
        "authInfoOptions": {
            "Name": row.cells[0].textContent,
            "Username": row.cells[1].textContent
        }
    };
}

function getTableId(node) {
    var element = node;
    while (element.tagName !== "TABLE") {
        element = element.parentNode;
    }
    return element.id;
}

function saveConfigDialog(element) { // eslint-disable-line no-unused-vars
    saveConfig(element);
    setTimeout(function(){ ctlGetConfig(); }, 250);
    closeDialog(element);
}

// This will use the modal described in the pagelet that is sent via the websocket from the backend
function addConfigModal(element) { // eslint-disable-line no-unused-vars
    let elementId = element.id;
    var id, template;
    switch(elementId) {
        case "ClusterBtn": id = "AddCluster"; template = "ClusterModalTemplate"; break;
        case "ContextBtn": id = "AddContext"; template = "ContextModalTemplate"; break;
        case "CredentialBtn": id = "AddCredential"; template = "CredentialModalTemplate"; break;
    }

    let dialog = document.createElement("DIALOG");
    document.body.appendChild(dialog);
    dialog.setAttribute("id", id);
    dialog.innerHTML = document.getElementById(template).innerHTML;
    dialog.showModal();
}

function closeDialog(element) { // eslint-disable-line no-unused-vars
    while (element = element.parentNode) {
        if (element.tagName === "DIALOG") {
            element.close();
            document.body.removeChild(element);
            break;
        }
    }
}