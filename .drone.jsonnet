{
	kind: 'pipeline',
	type: 'kubernetes',
	name: 'default',
	steps: [
		{
			name: 'test',
			image: 'alpine',
			commands: [
				'ls -alh',
			]
		}
	],
}
