local keycloakTestEnv() = {
	KEYCLOAK_CLIENT_ID: "terraform",
	KEYCLOAK_CLIENT_SECRET: "884e0f95-0f42-4a63-9b1f-94274655669e",
	KEYCLOAK_CLIENT_TIMEOUT: "5",
	KEYCLOAK_URL: "http://localhost:8080",
	KEYCLOAK_REALM: "master",
};

local pipeline(version) = {
	kind: 'pipeline',
	type: 'kubernetes',
	name: 'test-%(version)s' % { version: version },
	services: [
		{
			name: 'keycloak',
			image: 'jboss/keycloak:%(version)s' % { version: version },
			environment: {
				"DB_VENDOR": "H2",
				"KEYCLOAK_LOGLEVEL": "DEBUG",
				"KEYCLOAK_USER": "keycloak",
				"KEYCLOAK_PASSWORD": "password",
			},
			commands: [
				'/opt/jboss/tools/docker-entrypoint.sh -b 0.0.0.0 -Dkeycloak.profile.feature.upload_scripts=enabled'
			],
		},
	],
	steps: [
		{
			name: 'fetch dependencies',
			image: 'mrparkers/golang-test:1.13.11',
			volumes: [{
				name: "deps",
				path: "/go"
			}],
			commands: [
				'go mod download',
			]
		},
		{
			name: 'setup',
			image: 'mrparkers/golang-test:1.13.11',
			commands: [
				'./scripts/wait-for-local-keycloak.sh',
				'./scripts/create-terraform-client.sh',
			],
			environment: keycloakTestEnv(),
		},
		{
			name: 'example tf',
			image: 'mrparkers/golang-test:1.13.11',
			volumes: [{
				name: "deps",
				path: "/go"
			}],
			commands: [
				'mkdir -p ~/.terraform.d/plugins',
				'make build && mv terraform-provider-keycloak ~/.terraform.d/plugins',
				'cd example',
				'terraform init',
				'terraform apply -auto-approve',
				'terraform destroy -auto-approve'
			],
			environment: keycloakTestEnv(),
		},
		{
			name: 'test',
			image: 'mrparkers/golang-test:1.13.11',
			volumes: [{
				name: "deps",
				path: "/go"
			}],
			commands: [
				'make testacc',
			],
			environment: keycloakTestEnv(),
		},
	],
	trigger: {
		event: [
			'pull_request',
		],
	},
	volumes: [{
		name: "deps",
		temp: {},
	}],
};

[
	pipeline('10.0.2'),
]
