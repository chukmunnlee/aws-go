exports.handler = function(event, context, callback) {
	//console.log('DEBUG: context: ', context)
	console.log('DEBUG: event: ', event)
	console.log(`DEBUG: AWS_REGION: ${process.env.HELLO_AWS_REGION}`)
	console.log(`DEBUG: BUCKET_NAME: ${process.env.HELLO_BUCKET_NAME}`)
	//callback(null, `hello ${JSON.stringify(event)}`)
	callback(null, 
		{ 
			statusCode: 200, 
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify({ message: `Hello ${new Date()}` })
		}
	);
}
