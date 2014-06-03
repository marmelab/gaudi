[[ $dir := .Container.GetFirstMountedDir ]]
[[ $projectName := .Container.GetCustomValue "project_name" "project" ]]
[[ $version := (.Container.GetCustomValue "version" "1.9.3") ]]

if [ ! -d "[[$dir]]/[[ $projectName ]]" ]; then

    # Install ror
    /bin/bash -l -c 'gem install rails bundler'
    cd [[ $dir ]] && /bin/bash -l -c 'rails new [[ $projectName ]] [[ if(.Container.DependsOf "mysql" )]] -d mysql [[ end ]] -T'

    cd [[ $dir ]]/[[ $projectName ]]
    printf "gem 'execjs'\ngem 'therubyracer'" >> Gemfile
    /bin/bash -l -c 'bundle install'

    [[ if(.Container.DependsOf "mysql") ]]
        [[ $db := (.Collection.GetType "mysql")]]
        sed -i -e "s/host: localhost/host: <%= ENV['DB_PORT_3306_TCP_ADDR'] %>/" [[ $dir ]]/[[ $projectName ]]/config/database.yml
        sed -i -e "s/host: localhost/host: <%= ENV['DB_PORT_3306_TCP_ADDR'] %>/" [[ $dir ]]/[[ $projectName ]]/config/database.yml
        /bin/bash -c -l 'rake db:migrate'
    [[ end ]]
fi
