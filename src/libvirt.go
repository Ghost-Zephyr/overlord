package overlord

import (
	"log"

	libvirt "libvirt.org/go/libvirt"
)

type Libvirt struct {
	Hosts []*libvirt.Connect
	Auth  libvirt.ConnectAuth
	Flags libvirt.ConnectFlags
}

func (virt *Libvirt) Connect(cstr string) (*libvirt.Connect, error) {
	conn, err := libvirt.NewConnect(cstr) //, &virt.Auth, virt.Flags)
	if err != nil {
		log.Printf("Error connecting to libvirt host! %s", err)
	}
	return conn, err
}

func (virt *Libvirt) ConnectMany(cstrs []string) {
	for _, cstr := range cstrs {
		virt.Connect(cstr)
	}
}
