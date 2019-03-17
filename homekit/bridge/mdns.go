package bridge

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/brutella/dnssd"
)

func newService(config *Config) dnssd.Service {
	// 2016-03-14(brutella): Replace whitespaces (" ") from service name
	// with underscores ("_")to fix invalid http host header field value
	// produces by iOS.
	//
	// [Radar] http://openradar.appspot.com/radar?id=4931940373233664
	stripped := strings.Replace(config.name, " ", "_", -1)

	var ips []net.IP
	if config.IP != "" {
		if ip := net.ParseIP(config.IP); ip != nil {
			ips = append(ips, ip)
		}
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "hap"
	}

	sdConfig := dnssd.Config{
		Name:   stripped,
		Type:   "_hap._tcp",
		Host:   fmt.Sprintf("%s-%s", hostname, stripped),
		Domain: "local",
		IPs:    ips,
		Port:   config.servePort,
	}
	service, _ := dnssd.NewService(sdConfig)
	service.Text = config.txtRecords()

	return service
}
