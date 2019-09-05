# Netapp API Exporter
Prometheus exporter for Netapp ONTAP api. It fetches data from Netapp's filer and exports  them under the prometheus metric `netapp_capacity_svm` at port 9108. Labels of `netapp_capacity_svm` include
* Project_id
* Filer
* Vserver
* Volume
* Metric 
<!-- * ("total", "available", "used" or "percentage") -->

The value of Metric label may be any one of *total*, *available*, *used*, *percentage*.


## Use

### Build
```
go build
```

### Run
Provide list of netapp filers in configuration file "netapp_filers.yaml" and run
```
./netapp-api-exporter [-c netapp_filer_config_file] [-w wait_time] [-l listen_address]
```

### Flags
```
      --help              Show context-sensitive help (also try --help-long and --help-man).
  -w, --wait=300          Wait time
  -c, --config="./netapp_filers.yaml"  
                          Config file
  -l, --listen="0.0.0.0"  Listen address
  -d, --debug             Debug mode
```

### Configuration example
By default, the configuration file is "netapp_filers.yaml". It should contain blocks of following format,
```
- name: xxxx
  host: netapp-bb98.labx.company
  username: <username>
  password: <password>
```

