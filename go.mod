module github.com/lukasmalkmus/faktory_exporter

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4
	github.com/beorn7/perks v1.0.0
	github.com/contribsys/faktory v1.4.0-1
	github.com/golang/protobuf v1.3.2
	github.com/google/go-github/v25 v25.0.0 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.10.0
	github.com/prometheus/procfs v0.0.2
	github.com/prometheus/promu v0.2.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a // indirect
	golang.org/x/sys v0.0.0-20190422165155-953cdadca894
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.8.0
