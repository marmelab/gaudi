<?php

namespace Sandbox\MainBundle\Document;

use Doctrine\ODM\PHPCR\Mapping\Annotations as PHPCRODM;
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Cmf\Component\Routing\RouteReferrersReadInterface;

/**
 * A document that we map to a controller
 *
 * @PHPCRODM\Document(referenceable=true)
 */
class DemoClassContent implements RouteReferrersReadInterface
{
    /**
     * to create the document at the specified location. read only for existing documents.
     *
     * @PHPCRODM\Id
     */
    protected $path;

    /**
     * @PHPCRODM\Node
     */
    public $node;

    /**
     * @PHPCRODM\Parentdocument()
     */
    public $parent;

    /**
     * @Assert\NotBlank
     * @PHPCRODM\Nodename()
     */
    protected $name;

    /**
     * @Assert\NotBlank
     * @PHPCRODM\String()
     */
    protected $title;

    /**
     * @Assert\NotBlank
     * @PHPCRODM\String()
     */
    protected $body;

    /**
     * @PHPCRODM\Referrers(referringDocument="Symfony\Cmf\Bundle\RoutingBundle\Doctrine\Phpcr\Route", referencedBy="content")
     */
    public $routes;

    public function getName()
    {
        return $this->name;
    }

    public function setName($name)
    {
        $this->name = $name;
    }

    public function getTitle()
    {
        return $this->title;
    }

    public function setTitle($title)
    {
        $this->title = $title;
    }

    /**
     * Set repository path of this navigation item for creation
     */
    public function setPath($path)
    {
      $this->path = $path;
    }
    public function getPath()
    {
      return $this->path;
    }
    public function setParent($parent)
    {
        $this->parent = $parent;
    }

    public function getBody()
    {
        return $this->body;
    }

    public function setBody($content)
    {
        $this->body = $content;
    }

    /**
     * @return array of route objects that point to this content
     */
    public function getRoutes()
    {
        return $this->routes->toArray();
    }
}
