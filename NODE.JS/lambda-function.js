var http = require("http");

// default signature of a lambda function
exports.handler = (event, context, callback) => {
    //callback(null, {"Message":"Heisenberg"}); //SUCCESS WITH MESSAGE
    try {
        if (event.session.new) {
            console.log("NEW SESSION")
        }
        
        switch (event.request.type) {
            case "LaunchRequest":
                console.log(`LAUNCH REQUEST`)
                
                var options = {
                    hostname: 'shrouded-beyond-17924.herokuapp.com',
                    port: 80,
                    path: '/',
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    }                 
                };

                var req = http.request(options, function(res) {
                    console.log('Status: ' + res.statusCode);
                    console.log('Headers: ' + JSON.stringify(res.headers));
                    res.setEncoding('utf8');
                    res.on('data', function (body) {
                    console.log('Body: ' + body);
                    context.succeed(
                        generateResponse(
                            buildSpeechletResponse(body, true),
                            {}
                        )
                    )
                        });
                    });
                
                req.on('error', function(e) {
                console.log('problem with request: ' + e.message);
                });
                // write data to request body
                req.write('{"string": "Hello, World"}');
                req.end();
                
                
                
                break;
            case "IntentRequest":
                console.log(`INTENT REQUEST`)
                break;
            case "SessionEndedRequest":
                console.log(`SESSION ENDED REQUEST`)
                break;
            default:
                context.fail(`INVALID REQUEST TYPE ${event.request.type}`)
        }
    } catch (error) { context.fail(`Exception: ${error}`)}
    
    
}

buildSpeechletResponse = (outputText, shouldEndSession) => {
    return {
        outputSpeech: {
            type: "PlainText",
            text: outputText
        },
        shouldEndSession: shouldEndSession
    }
}

generateResponse = (speechletResponse, sessionAttributes) => {
    return {
        version: "1.0",
        sessionAttributes: sessionAttributes,
        response: speechletResponse
    }
}