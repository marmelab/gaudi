mkdir htdocs
tar zxvf magento-1.9.0.1.tar.gz -C htdocs --strip 1


echo "<?php echo 'ok';" > /htdocs/up.php
