#!/bin/bash

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

declare FILES_MISSING_COPYRIGHT=()

# Get the files of interst
getFiles() {
  FILES=$(find -L . -type f \( -iname \*.go -o -iname \*.yaml -o -iname \*.yml -o -iname \*.sh -o -iname \*.ts -o -iname \*.css \) \
    -not -path "./etc" \
    -not -path "./client/dist/*" \
    -not -path "./client/node_modules/*" \
    -not -path "./tools/*node*" \
    | grep -v "testdata" \
    | grep -v "manifests")

  for each in $FILES
  do
    if ! grep -Eq 'Apache License|License-Identifier: Apache' $each
    then
      FILES_MISSING_COPYRIGHT+=($each)
    fi
  done
}

# Find all files we care about and add licenses if needed
addLicense() {
  getFiles

  if [ ${#FILES_MISSING_COPYRIGHT[@]} -gt 0 ]
  then
    for each in $FILES
    do
      if ! grep -Eq 'Apache License|License-Identifier: Apache' $each
      then
        echo "Adding license header to $each"
        filename=$(basename -- "$each")
        case ${filename##*.} in
          css)
            cat tools/license_html.txt $each >$each.new
            ;;
          go)
            cat tools/license_go.txt $each >$each.new
            ;;
          sh)
            head -n 1 $each >>$each.new
            NUM_OF_LINES=$(< "tools/license_bash.txt" wc -l)
            head -n $NUM_OF_LINES tools/license_bash.txt >>$each.new
            tail -n+2 $each >>$each.new
            ;;
          ts)
            cat tools/license_html.txt $each >$each.new
            ;;
          yaml)
            cat tools/license_yaml.txt $each >$each.new
            ;;
          yal)
            cat tools/license_yaml.txt $each >$each.new
            ;;
        esac
        mv $each.new $each
      fi
    done
  else
    echo "no file with missing copyright header detected, make target completed successfully"
  fi
}

# Find all files we care about and check if the license is there
checkLicense() {
  getFiles

  if [ ${#FILES_MISSING_COPYRIGHT[@]} -gt 0 ]
  then
    printf "Copyright header missing for: %s\n" "${FILES_MISSING_COPYRIGHT[@]}"
    echo "please run make add-copyright"
    exit 1
  else
    echo "no file with missing copyright header detected, make target completed successfully"
  fi
}

case ${1} in
add)
    addLicense
    ;;
check)
    checkLicense
    ;;
*)
    echo "usage: ${0} { add | check }"
    exit 1
    ;;
esac