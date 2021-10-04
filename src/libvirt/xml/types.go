package xml

import "encoding/xml"

type DomainXML struct {
	XMLName    xml.Name `xml:"domain"`
	Type       string   `xml:"type,attr"`
	Name       string   `xml:"name"`
	UUID       string   `xml:"uuid"`
	Metadata   metadata `xml:"metadata"`
	Cpu        cpu      `xml:"cpu"`
	Vcpu       vcpu     `xml:"vcpu"`
	Memory     memory   `xml:"memory"`
	CurrentMem memory   `xml:"currentMemory"`
	Os         os       `xml:"os"`
	OnPowerOff string   `xml:"on_poweroff"`
	OnReboot   string   `xml:"on_reboot"`
	OnCrash    string   `xml:"on_crash"`
	Suspend    suspend  `xml:"pm"`
	Features   features `xml:"features"`
	Devices    devices  `xml:"devices"`
	RNG        rng      `xml:"rng"`
}

type cpu struct {
	XMLName  string    `xml:"cpu"`
	Mode     string    `xml:"mode,attr"`
	Check    string    `xml:"check,attr"`
	Topology *topology `xml:"topology,omitempty"`
}

type topology struct {
	XMLName string `xml:"topology"`
	Sockets string `xml:"sockets,attr"`
	Dies    string `xml:"dies,attr"`
	Cores   string `xml:"cores,attr"`
	Threads string `xml:"threads,attr"`
}

type vcpu struct {
	XMLName   string `xml:"vcpu"`
	Placement string `xml:"placement,attr"`
	Value     uint   `xml:",chardata"`
}

type memory struct {
	//XMLName string `xml:"memory"`
	Unit  string `xml:"unit,attr"`
	Value uint   `xml:",chardata"`
}

type suspend struct {
	XMLName string        `xml:"pm"`
	ToDisk  suspendToDisk `xml:"suspend-to-disk"`
	ToMem   suspendToMem  `xml:"suspend-to-mem"`
}

type suspendToDisk struct {
	XMLName string  `xml:"suspend-to-disk"`
	Enabled yesOrNo `xml:"enabled,attr"`
}

type suspendToMem struct {
	XMLName string  `xml:"suspend-to-mem"`
	Enabled yesOrNo `xml:"enabled,attr"`
}

type metadata struct {
	XMLName   string `xml:"metadata"`
	LibOSInfo osInfo `xml:"libosinfo"`
}

type osInfo struct {
	XMLName string `xml:"libosinfo"`
	Os      osMeta `xml:"os"`
}

type osMeta struct {
	XMLName string `xml:"os"`
	Url     string `xml:"id,attr"`
}

type os struct {
	XMLName string `xml:"os"`
	Type    osType `xml:"type"`
	Boot    osBoot `xml:"boot"`
}

type osType struct {
	XMLName string `xml:"type"`
	Arch    string `xml:"arch,attr"`
	Machine string `xml:"machine,attr"`
	Value   string `xml:",chardata"`
}

type osBoot struct {
	XMLName string `xml:"boot"`
	Device  string `xml:"dev,attr"`
}

type features struct {
	XMLName string `xml:"features"`
	Pae     string `xml:"pae,omitempty"`
	Acpi    string `xml:"acpi,omitempty"`
}

type devices struct {
	XMLName     string          `xml:"devices"`
	Emulator    string          `xml:"emulator"`
	Disks       *[]disk         `xml:"disk"`
	Controllers *[]controller   `xml:"controller"`
	Interfaces  *[]domInterface `xml:"interface"`
	HostDevices *[]hostDevice   `xml:"hostdev"`
}

type disk struct {
	XMLName string      `xml:"disk"`
	Driver  diskDriver  `xml:"driver"`
	Source  diskSource  `xml:"source"`
	Target  diskTarget  `xml:"target"`
	Address diskAddress `xml:"address"`
}

type diskDriver struct {
	XMLName string `xml:"driver"`
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	Cache   string `xml:"cache,attr"`
	IO      string `xml:"io,attr"`
}

type diskSource struct {
	XMLName string `xml:"source"`
	Device  string `xml:"dev,attr"`
}

type diskTarget struct {
	XMLName string `xml:"target"`
	Device  string `xml:"dev,attr"`
	Bus     string `xml:"bus,attr"`
}

type diskAddress struct {
	XMLName    string `xml:"address"`
	Type       string `xml:"type,attr"`
	Controller uint   `xml:"controller,attr"`
	Bus        uint   `xml:"bus,attr"`
	Target     uint   `xml:"target,attr"`
}

type controller struct {
	XMLName string  `xml:"controller"`
	Type    string  `xml:"type,attr"`
	Index   uint    `xml:"index,attr"`
	Model   string  `xml:"model,attr"`
	Ports   uint    `xml:"ports,attr,omitempty"`
	Address Address `xml:"address"`
}

type hostDevice struct {
	XMLName    string  `xml:"hostdev"`
	SourceAddr Address `xml:"source>address"`
	Address    Address `xml:"address"`
}

type domInterface struct {
	XMLName string          `xml:"interface"`
	Mac     interfaceMac    `xml:"mac"`
	Source  interfaceSource `xml:"source"`
	Model   interfaceModel  `xml:"model"`
	Addr    Address         `xml:"address"`
}

type interfaceMac struct {
	XMLName string `xml:"mac"`
	Address string `xml:"address"`
}
type interfaceSource struct {
	XMLName string `xml:"source"`
	Network string `xml:"network"`
}
type interfaceModel struct {
	XMLName string `xml:"model"`
	Type    string `xml:"type"`
}

type rng struct {
	XMLName string     `xml:"rng"`
	Backend rngBackend `xml:"backend"`
	Address Address    `xml:"address"`
}

type rngBackend struct {
	XMLName string `xml:"backend"`
	Model   string `xml:"model,attr"`
	Value   string `xml:",chardata"`
}

type Address struct {
	XMLName   string `xml:"address"`
	Type      string `xml:"type,attr"`
	Domain    string `xml:"domain,attr"`
	Bus       string `xml:"bus,attr"`
	Slot      string `xml:"slot,attr"`
	Function  string `xml:"function,attr"`
	MultiFunc string `xml:"multifunction,attr,omitempty"`
}
