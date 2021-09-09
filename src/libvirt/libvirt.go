package libvirt

import (
	log "bitsnthings.dev/overlord/src/log"
	libvirt "libvirt.org/go/libvirt"
)

type Libvirt struct {
	Hosts  map[string]*libvirt.Connect
	Guests map[string]map[string]libvirt.Domain
	// Auth   libvirt.ConnectAuth
	// Flags  libvirt.ConnectFlags
}

func New() Libvirt {
	return Libvirt{
		map[string]*libvirt.Connect{},
		map[string]map[string]libvirt.Domain{},
	}
}

func (virt *Libvirt) GetStatus() {
	log.PrintLog(log.TRACE, "Getting libvirt cluster status.")
	virt.getAllDomains()
	log.PrintLog(log.INFO, "All domains in the cluster:")
	for uri, doms := range virt.Guests {
		log.PrintLog(log.INFO, "With uri: %s", uri)
		for name, dom := range doms {
			os, _ := dom.GetOSType()
			log.PrintLog(log.INFO, "%s (%s)", name, os)
		}
	}
}

func (virt *Libvirt) Connect(cstr string) (*libvirt.Connect, error) {
	conn, err := libvirt.NewConnect(cstr) //, &virt.Auth, virt.Flags)
	if err != nil {
		log.PrintLog(log.WARN, "Error connecting to libvirt host with uri: \"%s\"! %s",
			cstr, err)
		return conn, err
	}
	virt.Hosts[cstr] = conn
	return conn, err
}

func (virt *Libvirt) ConnectMany(cstrs []string) {
	for _, cstr := range cstrs {
		virt.Connect(cstr)
	}
}
