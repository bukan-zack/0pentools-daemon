package domain

import (
	"bytes"
	"fmt"
	"html/template"

	lv "github.com/digitalocean/go-libvirt"
	"github.com/gofrs/uuid"
)

type Domain struct {
	Name  string `json:"name"`
	UUID  string `json:"uuid"`
	State string `json:"state"`
}

type XmlContext struct {
	UUID string
	RAM  uint64
	Port uint64
}

func GetAllDomains() ([]Domain, error) {
	var domains []Domain
	var lib lv.Libvirt

	res, _, err := lib.ConnectListAllDomains(1, 0)
	if err != nil {
		panic(err)
	}
	if len(res) < 1 {
		return domains, nil
	}

	for _, d := range res {
		state, _, _, _, _, _ := lib.DomainGetInfo(d)
		dom := Domain{
			Name:  d.Name,
			UUID:  fmt.Sprintf("%x", d.UUID),
			State: GetDomainState(state),
		}
		domains = append(domains, dom)
	}

	return domains, nil
}

func loadXML(path string, ctx XmlContext) (string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	tmpl.Execute(&buffer, ctx)
	return buffer.String(), nil
}

func Create(RAM uint64, Port uint64) (string, error) {
	var lib lv.Libvirt

	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	id := uuid.String()
	ctx := XmlContext{
		UUID: id,
		RAM:  RAM,
		Port: Port,
	}

	fxml, err := loadXML("volume.xml.txt", ctx)
	if err != nil {
		return "", err
	}
	pool, err := lib.StoragePoolLookupByName("default")
	if err != nil {
		return "", err
	}
	_, err = lib.StorageVolCreateXML(pool, fxml, 0)
	if err != nil {
		return "", err
	}

	vxml, err := loadXML("vm.xml.txt", ctx)
	if err != nil {
		return "", err
	}
	dom, err := lib.DomainDefineXML(vxml)
	if err != nil {
		return "", err
	}

	return dom.Name, nil
}

func Start(uuid string) error {
	var lib lv.Libvirt
	dom, err := lib.DomainLookupByName(uuid)
	if err != nil {
		return err
	}
	err = lib.DomainCreate(dom)
	if err != nil {
		return err
	}
	return nil
}

func Stop(uuid string) error {
	var lib lv.Libvirt
	dom, err := lib.DomainLookupByName(uuid)
	if err != nil {
		return err
	}
	err = lib.DomainDestroy(dom)
	if err != nil {
		return err
	}
	return nil
}

func Restart(uuid string) error {
	var lib lv.Libvirt
	dom, err := lib.DomainLookupByName(uuid)
	if err != nil {
		return err
	}
	err = lib.DomainReboot(dom, 0)
	if err != nil {
		return err
	}
	return nil
}
