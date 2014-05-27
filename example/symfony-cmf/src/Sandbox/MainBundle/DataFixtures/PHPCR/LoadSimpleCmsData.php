<?php

namespace Sandbox\MainBundle\DataFixtures\PHPCR;

use Doctrine\Common\DataFixtures\FixtureInterface;
use Doctrine\Common\DataFixtures\OrderedFixtureInterface;
use Doctrine\Common\Persistence\ObjectManager;

use PHPCR\Util\NodeHelper;

use Symfony\Component\DependencyInjection\ContainerAware;

use Symfony\Cmf\Bundle\SimpleCmsBundle\Doctrine\Phpcr\Page;

class LoadSimpleCmsData extends ContainerAware implements FixtureInterface, OrderedFixtureInterface
{
    public function getOrder()
    {
        return 50;
    }

    public function load(ObjectManager $manager)
    {
        $basepath = $this->container->getParameter('cmf_simple_cms.persistence.phpcr.menu_basepath');
        $base = $manager->find(null, $basepath);

        $root = $this->createPage($manager, $base, 'simple', 'root', 'root page of simple menu, never used', '');
        $this->createPage($manager, $root, 'about', 'About us', 'Some information about us', 'The about us page with some content');
        $this->createPage($manager, $root, 'contact', 'Contact', 'A contact page', 'Please send an email to cmf-devs@groups.google.com');

        $manager->flush();
    }

    /**
     * @return Page instance with the specified information
     */
    protected function createPage($manager, $parent, $name, $label, $title, $body)
    {
        $page = new Page();
        $page->setPosition($parent, $name);
        $page->setLabel($label);
        $page->setTitle($title);
        $page->setBody($body);

        $manager->persist($page);

        return $page;
    }
}
