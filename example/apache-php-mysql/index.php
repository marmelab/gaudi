<?php
$dbInfos = parse_url(getenv('DB_PORT'));

$dsn = 'mysql:dbname=my_db;host='.$dbInfos['host'];

$connection = new PDO($dsn, 'root', '');

$query = 'SELECT * FROM user';
$results = $connection->query($query);
?>
<html>
<head>
    <title>It's working !</title>
</head>
<body>
    <h1>Users</h1>
    <table>
        <?php foreach($results as $result): ?>
            <tr>
                <td><?php echo $result['username']; ?></td>
            </tr>
        <?php endforeach; ?>
    </table>
</body>
</html>
