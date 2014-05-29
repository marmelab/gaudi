<?php

namespace Sandbox\MainBundle\Controller;

use Symfony\Component\DependencyInjection\ContainerInterface;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Cmf\Bundle\ContentBundle\Controller\ContentController as BaseContentController;

/**
 * Special routes to demo the features of the Doctrine Router in the CmfRoutingBundle
 */
class ContentController extends BaseContentController
{
    /**
     * Action that is mapped in the controller_by_type map
     *
     * This can inject something else for the template for content with this type
     *
     * @param object $contentDocument
     *
     * @return \Symfony\Component\HttpFoundation\Response the response
     */
    public function typeAction($contentDocument)
    {
        if (!$contentDocument) {
            throw new NotFoundHttpException('Content not found');
        }

        $params = array(
            'cmfMainContent' => $contentDocument,
            'example' => 'Additional value injected by the controller for this type (this could work without content if we want)',
        );

        return $this->renderResponse('SandboxMainBundle:Demo:controller.html.twig', $params);
    }

    /**
     * Action that is mapped in the controller_by_class map
     *
     * This can inject something else for the template for this type of content.
     *
     * @param object $contentDocument
     *
     * @return \Symfony\Component\HttpFoundation\Response the response
     */
    public function classAction($contentDocument)
    {
        if (!$contentDocument) {
            throw new NotFoundHttpException('Content not found');
        }

        $params = array(
            'cmfMainContent' => $contentDocument,
            'example' => 'Additional value injected by the controller for all content mapped to classAction',
        );

        return $this->renderResponse('SandboxMainBundle:Demo:controller.html.twig', $params);
    }

    /**
     * Action that is explicitly referenced in the _controller field of a content
     *
     * This can inject something else for the template
     *
     * @param object $contentDocument
     *
     * @return \Symfony\Component\HttpFoundation\Response the response
     */
    public function specialAction($contentDocument)
    {
        if (!$contentDocument) {
            throw new NotFoundHttpException('Content not found');
        }

        $params = array(
            'cmfMainContent' => $contentDocument,
            'example' => 'Additional value injected by the controller when explicitly referenced',
        );

        return $this->renderResponse('SandboxMainBundle:Demo:controller.html.twig', $params);
    }
}
