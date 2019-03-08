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

1. Provide list of filers in `netapp_filers.yaml` 
```
  - name: xxxx
    host: netapp-bb98.labx.mo.sap.corp
    username: <username>
    password: <password>
```

2. Build 
```
  go build
```

3. Run
```
  ./netapp-api-exporter
```


## TODO
1. Extract manila micro version and use the highest version
2. Use Seeder to create os user
3. NO need to fetch manila share per filer
4. Convert value from string to float in exporter instead of in filer struct