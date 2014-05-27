#!/bin/bash

# Download installation script from phpcr-midgard2
wget https://raw.github.com/midgardproject/phpcr-midgard2/master/tests/travis_midgard.sh
chmod +x ./travis_midgard.sh
./travis_midgard.sh

# Copy PHPCR schemas to Midgard's global schema dir
sudo wget --directory-prefix=/usr/share/midgard2/schema https://github.com/midgardproject/phpcr-midgard2/raw/master/data/share/schema/midgard_namespace_registry.xml
sudo wget --directory-prefix=/usr/share/midgard2/schema https://github.com/midgardproject/phpcr-midgard2/raw/master/data/share/schema/midgard_tree_node.xml

composer require midgard/phpcr:dev-master --prefer-source
