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

### Parameters
* -c, --config sets the configuration file for netapp filers. Default value is `netapp_filers.yaml`. It should in the following format,
```
  - name: xxxx
    host: netapp-bb98.labx.company
    username: <username>
    password: <password>
```
* -w, --wait sets the time in seconds to wait between each query to the netapp filer. Default value is `300`.
* -l, --listen sets the allowed listen address. Default is `0.0.0.0`.


## TODO
1. Extract manila micro version and use the highest version
2. Use Seeder to create os user
3. NO need to fetch manila share per filer
4. Convert value from string to float in exporter instead of in filer struct
