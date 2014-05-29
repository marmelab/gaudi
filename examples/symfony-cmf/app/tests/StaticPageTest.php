<?php

namespace Sandbox;

class StaticPageTest extends WebTestCase
{
    /**
     * @dataProvider contentDataProvider
     */
    public function testContent($url, $title)
    {
        $client = $this->createClient();

        $crawler = $client->request('GET', $url);

        $this->assertEquals(200, $client->getResponse()->getStatusCode());

        $this->assertCount(1, $crawler->filter(sprintf('h1:contains("%s")', $title)));
    }

    public function contentDataProvider()
    {
        return array(
            array('/en/projects', 'The projects'),
            array('/en/projects/cmf', 'Content Management Framework'),
            array('/en/company', 'The Company'),
            array('/en/company/team', 'The Team'),
            array('/en/company/more', 'More Information'),
            array('/demo', 'Routing demo'),
            array('/demo/controller', 'Explicit Controller'),
            array('/demo/atemplate', 'Explicit template'),
            array('/demo/type', 'Controller by type'),
            array('/demo/class', 'Controller by class'),
            array('/hello', 'Hello World!'),
            array('/about', 'Some information about us'),
        );
    }
}
