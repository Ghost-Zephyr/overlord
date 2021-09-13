package libvirt

import (
	log "bitsnthings.dev/overlord/src/log"
	"libvirt.org/go/libvirt"
)

func (virt *Libvirt) getAllDomains() {
	virt.Domains.Active = make(map[string]map[string]libvirt.Domain)
	virt.Domains.Inactive = make(map[string]map[string]libvirt.Domain)
	for uri, conn := range virt.Hosts {
		virt.Domains.Active[uri] = getDomainsByFlag(uri, conn, libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
		virt.Domains.Inactive[uri] = getDomainsByFlag(uri, conn, libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	}
}

func getDomainsByFlag(uri string, conn *libvirt.Connect, flag libvirt.ConnectListAllDomainsFlags) map[string]libvirt.Domain {
	domMap := make(map[string]libvirt.Domain)
	doms, err := conn.ListAllDomains(flag)
	if err != nil {
		log.PrintLog(
			log.ERROR,
			"Error fetching domain list from node with connection uri: \"%s\"! %s",
			uri, err)
	}
	for _, dom := range doms {
		uuid, err := dom.GetUUIDString()
		if err != nil {
			log.PrintLog(
				log.ERROR,
				"Error getting id of domain on node with uri: \"%s\"! %s",
				uri, err)
		} else {
			domMap[uuid] = dom
		}
	}
	return domMap
}
