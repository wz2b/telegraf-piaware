module telegraf_piaware

go 1.16

require (
	github.com/go-kit/kit v0.10.0
	github.com/influxdata/line-protocol v0.0.0-20210311194329-9aa0e372d097
	github.com/slim-bean/adsb-loki v0.0.0-20210709165408-4e075c72fe4b
)

replace k8s.io/client-go => k8s.io/client-go v12.0.0+incompatible
