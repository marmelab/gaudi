<?php

namespace Sandbox\MainBundle\DataFixtures\PHPCR;

use Doctrine\Common\DataFixtures\FixtureInterface;
use Symfony\Cmf\Bundle\BlogBundle\Util\PostUtils;
use Doctrine\Common\DataFixtures\OrderedFixtureInterface;
use Doctrine\Common\Persistence\ObjectManager;

use PHPCR\Util\NodeHelper;

use Symfony\Cmf\Bundle\RoutingBundle\Doctrine\Phpcr\RedirectRoute;
use Symfony\Cmf\Bundle\RoutingBundle\Doctrine\Phpcr\Route;
use Symfony\Component\DependencyInjection\ContainerAware;

use Symfony\Cmf\Component\Routing\RouteObjectInterface;

class LoadRoutingData extends ContainerAware implements FixtureInterface, OrderedFixtureInterface
{
    public function getOrder()
    {
        return 21;
    }

    /**
     * Load routing data into the document manager.
     *
     * NOTE: We demo all possibilities. Of course, you should try to be
     * consistent in what you use and only use different things for special
     * cases.
     *
     * @param $dm \Doctrine\ODM\PHPCR\DocumentManager
     */
    public function load(ObjectManager $dm)
    {
        $session = $dm->getPhpcrSession();

        $basepath = $this->container->getParameter('cmf_routing.dynamic.persistence.phpcr.route_basepath');
        if ($session->itemExists($basepath)) {
            $session->removeItem($basepath);
        }

        NodeHelper::createPath($session, $basepath);
        $parent = $dm->find(null, $basepath);

        $content_path = $this->container->getParameter('cmf_content.persistence.phpcr.content_basepath');
        $locales = $this->container->getParameter('locales');
        foreach ($locales as $locale) {
            $home = new Route();
            $home->setPosition($parent, $locale);
            $home->setDefault(RouteObjectInterface::TEMPLATE_NAME, 'SandboxMainBundle:Homepage:index.html.twig');
            $home->setContent($dm->find(null, "$content_path/home"));
            $dm->persist($home);

            $blog = new Route;
            $blog->setPosition($home, PostUtils::slugify('CMF Blog'));
            $blog->setContent($dm->find(null, "$content_path/CMF Blog"));
            $dm->persist($blog);

            $company = new Route;
            $company->setPosition($home, 'company');
            $company->setContent($dm->find(null, "$content_path/company"));
            $dm->persist($company);

            $team = new Route;
            $team->setPosition($company, 'team');
            $team->setContent($dm->find(null, "$content_path/team"));
            $dm->persist($team);

            $more = new Route;
            $more->setPosition($company, 'more');
            $more->setContent($dm->find(null, "$content_path/more"));
            $dm->persist($more);

            $projects = new Route;
            $projects->setPosition($home, 'projects');
            $projects->setContent($dm->find(null, "$content_path/projects"));
            $dm->persist($projects);

            $cmf = new Route;
            $cmf->setPosition($projects, 'cmf');
            $cmf->setContent($dm->find(null, "$content_path/cmf"));
            $dm->persist($cmf);
        }

        // demo features of routing

        // we can create routes without locales, but will lose the language context of course

        $demo = new Route;
        $demo->setPosition($parent, 'demo');
        $demo->setContent($dm->find(null, "$content_path/demo"));
        $demo->setDefault(RouteObjectInterface::TEMPLATE_NAME, 'SandboxMainBundle:Demo:template_explicit.html.twig');
        $dm->persist($demo);

        // explicit template
        $template = new Route;
        $template->setPosition($demo, 'atemplate');
        $template->setContent($dm->find(null, "$content_path/demo_template"));
        $template->setDefault(RouteObjectInterface::TEMPLATE_NAME, 'SandboxMainBundle:Demo:template_explicit.html.twig');
        $dm->persist($template);

        // explicit controller
        $controller = new Route;
        $controller->setPosition($demo, 'controller');
        $controller->setContent($dm->find(null, "$content_path/demo_controller"));
        $controller->setDefault('_controller', 'sandbox_main.controller:specialAction');
        $dm->persist($controller);

        // type to controller mapping
        $type = new Route;
        $type->setPosition($demo, 'type');
        $type->setContent($dm->find(null, "$content_path/demo_type"));
        $type->setDefault('type', 'demo_type');
        $dm->persist($type);

        // class to controller mapping
        $class = new Route;
        $class->setPosition($demo, 'class');
        $class->setContent($dm->find(null, "$content_path/demo_class"));
        $dm->persist($class);

        // redirections

        // redirect to uri
        $redirect = new RedirectRoute();
        $redirect->setPosition($parent, 'external');
        $redirect->setUri('http://cmf.symfony.com');
        $dm->persist($redirect);

        // redirect to other doctrine route
        $redirectRoute = new RedirectRoute();
        $redirectRoute->setPosition($parent, 'short');
        $redirectRoute->setRouteTarget($cmf);
        $dm->persist($redirectRoute);

        // redirect to Symfony route
        $redirectS = new RedirectRoute();
        $redirectS->setPosition($parent, 'short1');
        $redirectS->setRouteName('test');
        $dm->persist($redirectS);

        // class to template mapping is used for all the rest

        $default_locale = $this->container->getParameter('locale');
        $singlelocale = new Route;
        $singlelocale->setPosition($dm->find(null, "$basepath/$default_locale"), 'singlelocale');
        $singlelocale->setDefault('_locale', $default_locale);
        $singlelocale->setRequirement('_locale', $default_locale);
        $singlelocale->setContent($dm->find(null, "$content_path/singlelocale"));
        $dm->persist($singlelocale);

        $dm->flush();
    }
}
