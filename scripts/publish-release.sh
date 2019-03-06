#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "$0")"

go get github.com/tcnksm/ghr

if [[ ${CIRCLE_TAG} == *"rc"* ]]; then
	ghr \
		-u ${CIRCLE_PROJECT_USERNAME} \
		-r ${CIRCLE_PROJECT_REPONAME} \
		-n "v${CIRCLE_TAG}" \
		-prerelease \
		-replace \
		${CIRCLE_TAG} ../artifacts
else
	releaseDate=$(date '+%B-%-d-%Y' | tr '[:upper:]' '[:lower:]')
	releaseVersion=$(echo ${CIRCLE_TAG} | tr -d '.')

	ghr \
		-u ${CIRCLE_PROJECT_USERNAME} \
		-r ${CIRCLE_PROJECT_REPONAME} \
		-n "v${CIRCLE_TAG}" \
		-b "[Release Notes](https://github.com/mrparkers/terraform-provider-keycloak/blob/master/CHANGELOG.md#${releaseVersion}-${releaseDate})" \
		-replace \
		${CIRCLE_TAG} ../artifacts

	sudo apt-get update && sudo apt-get install mkdocs
	cd .. && mkdocs gh-deploy
fi
