<?php

namespace Sandbox;

use Symfony\Component\Routing\Exception\MissingMandatoryParametersException;

class AdminTest extends WebTestCase
{
    protected $pool;
    protected $router;

    protected $verifiablePatterns = array(
        '/cmf/content/staticcontent/list',
        '/cmf/content/staticcontent/create',
        '/cmf/content/staticcontent/{id}/edit',
        '/cmf/content/staticcontent/{id}/delete',
        '/cmf/block/simpleblock/list',
        '/cmf/block/simpleblock/create',
        '/cmf/block/simpleblock/{id}/edit',
        '/cmf/block/simpleblock/{id}/delete',
        '/cmf/block/containerblock/list',
        '/cmf/block/containerblock/create',
        '/cmf/block/containerblock/{id}/edit',
        '/cmf/block/containerblock/{id}/delete',
        '/cmf/block/referenceblock/list',
        '/cmf/block/referenceblock/create',
        '/cmf/block/referenceblock/{id}/edit',
        '/cmf/block/referenceblock/{id}/delete',
        '/cmf/block/actionblock/list',
        '/cmf/block/actionblock/create',
        '/cmf/block/actionblock/{id}/edit',
        '/cmf/block/actionblock/{id}/delete',
        '/cmf/block/imagineblock/list',
        '/cmf/block/imagineblock/create',
        '/cmf/routing/route/list',
        '/cmf/routing/route/create',
        '/cmf/routing/route/{id}/edit',
        '/cmf/routing/route/{id}/delete',
        '/cmf/routing/redirectroute/list',
        '/cmf/routing/redirectroute/create',
        '/cmf/routing/redirectroute/{id}/edit',
        '/cmf/routing/redirectroute/{id}/delete',
        '/cmf/menu/menu/list',
        '/cmf/menu/menu/create',
        '/cmf/menu/menu/{id}/edit',
        '/cmf/menu/menu/{id}/delete',
        '/cmf/menu/menunode/list',
        '/cmf/menu/menunode/create',
        '/cmf/menu/menunode/{id}/edit',
        '/cmf/menu/menunode/{id}/delete',
        '/cmf/blog/blog/list',
        '/cmf/blog/blog/create',
        '/cmf/blog/blog/{id}/edit',
        '/cmf/blog/blog/{id}/delete',
        '/cmf/blog/post/list',
        '/cmf/blog/post/create',
        '/cmf/blog/post/{id}/edit',
        '/cmf/blog/post/{id}/delete',
        '/cmf/simplecms/page/list',
        '/cmf/simplecms/page/create',
        '/cmf/simplecms/page/{id}/edit',
        '/cmf/simplecms/page/{id}/delete',
    );

    protected $testedPatterns = array();

    public function setUp()
    {
        parent::setUp();
        $this->pool = $this->getContainer()->get('sonata.admin.pool');
        $this->router = $this->getContainer()->get('router');
        $this->client = $this->createClientAuthenticated();
        $this->dm = $this->getContainer()->get('doctrine_phpcr.odm.default_document_manager');
    }

    public function testAdmin()
    {
        $adminGroups = $this->pool->getAdminGroups();
        $admins = array();

        foreach (array_keys($adminGroups) as $adminName) {
            $admins = array_merge($admins, $this->pool->getAdminsByGroup($adminName));
        }

        foreach ($admins as $admin) {
            $this->doTestReachableAdminRoutes($admin);
        }

        // verify that we have tested everything we wanted to test.
        $this->assertEquals($this->verifiablePatterns, $this->testedPatterns);
    }

    protected function doTestReachableAdminRoutes($admin)
    {
        $routeCollection = $admin->getRoutes();
        $class = $admin->getClass();
        $routeParams = array('_locale' => 'en');

        foreach ($routeCollection->getElements() as $route) {
            $requirements = $route->getRequirements();

            // fix this one later
            if (strpos($route->getPattern(), 'export')) {
                continue;
            }

            // these don't all work atm
            if (strpos($route->getPattern(), 'show')) {
                continue;
            }

            // do not test POST routes
            if (isset($requirements['_method'])) {
                if ($requirements['_method'] != 'GET') {
                    continue;
                }
            }

            // if an ID is required, try and find a document to test
            if (isset($requirements['id'])) {
                if ($document = $this->dm->getRepository($class)->findOneBy(array())) {
                    $node = $this->dm->getNodeForDocument($document);
                    $routeParams['id'] = $node->getPath();
                } else {
                    // we should throw an exception here maybe and fix the missing fixtures
                }
            }

            try {
                $url = $this->router->generate($route->getDefault('_sonata_name'), $routeParams);
            } catch (MissingMandatoryParametersException $e) {
                // do not try and load pages with parameters, e.g. edit, show, etc.
                continue;
            }

            $crawler = $this->client->request('GET', $url);
            $res = $this->client->getResponse();
            $statusCode = $res->getStatusCode();

            $this->assertEquals(200, $statusCode, sprintf(
                'URL %s returns 200 OK HTTP Code', $url
            ));

            $this->testedPatterns[] = $route->getPattern();
        }
    }
}
