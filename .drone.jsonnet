local pipeline(version) = {
	kind: 'pipeline',
	type: 'kubernetes',
	name: 'test-%(version)s' % { version: version },
	services: [
		{
			name: 'keycloak',
			image: 'jboss/keycloak:%(version)s' % { version: version },
		}
	],
	steps: [
		{
			name: 'list contents',
			image: 'alpine',
			commands: [
				'ls -alh',
			]
		},
		{
			name: 'sleep',
			image: 'alpine',
			commands: [
				'sleep 60',
			]
		},
	],
	trigger: {
		event: [
			'pull_request',
		],
	},
};

[
	pipeline('8.0.1'),
]
