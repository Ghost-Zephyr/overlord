package libvirt

import (
	log "bitsnthings.dev/overlord/src/log"
	libvirt "libvirt.org/go/libvirt"
)

type libvirtConnectFunc func(string) (*libvirt.Connect, error)

type Domains struct {
	Active   map[string]map[string]libvirt.Domain
	Inactive map[string]map[string]libvirt.Domain
}

type Stats struct {
}

type Libvirt struct {
	Hosts   map[string]*libvirt.Connect
	ROHosts map[string]*libvirt.Connect
	Domains Domains
	// Auth   libvirt.ConnectAuth
	// Flags  libvirt.ConnectFlags
	Stats Stats
}

func New() Libvirt {
	return Libvirt{
		map[string]*libvirt.Connect{},
		map[string]*libvirt.Connect{},
		Domains{},
		Stats{},
	}
}

func (virt *Libvirt) GetStatus() {
	virt.getAllDomains()
	log.PrintLog(log.TRACE, "Fetched libvirt cluster status.")
}

func (virt *Libvirt) IsReadOnly(uri string) bool {
	if _, ok := virt.ROHosts[uri]; ok {
		return true
	} else {
		return false
	}
}

func (virt *Libvirt) Connect(cstr string, connFn libvirtConnectFunc) (*libvirt.Connect, error) {
	conn, err := connFn(cstr) //, &virt.Auth, virt.Flags)
	if err != nil {
		log.PrintLog(log.WARN, "Error connecting to libvirt host with uri: \"%s\"! %s",
			cstr, err)
		return conn, err
	}
	virt.Hosts[cstr] = conn
	return conn, err
}

func (virt *Libvirt) ConnectMany(cstrs []string, connFn libvirtConnectFunc) {
	for _, cstr := range cstrs {
		virt.Connect(cstr, connFn)
	}
}
