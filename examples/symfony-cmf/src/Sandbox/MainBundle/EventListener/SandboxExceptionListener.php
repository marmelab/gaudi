<?php

namespace Sandbox\MainBundle\EventListener;

use PHPCR\RepositoryException;

use Symfony\Component\DependencyInjection\ContainerAware;
use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Symfony\Component\HttpFoundation\Response;
use Symfony\Component\HttpKernel\KernelEvents;
use Symfony\Component\HttpKernel\Event\GetResponseForExceptionEvent;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;

/**
 * Exception listener that will handle not found exceptions and try to give the
 * first time installer some clues what is wrong.
 */
class SandboxExceptionListener extends ContainerAware implements EventSubscriberInterface
{
    public function onKernelException(GetResponseForExceptionEvent $event)
    {
        if (! $event->getException() instanceof NotFoundHttpException) {
            return;
        }

        if (! $this->container->has('doctrine_phpcr.odm.default_document_manager')) {
            $error = 'Missing the service doctrine_phpcr.odm.default_document_manager.';
        } else {
            try {
                $om = $this->container->get('doctrine_phpcr.odm.default_document_manager');
                $doc = $om->find(null, $this->container->getParameter('cmf_menu.persistence.phpcr.menu_basepath'));
                if ($doc) {
                    $error = 'Hm. No clue what goes wrong. Maybe this is a real 404?<pre>'.$event->getException()->__toString().'</pre>';
                } else {
                    $error = 'Did you load the fixtures? See README for how to load them. I found no node at menu_basepath: '.$this->container->getParameter('cmf_menu.persistence.phpcr.menu_basepath');
                }
            } catch(RepositoryException $e) {
                $error = 'There was an exception loading the document manager: <strong>' . $e->getMessage() .
                    "</strong><br/>\n<em>Make sure you have a phpcr backend properly set up and running.</em><br/><pre>".
                    $e->__toString() .'</pre>';
            }
        }
        // do not even trust the templating system to work
        $response = new Response("<html><body>
            <h2>Sandbox</h2>
            <p>If you see this page, it means your sandbox is not correctly set up.
               Please see the README file in the sandbox root folder and if you can't figure out
               what is wrong, ask us on freenode irc #symfony-cmf or the mailinglist cmf-users@groups.google.com.
            </p>

            <p>If you are seeing this page as the result of an edit in the admin tool, please report what you were doing
                to our <a href=\"https://github.com/symfony-cmf/cmf-sandbox/issues/new\">ticket system</a>,
                so that we can add means to prevent this issue in the future. But to get things working again
                for now, please just <a href=\"".$event->getRequest()->getSchemeAndHttpHost()."/reload-fixtures.php\">click here</a>
                to reload the data fixtures.
            </p><p style='color:red;'>
               <strong>Detected the following problem</strong>: $error
            </p>
            </body></html>
            ");

        $event->setResponse($response);
    }


    static public function getSubscribedEvents()
    {
        return array(
            KernelEvents::EXCEPTION => array('onKernelException', 0),
        );
    }
}
