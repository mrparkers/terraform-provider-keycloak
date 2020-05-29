{
	kind: 'pipeline',
	type: 'kubernetes',
	name: 'test',
	steps: [
		{
			name: 'list contents',
			image: 'alpine',
			commands: [
				'ls -alh',
			]
		}
	],
	trigger: {
		event: [
			'pull_request'
		]
	}
}
