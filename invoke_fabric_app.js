var module = require('./invoke.js');
var http = require('http');
var querystring = require('querystring');
var bodyParser = require('body-parser');
var json = require('JSON');

var server = http.createServer(function(request, response){
  var postdata = '';
	response.writeHead(200, {'Content-Type' : 'text/plain'});
	response.write("h");
	response.end();
    request.on('data', function(data){
	console.log(typeof(data));
   	 postdata = postdata + data;
    });
    request.on('end',function(){
	console.log(postdata);
   	var request=json.parse(postdata);
	
	
  	const request1 = {
   		 //targets : --- letting this default to the peers assigned to the channel
   		 chaincodeId: 'iot_federation',
   		 fcn: 'CreateLedger',
   		 args: [request.schoolName, request.deviceIp, request.time, request.con],
    	chainId: 'mychannel',
    	txId: ''
   	 }; 

  	console.log(request1);

  	module.dataprocess(request1);
	});
});


server.listen(8080, function(){
    console.log('Server is running...');
});


