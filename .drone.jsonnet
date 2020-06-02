local pipeline(version) = {
	kind: 'pipeline',
	type: 'kubernetes',
	name: 'test-%(version)s' % { version: version },
	services: [
		{
			name: 'keycloak',
			image: 'jboss/keycloak:%(version)s' % { version: version },
		},
	],
	steps: [
		{
			name: 'create test file',
			image: 'alpine',
			commands: [
				'echo "hi" >> test.txt',
			]
		},
		{
			name: 'print test file',
			image: 'alpine',
			commands: [
				'cat test.txt',
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
	pipeline('7.0.1'),
]
