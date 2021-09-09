package libvirt

import (
	log "bitsnthings.dev/overlord/src/log"
	"libvirt.org/go/libvirt"
)

func (virt *Libvirt) getAllDomains() {
	virt.getDomainsByFlag(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	virt.getDomainsByFlag(libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
}

func (virt *Libvirt) getDomainsByFlag(flag libvirt.ConnectListAllDomainsFlags) {
	for _, conn := range virt.Hosts {
		doms, err := conn.ListAllDomains(flag)
		uri, _ := conn.GetURI()
		if err != nil {
			log.PrintLog(
				log.ERROR,
				"Error fetching domain list from node with connection uri: \"%s\"! %s",
				uri, err)
		}
		virt.Guests[uri] = make(map[string]libvirt.Domain)
		for _, dom := range doms {
			name, err := dom.GetName()
			if err != nil {
				log.PrintLog(log.ERROR, "Error getting name of domain on node with connection string \"%s\"! %s",
					uri, err)
			}
			virt.Guests[uri][name] = dom
		}
	}
}
