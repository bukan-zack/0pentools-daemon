package domain

import (
	"bytes"
	"html/template"
	"fmt"
	"github.com/digitalocean/go-libvirt"
	"github.com/gofrs/uuid"
)

type Manager struct {
	lv *libvirt.Libvirt
}

func NewManager(lv *libvirt.Libvirt) *Manager {
	return &Manager{
		lv: lv,
	}
}

func (m *Manager) GetDomains() ([]Domain, error) {
	var ds []Domain
	lds, _, err := m.lv.ConnectListAllDomains(1, 0)
	if err != nil {
		return nil, err
	}

	for _, d := range lds {
		state, _, _, _, _, _ := m.lv.DomainGetInfo(d)
		ds = append(ds, Domain{
			Name: d.Name,
			UUID: fmt.Sprintf("%x", d.UUID),
			ID: d.ID,
			State: ConvDomainState(state),
		})
	}

	return ds, nil
}

func (m *Manager) ParseXML(path string, data DomainXML) (string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	b := bytes.Buffer{}
	tmpl.Execute(&b, data)

	return b.String(), nil
}

func (m *Manager) CreateDomain(memory uint64, port uint64) (string, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	ctx := DomainXML{
		UUID: uuid.String(),
		Memory: memory,
		Port: port,
	}

	fxml, err := m.ParseXML("volume.xml.txt", ctx)
	if err != nil {
		return "", err
	}

	pool, err := m.lv.StoragePoolLookupByName("default")
	if err != nil {
		return "", err
	}

	if _, err := m.lv.StorageVolCreateXML(pool, fxml, 0); err != nil {
		return "", err
	}

	dxml, err := m.ParseXML("vm.xml.txt", ctx)
	if err != nil {
		return "", err
	}

	d, err := m.lv.DomainDefineXML(dxml)
	if err != nil {
		return "", err
	}

	return d.Name, nil
}

func (m *Manager) StartDomain(uuid string) error {
	d, err := m.lv.DomainLookupByName(uuid)
	if err != nil {
		return err
	}

	return m.lv.DomainCreate(d)
}

func (m *Manager) StopDomain(uuid string) error {
	d, err := m.lv.DomainLookupByName(uuid)
	if err != nil {
		return err
	}

	return m.lv.DomainDestroy(d)
}

func (m *Manager) RestartDomain(uuid string) error {
	d, err := m.lv.DomainLookupByName(uuid)
	if err != nil {
		return err
	}

	return m.lv.DomainReboot(d, 0)
}