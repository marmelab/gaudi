var http        = require('http');
var redis       = require("redis");
var redisHost   = process.env.REDIS_PORT_6379_TCP_ADDR;
var client      = redis.createClient(6379, redisHost);

client.on("error", function (err) {
    console.log("Redis error : " + err);
});

http.createServer(function (req, res) {
    res.writeHead(200, {'Content-Type': 'text/plain'});

    client.set("Fou", "barre");
    client.get("Fou", function(err, reply) {
        if (err) {
            return res.end(err)
        }

        res.end(reply)
    });
}).listen(80);
