
var module = require('./queryAll.js');
var http = require('http');
var querystring = require('querystring');
var bodyParser = require('body-parser');
var json = require('JSON');

var server = http.createServer(function(request, response){
  var postdata = '';
    request.on('data', function(data){
   	 postdata = postdata + data;
    });
    request.on('end',function(){
   	var request=json.parse(postdata);
	console.log(postdata);
	
  
    const request1 = {
        chaincodeId: 'iot_federation',
        fcn: 'QueryAllEvents',
        args: [request.endKey]
    };
  	console.log(request1);

  	module.dataprocess(request1);
    });
    exports.dataprocess = function(data){
        response.writeHead(200);
        response.end(data);
    }
});


server.listen(8081, function(){
    console.log('Server is running...');
});


