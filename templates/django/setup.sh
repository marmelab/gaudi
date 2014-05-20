if [ ! -d "/app/project/myapp" ]; then
	# Install django & configure it
	cd /app

	django-admin.py startproject project
	python ./project/manage.py startapp myapp

	cd /app/project

	sed -i "1s/^/import os\n/" ./settings.py
	sed -i -e "s/'ENGINE': 'django.db.backends.'/'ENGINE': 'django.db.backends.mysql'/" ./settings.py
	sed -i -e "s/'NAME': ''/'NAME': 'django'/" ./settings.py
	sed -i -e "s/'USER': ''/'User': 'root'/" ./settings.py
	sed -i -e "s/'HOST': ''/'HOST': os.environ['DB_PORT_3306_TCP_ADDR']/" ./settings.py
	sed -i -e "s/# 'django.contrib.admin'/'django.contrib.admin'/" ./settings.py
	sed -i -e "s/'django.contrib.sites'/# 'django.contrib.sites'/" ./settings.py
	sed -i -e "s/'django.contrib.admindocs'/'django.contrib.admindocs',\n\t'myapp'/" ./settings.py

	sed -i -e "s/# from django.contrib import admin/from django.contrib import admin/" ./urls.py
	sed -i -e "s/# admin.autodiscover()/admin.autodiscover()/" ./urls.py
	sed -i -e "s/# url(r'^admin\/', include(admin.site.urls))/url(r'^admin\/', include(admin.site.urls))/" ./urls.py
fi
