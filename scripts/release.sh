#!/usr/bin/env bash

set -eu

cd "$(dirname "$0")"

for config in $(cat release-targets.json | jq -rc '.[]'); do
	os=$(echo ${config} | jq -r '.os')
	platform=$(echo ${config} | jq -r '.platform')

	echo "Building for ${os}_${platform}..."

	GOOS=${os} GOARCH=${platform} go build -o terraform-provider-keycloak_v${VERSION} ..
	zip terraform-provider-keycloak_v${VERSION}_${os}_${platform}.zip terraform-provider-keycloak_v${VERSION} ../LICENSE
	rm terraform-provider-keycloak_v${VERSION}
done;
