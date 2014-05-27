<?php

namespace Sandbox;

class HomepageTest extends WebTestCase
{
    public function testRedirectToDefaultLanguage()
    {
        $client = $this->createClient();

        $client->request('GET', '/');

        $this->assertEquals(301, $client->getResponse()->getStatusCode());

        $client->followRedirect();

        $this->assertContains('http://localhost/en', $client->getRequest()->getUri());
    }

    public function testContents()
    {
        $client = $this->createClient();

        $crawler = $client->request('GET', '/en');

        $this->assertEquals(200, $client->getResponse()->getStatusCode());

        $this->assertCount(3, $crawler->filter('.cmf-block'));
        $this->assertCount(1, $crawler->filter('h1:contains(Homepage)'));
        $this->assertCount(1, $crawler->filter('h2:contains("Welcome to the Symfony CMF Demo")'));

        $menuCount = $this->isSearchSupported() ? 17 : 16;
        $this->assertCount($menuCount, $crawler->filter('ul.menu_main li'));
    }

    public function testJsonContents()
    {
        $client = $this->createClient();

        $client->request(
            'GET',
            '/en',
            array(),
            array(),
            array(
                'HTTP_ACCEPT'  => 'application/json',
                'CONTENT_TYPE' => 'application/json'
            )
        );
        $this->assertEquals(200, $client->getResponse()->getStatusCode());

        $json = @json_decode($client->getResponse()->getContent());
        $this->assertNotEmpty($json);
    }
}
