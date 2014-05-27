<?php
$container->setParameter('phpcr_backend', [
    "type" => "jackrabbit",
    "url" => 'http://'.getenv('JACKRABBIT_PORT_8082_TCP_ADDR').':'.getenv('JACKRABBIT_PORT_8082_TCP_PORT')."/server/"
]);
