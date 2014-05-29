<?php

($autoload = @include_once __DIR__ . '/../vendor/autoload.php') || $autoload = @include_once __DIR__ . '/../../../autoload.php';
if (!$autoload) {
    throw new RuntimeException('Install dependencies to run phpcr.');
}

use Doctrine\Common\Annotations\AnnotationRegistry;
AnnotationRegistry::registerLoader(array($autoload, 'loadClass'));
AnnotationRegistry::registerFile(__DIR__.'/../lib/Doctrine/ODM/PHPCR/Mapping/Annotations/DoctrineAnnotations.php');

$configFile = getcwd() . DIRECTORY_SEPARATOR . 'cli-config.php';

$helperSet = null;
if (file_exists($configFile)) {
    if (!is_readable($configFile)) {
        trigger_error(
            'Configuration file [' . $configFile . '] does not have read permission.', E_USER_ERROR
        );
    }

    require $configFile;

    foreach ($GLOBALS as $helperSetCandidate) {
        if ($helperSetCandidate instanceof \Symfony\Component\Console\Helper\HelperSet) {
            $helperSet = $helperSetCandidate;
            break;
        }
    }
} else {
    trigger_error(
        'Configuration file [' . $configFile . '] does not exist. See https://github.com/doctrine/phpcr-odm/wiki/Command-line-tool-configuration', E_USER_ERROR
    );
}

$helperSet = ($helperSet) ?: new \Symfony\Component\Console\Helper\HelperSet();

$cli = new \Symfony\Component\Console\Application('Doctrine ODM PHPCR Command Line Interface', Doctrine\ODM\PHPCR\Version::VERSION);
$cli->setCatchExceptions(true);
$cli->setHelperSet($helperSet);
$cli->addCommands(array(
    new \PHPCR\Util\Console\Command\WorkspaceCreateCommand(),
    new \PHPCR\Util\Console\Command\NodeDumpCommand(),
    new \PHPCR\Util\Console\Command\WorkspaceExportCommand(),
    new \PHPCR\Util\Console\Command\WorkspaceImportCommand(),
    new \PHPCR\Util\Console\Command\WorkspaceListCommand(),
    new \PHPCR\Util\Console\Command\WorkspacePurgeCommand(),
    new \PHPCR\Util\Console\Command\WorkspaceQueryCommand(),
    new \PHPCR\Util\Console\Command\NodeTypeRegisterCommand(),
    new \Doctrine\ODM\PHPCR\Tools\Console\Command\RegisterSystemNodeTypesCommand(),
    new \Doctrine\ODM\PHPCR\Tools\Console\Command\InfoDoctrineCommand(),
    new \Doctrine\ODM\PHPCR\Tools\Console\Command\DumpQueryBuilderReferenceCommand(),
));
if (isset($extraCommands) && ! empty($extraCommands)) {
    $cli->addCommands($extraCommands);
}
$cli->run();
