var http = require("http");

http.createServer(function(request,response){
    response.writeHeader(200, {"Content-Type": "text/plain"});
    response.write("Hello from nodejs !");
    response.end();
}).listen(8080);

console.log("Server Running on port 8080");
