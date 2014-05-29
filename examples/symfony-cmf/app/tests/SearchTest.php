<?php

namespace Sandbox;

class SearchTest extends WebTestCase
{
    public function testSearch()
    {
        if (!$this->isSearchSupported()) {
            $this->markTestSkipped('Fulltext search is not supported.');
        }

        $client = $this->createClient();

        $client->request('GET', '/search?query=cmf');

        $this->assertEquals(200, $client->getResponse()->getStatusCode());

        $this->assertContains('results for &quot;cmf&quot; have been found', $client->getResponse()->getContent());
    }
}
