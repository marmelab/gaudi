<?php

namespace Sandbox\MainBundle\DataFixtures\PHPCR;

use PHPCR\RepositoryInterface;
use Doctrine\ODM\PHPCR\DocumentManager;
use Doctrine\Common\DataFixtures\FixtureInterface;
use Doctrine\Common\DataFixtures\OrderedFixtureInterface;
use Doctrine\Common\Persistence\ObjectManager;

use PHPCR\Util\NodeHelper;

use Symfony\Component\DependencyInjection\ContainerAware;

use Symfony\Cmf\Bundle\MenuBundle\Doctrine\Phpcr\MenuNode;
use Symfony\Cmf\Bundle\MenuBundle\Doctrine\Phpcr\Menu;

class LoadMenuData extends ContainerAware implements FixtureInterface, OrderedFixtureInterface
{
    public function getOrder()
    {
        return 60;
    }

    /**
     * @param \Doctrine\ODM\PHPCR\DocumentManager $dm
     */
    public function load(ObjectManager $dm)
    {
        $session = $dm->getPhpcrSession();

        $basepath = $this->container->getParameter('cmf_menu.persistence.phpcr.menu_basepath');
        $content_path = $this->container->getParameter('cmf_content.persistence.phpcr.content_basepath');

        NodeHelper::createPath($session, $basepath);
        $root = $dm->find(null, $basepath);

        $labels = array('en' => 'Home', 'de' => 'Start', 'fr' => 'Accueil');
        /** @var $main Menu */
        $main = $this->createMenuNode($dm, $root, 'main', $labels, $dm->find(null, "$content_path/home"));
        $main->setChildrenAttributes(array("class" => "menu_main"));

        if ($session->getRepository()->getDescriptor(RepositoryInterface::QUERY_FULL_TEXT_SEARCH_SUPPORTED)) {
            $this->createMenuNode($dm, $main, 'search-item', 'Search', null, null, 'liip_search');
        }

        $this->createMenuNode($dm, $main, 'admin-item', 'Admin', null, null, 'sonata_admin_dashboard');

        $this->createMenuNode($dm, $main, 'blog-item', 'Blog', $dm->find(null, "$content_path/CMF Blog"));

        $projects = $this->createMenuNode($dm, $main, 'projects-item', array('en' => 'Projects', 'de' => 'Projekte', 'fr' => 'Projets'), $dm->find(null, "$content_path/projects"));
        $this->createMenuNode($dm, $projects, 'cmf-item', 'Symfony CMF', $dm->find(null, "$content_path/cmf"));

        $company = $this->createMenuNode($dm, $main, 'company-item', array('en' => 'Company', 'de' => 'Firma', 'fr' => 'Entreprise'), $dm->find(null, "$content_path/company"));
        $this->createMenuNode($dm, $company, 'team-item', array('en' => 'Team', 'de' => 'Team', 'fr' => 'Equipe'), $dm->find(null, "$content_path/team"));
        $this->createMenuNode($dm, $company, 'more-item', array('en' => 'More', 'de' => 'Mehr', 'fr' => 'Plus'), $dm->find(null, "$content_path/more"));

        $demo = $this->createMenuNode($dm, $main, 'demo-item', 'Demo', $dm->find(null, "$content_path/demo"));
        //TODO: this should be possible without a content as the controller might not need a content. support directly having the route document as "content" in the menu document?
        $this->createMenuNode($dm, $demo, 'controller-item', 'Explicit controller', $dm->find(null, "$content_path/demo_controller"));
        $this->createMenuNode($dm, $demo, 'template-item', 'Explicit template', $dm->find(null, "$content_path/demo_template"));
        $this->createMenuNode($dm, $demo, 'type-item', 'Route type to controller', $dm->find(null, "$content_path/demo_type"));
        $this->createMenuNode($dm, $demo, 'class-item', 'Class to controller', $dm->find(null, "$content_path/demo_class"));
        $this->createMenuNode($dm, $demo, 'test-item', 'Normal Symfony Route', null, null, 'test');
        $this->createMenuNode($dm, $demo, 'external-item', 'External Link', null, 'http://cmf.symfony.com/');

        $singlelocale = $this->createMenuNode($dm, $main, 'singlelocale-item', array('en' => 'singlelocale'), $dm->find(null, "$content_path/singlelocale"));
        $this->createMenuNode($dm, $singlelocale, 'singlelocale-sub-item', array('en' => 'singlelocale child'), $dm->find(null, "$content_path/singlelocale"));

        $dm->flush();
    }

    /**
     * @return MenuNode a Navigation instance with the specified information
     */
    protected function createMenuNode(DocumentManager $dm, $parent, $name, $label, $content, $uri = null, $route = null)
    {
        if (!$parent instanceof MenuNode && !$parent instanceof Menu) {
            $menuNode = new Menu();
        } else {
            $menuNode = new MenuNode();
        }

        $menuNode->setParent($parent);
        $menuNode->setName($name);

        $dm->persist($menuNode); // do persist before binding translation

        if (null !== $content) {
            $menuNode->setContent($content);
        } else if (null !== $uri) {
            $menuNode->setUri($uri);
        } else if (null !== $route) {
            $menuNode->setRoute($route);
        }

        if (is_array($label)) {
            foreach ($label as $locale => $l) {
                $menuNode->setLabel($l);
                $dm->bindTranslation($menuNode, $locale);
            }
        } else {
            $menuNode->setLabel($label);
        }

        return $menuNode;
    }
}
