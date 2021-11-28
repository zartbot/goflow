REF: https://github.com/elastic/beats/tree/master/metricbeat/module/golang
 #sudo cp  /etc/metricbeat/modules.d/golang.yml.disabled /etc/metricbeat/modules.d/golang.yml
 sudo cp ./metricbeat.yml  /etc/metricbeat/modules.d/golang.yml
 sudo metricbeat -e -c  metricbeat.yml
