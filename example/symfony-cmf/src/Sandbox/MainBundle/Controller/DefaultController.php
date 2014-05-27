<?php

namespace Sandbox\MainBundle\Controller;

use Symfony\Bundle\FrameworkBundle\Controller\Controller,
    Symfony\Component\HttpFoundation\Response,
    Sonata\BlockBundle\Model\BlockInterface;

class DefaultController extends Controller
{
    /**
     * Action that is referenced in an ActionBlock
     *
     * @param \Sonata\BlockBundle\Model\BlockInterface $block
     *
     * @return \Symfony\Component\HttpFoundation\Response the response
     */
    public function blockAction($block)
    {
        return $this->render('SandboxMainBundle:Block:demo_action_block.html.twig', array(
            'block' => $block
        ));
    }
}
