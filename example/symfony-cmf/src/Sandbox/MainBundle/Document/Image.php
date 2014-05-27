<?php
namespace Sandbox\MainBundle\Document;

use Doctrine\ODM\PHPCR\Mapping\Annotations as PHPCRODM;
use Doctrine\Common\Collections\ArrayCollection;

use Symfony\Component\Validator\Constraints as Assert;

/**
 * @PHPCRODM\Document(referenceable=true)
 */
class Image
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
     * @Assert\NotBlank
     * @Assert\Regex("{^[a-z]+$}")
     * @PHPCRODM\String()
     */
    public $name;

    /**
     * @Assert\NotBlank()
     * @var stream
     * @PHPCRODM\Binary
     */
    public $content;

    /**
     * @PHPCRODM\String(multivalue=true)
     */
    public $tags;

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
}
