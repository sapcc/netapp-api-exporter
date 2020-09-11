# Netapp API Exporter
Prometheus exporter for Netapp ONTAP API. It fetches data from Netapp's filer and exports them as prometheus metrics. There are mainly two groups of metrics that have been implemented.

__Volume Metrics__ with labels `availability_zone`, `filer`, `project_id`, `share_id`, `volume` and `vserver`. <sup>1</sup>

* netapp_volume_state <sup>2</sup>
* netapp_volume_total_bytes
* netapp_volume_used_bytes
* netapp_volume_available_bytes
* netapp_volume_snapshot_used_bytes
* netapp_volume_snapshot_reserved_bytes
* netapp_volume_snapshot_available_bytes
* netapp_volume_used_percentage
* netapp_volume_saved_total_percentage
* netapp_volume_saved_compression_percentage
* netapp_volume_saved_deduplication_percentage

<sup>1</sup> The label `project_id` is openstack specific, and `share_id` is openstack manila specific.

<sup>2</sup> The metric netapp_volume_state being 1 means "online"; being -1 means "offline".

__Aggregate Metrics__ with labels `availability_zone`, `filer`, `node` and `aggregate`.
* netapp_aggregate_total_bytes
* netapp_aggregate_used_bytes
* netapp_aggregate_available_bytes
* netapp_aggregate_used_percentage
* netapp_aggregate_physical_used_bytes
* netapp_aggregate_physical_percentage

In addition, filer status metrics (labes `availability_zone`, `filer`).
* netapp_filer_scrape_failure

## Usage

This collector includes three groups of metrics: volume metrics, aggregate metrics and system info metrics. Each group
can be disabled with the --no-\* flag.

### Flags
```
Flags:
      --help              Show context-sensitive help (also try --help-long and --help-man).
  -c, --config="./netapp_filers.yaml"  
                          Config file
  -l, --listen="0.0.0.0"  Listen address
  -d, --debug             Debug mode
      --no-aggregate           Disable aggregate collector
      --no-volume              Disable volume collector
      --no-system              Disable system collector
      --aggregateRetention=5m  Aggregate collector retention period
      --volumeRetention=2m     Volume collector retention period
```

### Configuration 
Configuration file is in yaml format (default path "./config/netapp_filers.yaml"). It should contain blocks in following format,
```
- name: xxxx
  host: netapp-bb98.labx.company
  username: <username>
  password: <password>
```

