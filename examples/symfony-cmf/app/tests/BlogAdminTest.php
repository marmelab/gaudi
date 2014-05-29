<?php

namespace Sandbox;

class BlogAdminTest extends WebTestCase
{
    public function testList()
    {
        $client = $this->createClientAuthenticated();

        $client->request('GET', '/en/admin/cmf/blog/blog/list');

        $response = $client->getResponse();
        $this->assertEquals(200, $response->getStatusCode());
        $this->assertContains('Blogs', $response->getContent());
    }
}

