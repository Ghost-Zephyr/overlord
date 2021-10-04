/*
Package libvirt for the overlord VM management software!
Has a useful XML export for libvirt domain xml parsing.
*/
package libvirt

import (
	"errors"

	log "bitsnthings.dev/overlord/src/log"
	libvirt "libvirt.org/go/libvirt"
)

var errDomainNotFound = errors.New("Clould not find the domain in internal state!")

type libvirtConnectFunc func(string) (*libvirt.Connect, error)

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

func (virt *Libvirt) GetXMLFlags(uri string) libvirt.DomainXMLFlags {
	var flags libvirt.DomainXMLFlags
	if !virt.IsReadOnly(uri) {
		flags = libvirt.DOMAIN_XML_INACTIVE + libvirt.DOMAIN_XML_INACTIVE
	}
	return flags
}

func (virt *Libvirt) GetDomainAndStatus(uri string, uuid string) (libvirt.Domain, bool, error) {
	if host, exists := virt.Domains.Active[uri]; exists {
		dom, exists := host[uuid]
		if exists {
			return dom, true, nil
		}
		return libvirt.Domain{}, false, errDomainNotFound
	}
	if host, exists := virt.Domains.Inactive[uri]; exists {
		dom, exists := host[uuid]
		if exists {
			return dom, false, nil
		}
		return libvirt.Domain{}, false, errDomainNotFound
	}
	return libvirt.Domain{}, false, errDomainNotFound
}
