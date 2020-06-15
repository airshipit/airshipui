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

function baremetalAction(element) { // eslint-disable-line no-unused-vars
    let elementId = element.id;

    // change text & disable the button while the process happens
    buttonHelper(elementId, "In Progress", true);

    var json = { "type": "airshipctl", "component": "baremetal" };
    switch(elementId) {
        case "GenIsoBtn": json = Object.assign(json, { "subComponent": "generateISO" }); break;
    }
    ws.send(JSON.stringify(json));
}

function ctlParseBaremetal(json) { // eslint-disable-line no-unused-vars
    switch(json["subComponent"]) {
        case "getDefaults": displayCTLInfo(json); break;
        case "generateISO": buttonHelper("GenIsoBtn", "Generate ISO",false); handleCTLResponse(json); break;
        default: handleCTLResponse(json)
    }
}
