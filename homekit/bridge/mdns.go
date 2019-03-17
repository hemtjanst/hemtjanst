package bridge

import (
	"context"
	"net"
	"strings"

	"github.com/brutella/dnssd"
	"github.com/brutella/hc/log"
)

// MDNSService represents a mDNS service.
type MDNSService struct {
	config    *Config
	responder dnssd.Responder
	handle    dnssd.ServiceHandle
}

func newService(config *Config) dnssd.Service {
	// 2016-03-14(brutella): Replace whitespaces (" ") from service name
	// with underscores ("_")to fix invalid http host header field value
	// produces by iOS.
	//
	// [Radar] http://openradar.appspot.com/radar?id=4931940373233664
	stripped := strings.Replace(config.name, " ", "_", -1)

	var ips []net.IP
	if ip := net.ParseIP(config.IP); ip != nil {
		ips = append(ips, ip)
	}

	sdConfig := dnssd.Config{
		Name:   stripped,
		Type:   "_hap._tcp",
		Domain: "local",
		IPs:    ips,
		Port:   config.servePort,
	}
	service, _ := dnssd.NewService(sdConfig)
	service.Text = config.txtRecords()

	return service
}

// NewMDNSService returns a new service based for the bridge name, id and port.
func NewMDNSService(config *Config) *MDNSService {
	// TODO handle error
	responder, _ := dnssd.NewResponder()

	return &MDNSService{
		config:    config,
		responder: responder,
	}
}

// Publish announces the service for the machine's ip address on a random port using mDNS.
func (s *MDNSService) Publish(ctx context.Context) error {
	// 2016-03-14(brutella): Replace whitespaces (" ") from service name
	// with underscores ("_")to fix invalid http host header field value
	// produces by iOS.
	//
	// [Radar] http://openradar.appspot.com/radar?id=4931940373233664
	stripped := strings.Replace(s.config.name, " ", "_", -1)

	var ips []net.IP
	if ip := net.ParseIP(s.config.IP); ip != nil {
		ips = append(ips, ip)
	}

	sdConfig := dnssd.Config{
		Name:   stripped,
		Type:   "_hap._tcp",
		Domain: "local",
		IPs:    ips,
		Port:   s.config.servePort,
	}

	service, _ := dnssd.NewService(sdConfig)
	service.Text = s.config.txtRecords()
	handle, err := s.responder.Add(service)
	if err != nil {
		log.Info.Panic(err)
	}

	s.handle = handle

	return s.responder.Respond(ctx)
}

// Update updates the mDNS txt records.
func (s *MDNSService) Update() {
	if s.handle != nil {
		txt := s.config.txtRecords()
		s.handle.UpdateText(txt, s.responder)
		log.Debug.Println(txt)
	}
}

// Stop stops the running mDNS service.
func (s *MDNSService) Stop() {
	if s.handle != nil {
		s.responder.Remove(s.handle)
		s.handle = nil
	}
}
