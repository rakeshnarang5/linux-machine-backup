var querystring = require('querystring');
var http = require('http');
var fs = require('fs');

function PostCode(codestring) {
  var post_data = querystring.stringify({
    'js_code': codestring
  });

  var post_options = {
    host: 'http://shrouded-beyond-17924.herokuapp.com/al',
    port: '80',
    path: '/',
    method: 'POST',
    headers: {
      'Content-type': 'application/json'
    }
  };

  var post_req = http.request(post_options, function(res) {
    console.log("request successful")
  });

  post_req.write(post_data);
  post_req.end();
}

fs.readFile('sampledata.json', 'utf-8', function(err,data) {
  if (err) {
    console.log(err);
    process.exit(-2);
  }
  if (data) {
    PostCode(data);
  } else {
    console.log("no data");
    process.exit(-1);
  }
});