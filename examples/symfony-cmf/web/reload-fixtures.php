<pre>
<?php

$output = $return_var= null;
$command = __DIR__.'/../app/console -v doctrine:phpcr:fixtures:load -e=prod';
echo "Running: $command\n";
exec($command, $output, $return_var);

if (!empty($output) && is_array($output)) {
    echo "Output:\n";
    foreach ($output as $line) {
        echo $line."\n";
    }
} else {
    echo 'Fixtures could not be loaded: '.var_export($return_var, true);
}
