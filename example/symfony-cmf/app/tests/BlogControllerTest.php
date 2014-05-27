<?php

namespace Sandbox;

class BlogControllerTest extends WebTestCase
{
    public function testList()
    {
        $client = $this->createClientAuthenticated();

        $crawler = $client->request('GET', '/blog/cmf-blog');

        $response = $client->getResponse();

        $this->assertEquals(200, $response->getStatusCode());
        $this->assertGreaterThan(1, $crawler->filter('div.post')->count());
    }
}

