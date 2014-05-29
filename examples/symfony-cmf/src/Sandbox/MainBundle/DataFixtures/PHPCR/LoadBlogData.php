<?php

namespace Sandbox\MainBundle\DataFixtures\PHPCR;

use Doctrine\Common\DataFixtures\FixtureInterface;
use PHPCR\RepositoryInterface;
use Doctrine\Common\DataFixtures\OrderedFixtureInterface;
use Doctrine\Common\Persistence\ObjectManager;

use PHPCR\Util\NodeHelper;

use Symfony\Component\DependencyInjection\ContainerAware;

use Symfony\Cmf\Bundle\BlogBundle\Document\Blog;
use Symfony\Cmf\Bundle\BlogBundle\Document\Post;

class LoadBlogData extends ContainerAware implements FixtureInterface, OrderedFixtureInterface
{
    public function getOrder()
    {
        return 30;
    }

    /**
     * @param \Doctrine\ODM\PHPCR\DocumentManager $dm
     */
    public function load(ObjectManager $dm)
    {
        $session = $dm->getPhpcrSession();

        $basepath = $this->container->getParameter('cmf_blog.blog_basepath');
        NodeHelper::createPath($session, $basepath);
        $root = $dm->find(null, $basepath);

        // generate some blog data here..
        $blog = new Blog;
        $blog->setParent($root);
        $blog->setName('CMF Blog');
        $dm->persist($blog);

        for ($i = 1; $i <= 20; $i++) {
            $p = new Post;
            $p->setTitle(ucfirst($this->getWords(rand(2,5))));
            $p->setDate(new \DateTime());
            $p->setBody($this->getWords());
            $p->setBlog($blog);
            $p->setPublishable(true);
            $dm->persist($p);
        }

        $dm->flush();
    }

    protected function getWords($nbWords = null)
    {
        $text = <<<HERE
This tutorial shows how to install the Symfony CMF Sandbox, a demo platform aimed at showing the tool's basic features running on a demo environment. This can be used to evaluate the platform or to see actual code in action, helping you understand the tool's internals.
While it can be used as such, this sandbox does not intend to be a development platform. If you are looking for installation instructions for a development setup, please refer to:

As Symfony CMF Sandbox is based on Symfony2, you should make sure you meet the Requirements for running Symfony2. Git 1.6+, Curl and PHP Intl are also needed to follow the installation steps listed below.
If you wish to use Jackalope + Apache JackRabbit as the storage medium recommended, you will also need Java JRE. For other mechanisms and its requirements, please refer to their respective sections.

ontent Repository API for Java JCR is a specification for a Java platform application programming interface API to access content repositories in a uniform manner.1dead link2not in citation given The content repositories are used in content management systems to keep the content data and also the metadata used in content management systems CMS such as versioning metadata. The specification was developed under the Java Community Process as JSR-170 Version 1.34 and as JSR-283 version 25 The main Java package is javax.jcr.

he PHP Content Repository is an adaption of the Java Content Repository JCR standard, an open API specification defined in JSR-283. 
The API defines how to handle hierarchical semi-structured data in a consistent way. The typical use case is content management systems. PHPCR combines the best out of document-oriented databases weak structured data and of XML databases hierarchical trees. On top of that, it adds useful features like searching, versioning, access control and locking on top of it.
HERE;

        $words = explode(' ', $text);

        $nbWords = $nbWords ? $nbWords : rand(100,500);
        $sWords = array();
        $paraBreak = rand(100,500);
        $paraIdx = 0;
        for ($i = 0; $i < $nbWords; $i ++) {
            $sWords[] = $words[rand(0, (count($words) - 1))];
            if ($paraIdx++ == $paraBreak) {
                $sWords[] = "\n\n";
                $paraIdx = 0;
            };
        }

        return implode(' ', $sWords);
    }
}
