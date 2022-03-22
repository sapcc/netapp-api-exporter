# Netapp API Exporter

Prometheus exporter for Netapp ONTAP API (not the ONTAP REST API). The package
is tested against ONTAP version 9.2 and up. It fetches data from Netapp's filer
and exports them as prometheus metrics.

## Usage

This collector includes three groups of metrics: volume metrics, aggregate
metrics and system info metrics. See below section for a complete list of
metrics. Each group can be disabled with the --no-<group-name> flag.

### CLI Flags

```
Flags:
      --help                    Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=""               Config file
  -l, --listen="0.0.0.0"        Listen address
  -d, --debug                   Debug mode
  -v, --volume-fetch-period=2m  Period of asynchronously fetching volumes
      --no-aggregate            Disable aggregate collector
      --no-volume               Disable volume collector
      --no-system               Disable system collector
```

### Configuration

A configuration file needs to be provided via the `-c` or `--config` flag. By
default, the collector tries to use the "./netapp_filers.yaml". It should be a
yaml file with list of filer definitions with following fields.

```
- name: netapp-123
  host: netapp-123.labx.company
  availability_zone: az-a
  username: <username>
  password: <password>
```

The `username` and `password` field can be omitted in the yaml file, and set via
the env variables `NETAPP_USERNAME` and `NETAPP_PASSWORD`.

## Metrics

**Volume Metrics** with labels `availability_zone`, `filer`, `project_id`,
`share_id`, `volume` and `vserver`. <sup>1</sup>

- netapp_volume_state <sup>2</sup>
- netapp_volume_total_bytes
- netapp_volume_used_bytes
- netapp_volume_available_bytes
- netapp_volume_snapshot_used_bytes
- netapp_volume_snapshot_reserved_bytes
- netapp_volume_snapshot_available_bytes
- netapp_volume_percentage_snapshot_reserve
- netapp_volume_used_percentage
- netapp_volume_saved_total_percentage
- netapp_volume_saved_compression_percentage
- netapp_volume_saved_deduplication_percentage
- netapp_volume_is_encrypted
- netapp_volume_inode_files_total
- netapp_volume_inode_files_used
- netapp_volume_inode_files_used_percentage

<sup>1</sup> The label `project_id` is openstack specific, and `share_id` is
openstack manila specific.

<sup>2</sup> The metric netapp_volume_state being 1 means "online"; being -1
means "offline".

**Aggregate Metrics** with labels `availability_zone`, `filer`, `node` and
`aggregate`.

- netapp_aggregate_total_bytes
- netapp_aggregate_used_bytes
- netapp_aggregate_available_bytes
- netapp_aggregate_used_percentage
- netapp_aggregate_physical_used_bytes
- netapp_aggregate_physical_percentage
- netapp_aggregate_is_encrypted

**System Metrics** with labels `availability_zone` and `filer`.

- netapp_filer_system_version

## Version

Code is currently on v2, and is largely refactored to make extension easier. Old
version is available under tag
[v1](https://github.com/sapcc/netapp-api-exporter/releases/tag/v1).
