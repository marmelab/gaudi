var cradle      = require('cradle');
var connection  = new(cradle.Connection)(process.env.COUCHDB_PORT_5984_TCP_ADDR, 5984, {
    cache: true,
    raw: false,
    forceSave: true
});

var db = connection.database('books');
db.create();


db.save({
    title: 'Harry Potter and the Goblet of Fire', author: 'J. K. Rowling'
}, function (err, res) {
    if (err) {
        return console.error(err);
    }

    console.log(res);
});
