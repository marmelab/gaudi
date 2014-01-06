<?php
$dbInfos = parse_url(getenv('DB_PORT'));

$dsn = 'mysql:dbname=users;host='.$dbInfos['host'];

$connection = new PDO($dsn, 'root', '');

$query = 'SELECT * FROM users';
$results = $connection->query($query);

foreach ($results as $result) {
    var_dump($result['username']);
}
