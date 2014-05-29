<?php

namespace Sandbox\TestBundle\Controller;

use Symfony\Bundle\FrameworkBundle\Controller\Controller;

class TestController extends Controller
{
    public function indexAction()
    {
        return $this->render('SandboxTestBundle:Test:index.html.twig', array('title'=>'Normal Symfony Route'));
    }

}
