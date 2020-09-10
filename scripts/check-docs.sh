#!/bin/bash
#
# Hacky shell script to find potentially missing documentation. 
# Relies on the fact that, so far, most of the source code 
# filenames mirror their respective doc filename. Does not do any 
# parsing (of either golang or markdown) to try and build a 100% 
# accurate list, the intention is this quickly gets you a short list
# of 10-15 things to check manually. 
#
# Usage: 
#    $ cd <git-wd-root>
#    $ ./scripts/check-docs.sh


for code_file in $(ls provider)
do

	if [[ "$code_file" == *"test.go" ]]; then
	  # echo "$code_file is a test file, ignoring"
	  continue
	fi

	doc_file=$(echo $code_file | sed -e 's|\.go|\.md|g')
	if [[ "$code_file" == "resource_"* ]]; then
		doc_file=$(echo $doc_file | sed -e 's|^resource_|docs/resources/|g')
		# echo "Looking for $doc_file"
		if [[ ! -f "$doc_file" ]]; then
   			 echo "$doc_file does not exist."
		fi

		continue
	fi

	if [[ "$code_file" == "data_source_"* ]]; then
		doc_file=$(echo $doc_file | sed -e 's|^data_source_|docs/data-sources/|g')
		# echo "Looking for $doc_file"
		if [[ ! -f "$doc_file" ]]; then
   			 echo "$doc_file does not exist."
		fi

		continue
	fi

	echo "No doc file found automatically for $code_file"

done