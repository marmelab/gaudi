mkdir htdocs
wget http://www.magentocommerce.com/downloads/assets/1.9.0.1/magento-1.9.0.1.tar.gz
tar zxvf magento-1.9.0.1.tar.gz -C htdocs --strip 1


echo "<?php echo 'ok';" > htdocs/up.php
