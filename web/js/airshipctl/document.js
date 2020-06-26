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

var editor = null;
var editorContents = null;

function documentAction(element) { // eslint-disable-line no-unused-vars
    let elementId = element.id;

    // change text & disable the button while the process happens
    buttonHelper(elementId, "In Progress", true);

    var json = { "type": "airshipctl", "component": "document" };
    switch(elementId) {
        case "DocPullBtn": Object.assign(json, { "subComponent": "docPull" }); break;
        case "KubeConfigBtn":
            Object.assign(json, { "subComponent": "yaml" });
            Object.assign(json, { "message": "kube" });
            break;
        case "AirshipConfigBtn":
            Object.assign(json, { "subComponent": "yaml" });
            Object.assign(json, { "message": "airship" });
            break;
        case "SaveYamlBtn":
            Object.assign(json, { "subComponent": "yamlWrite" });
            Object.assign(json, { "message": editorContents });
            Object.assign(json, { "yaml": window.btoa(editor.getValue()) });
            break;
    }
    ws.send(JSON.stringify(json));
}

function ctlParseDocument(json) { // eslint-disable-line no-unused-vars
    switch(json["subComponent"]) {
        case "getDefaults": displayCTLInfo(json); addFolderToggles(); break;
        case "yaml": insertEditor(json); break;
        case "yamlWrite": insertEditor(json); buttonHelper("SaveYamlBtn", "Save", true); break;
        case "docPull": buttonHelper("DocPullBtn", "Document Pull",false); handleCTLResponse(json); break;
        default: handleCTLResponse(json)
    }
}

// adds the monaco editor to the UI and populates it with yaml
function insertEditor(json) {
    // dispose of any detritus that may not have been disposed of before reuse
    if (editor !== null) { editor.dispose(); editorContents = null; }

    // disable the save button if it's not already
    let saveBtn = document.getElementById("SaveYamlBtn");
    saveBtn.disabled = true;

    // create and populate the monaco editor
    let div = document.getElementById("DocYamlDIV");

    editor = monaco.editor.create(div, {
        value:  window.atob(json["yaml"]),
        language: "yaml",
        automaticLayout: true
    });

    toggleDocument();

    // toggle the buttons back to the original message
    switch(json["message"]) {
        case "kube":
            buttonHelper("KubeConfigBtn", " - kubeconfig", false);
            editorContents = "kube";
            document.getElementById("KubeConfigSpan").classList.toggle("document-open");
            break;
        case "airship":
            buttonHelper("AirshipConfigBtn", " - config", false);
            editorContents = "airship";
            document.getElementById("AirshipConfigSpan").classList.toggle("document-open");
            break;
    }

    // on change event for the editor
    editor.onDidChangeModelContent(() => {
        saveBtn.disabled = false;
    });
}

function addFolderToggles() {
    var toggler = document.getElementsByClassName("folder");
    for (let i = 0; i < toggler.length; i++) {
        toggler[i].addEventListener("click", function() {
            this.parentElement.querySelector(".nested").classList.toggle("active");
            this.classList.toggle("folder-open");
        });
    }
}

function toggleDocument() {
    var toggler = document.getElementsByClassName("document");
    for (let i = 0; i < toggler.length; i++) {
        toggler[i].className = "document";
    }
}
