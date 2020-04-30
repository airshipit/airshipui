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

const remote = require('electron').remote;
const app = remote.app;
const fs = require('fs')
const config = require('electron-json-config')


// add the footer and header when the page loads
if (document.addEventListener) {
	document.addEventListener("DOMContentLoaded", function() {
        window.onscroll = function() {
            let header = document.getElementById("HeaderDiv");
            let sticky = header.offsetTop;

            if (window.pageYOffset > sticky) {
                header.classList.add("sticky");
            } else {
                header.classList.remove("sticky");
            }
        };
        addPlugins()
    }, false);
}

// add dashboard links to Plugins if present in $HOME/.airshipui/plugins.json
function addPlugins() {
    try {
        var f = fs.readFileSync(app.getPath('home') + '/.airshipui/plugins.json')

        var dashboards = JSON.parse(f)

        for (var i = 0; i < dashboards.external_dashboards.length; i++) {
            var path = app.getAppPath()
            var p = document.getElementById("plugins")
            var a = document.createElement("a")

            var dash = dashboards.external_dashboards[i]
            config.set(i.toString(), dash.url)
            a.setAttribute('href', `${path}/plugins/dashboards/index.html`)
            a.setAttribute('onclick', `config.set('dashboard', ${i.toString()})`)
            a.innerText = dash.name
            p.appendChild(a)
        }

    } catch (e) {
        console.log("Plugins file not found")
    }
}

function removeElement(id) {
    if (document.contains(document.getElementById(id))) {
        document.getElementById(id).remove();
    }
}

// based on w3school: https://www.w3schools.com/howto/howto_js_sort_table.asp
function sortTable(tableID, column) {
    var table, rows, switching, i, x, y, shouldSwitch, dir, switchcount = 0;
    table = document.getElementById(tableID);
    switching = true;
    // Set the sorting direction to ascending:
    dir = "asc";
    /* Make a loop that will continue until
    no switching has been done: */
    while (switching) {
        // Start by saying: no switching is done:
        switching = false;
        rows = table.rows;
        /* Loop through all table rows (except the
        first, which contains table headers): */
        for (i = 1; i < (rows.length - 1); i++) {
            // Start by saying there should be no switching:
            shouldSwitch = false;
            /* Get the two elements you want to compare,
            one from current row and one from the next: */
            x = rows[i].getElementsByTagName("TD")[column];
            y = rows[i + 1].getElementsByTagName("TD")[column];

            if (x !== undefined && y !== undefined) {
                /* Check if the two rows should switch place,
                based on the direction, asc or desc: */
                if (dir == "asc") {
                    if (x.innerHTML.toLowerCase() > y.innerHTML.toLowerCase()) {
                        // If so, mark as a switch and break the loop:
                        shouldSwitch = true;
                        break;
                    }
                } else if (dir == "desc") {
                    if (x.innerHTML.toLowerCase() < y.innerHTML.toLowerCase()) {
                        // If so, mark as a switch and break the loop:
                        shouldSwitch = true;
                        break;
                    }
                }
            }
        }
        if (shouldSwitch) {
            /* If a switch has been marked, make the switch
            and mark that a switch has been done: */
            rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
            switching = true;
            // Each time a switch is done, increase this count by 1:
            switchcount++;
        } else {
            /* If no switching has been done AND the direction is "asc",
            set the direction to "desc" and run the while loop again. */
            if (switchcount == 0 && dir == "asc") {
                dir = "desc";
                switching = true;
            }
        }
    }
}