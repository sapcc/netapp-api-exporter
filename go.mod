module github.com/sapcc/netapp-api-exporter

go 1.14

require (
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883 // indirect
	github.com/pepabo/go-netapp v0.0.0-20200708032902-3c5b98f52cf4
	github.com/prometheus/client_golang v1.7.1
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/pepabo/go-netapp v0.0.0-20200708032902-3c5b98f52cf4 => ../../../chuan137/go-netapp
)
