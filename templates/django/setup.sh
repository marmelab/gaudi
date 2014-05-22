if [ ! -d "/app/[[ .Container.GetCustomValue "project_name" "project" ]]/[[ .Container.GetCustomValue "app_name" "myapp" ]]" ]; then
	# Install django & configure it
	cd /app
	django-admin.py startproject [[ .Container.GetCustomValue "project_name" "project" ]] .

	mkdir ./[[ .Container.GetCustomValue "project_name" "project" ]]/[[ .Container.GetCustomValue "app_name" "myapp" ]]
	python ./manage.py startapp [[ .Container.GetCustomValue "app_name" "myapp" ]] ./[[ .Container.GetCustomValue "project_name" "project" ]]/[[ .Container.GetCustomValue "app_name" "myapp" ]]

	cd /app/[[ .Container.GetCustomValue "project_name" "project" ]]
	sed -i -e "s/'django.db.backends.sqlite3'/'django.db.backends.mysql'/" ./settings.py
	sed -i -e "s/'NAME': os.path.join(BASE_DIR, 'db.sqlite3'),/'NAME': 'django',\n\t\t'USER': 'root',\n\t\t'PASSWORD': '',\n\t\t'HOST': os.environ['DB_PORT_3306_TCP_ADDR']/" ./settings.py

	sed -i -e "s/# from django.contrib import admin/from django.contrib import admin/" ./urls.py
	sed -i -e "s/# admin.autodiscover()/admin.autodiscover()/" ./urls.py
	sed -i -e "s/# url(r'^admin\/', include(admin.site.urls))/url(r'^admin\/', include(admin.site.urls))/" ./urls.py

	echo -e "import os, sys\nbase = os.path.dirname(os.path.dirname(__file__))\nbase_parent = os.path.dirname(base)\nsys.path.append(base)\nsys.path.append(base_parent)\n\n$(cat /app/[[ .Container.GetCustomValue "project_name" "project" ]]/wsgi.py)" > /app/[[ .Container.GetCustomValue "project_name" "project" ]]/wsgi.py
fi
