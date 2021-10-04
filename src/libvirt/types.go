package libvirt

import (
	"bitsnthings.dev/overlord/src/libvirt/xml"
	libvirt "libvirt.org/go/libvirt"
)

type Domains struct {
	Active   map[string]map[string]libvirt.Domain
	Inactive map[string]map[string]libvirt.Domain
}

type Stats struct {
}

type Libvirt struct {
	Hosts      map[string]*libvirt.Connect
	ROHosts    map[string]*libvirt.Connect
	DomainXMLs map[string]map[string]xml.DomainXML
	Domains    Domains
	// Auth   libvirt.ConnectAuth
	// Flags  libvirt.ConnectFlags
	Stats Stats
}

func New() Libvirt {
	return Libvirt{
		map[string]*libvirt.Connect{},
		map[string]*libvirt.Connect{},
		map[string]map[string]xml.DomainXML{},
		Domains{},
		Stats{},
	}
}
