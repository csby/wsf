package example

import (
	"github.com/csby/security/certificate"
	"os"
	"path/filepath"
	"runtime"
)

var (
	caCrtFilePath     = getCrtFilePath("ca.crt")
	serverCrtFilePath = getCrtFilePath("server.pfx")
	serverCrtPassword = "server123"
	serverCrtOU       = "server-test"
	clientCrtFilePath = getCrtFilePath("client.pfx")
	clientCrtPassword = "client123"
	clientCrtOU       = "client-test"
)

func createCertificates() error {
	// ca
	caPrivate := &certificate.RSAPrivate{}
	err := caPrivate.Create(2048)
	if err != nil {
		return err
	}
	caPublic, err := caPrivate.Public()
	if err != nil {
		return err
	}
	caTemplate := &certificate.CrtTemplate{
		Organization:       "ca",
		OrganizationalUnit: "ca-test",
	}
	template, err := caTemplate.Template()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	caCrt := &certificate.Crt{}
	err = caCrt.Create(template, template, caPublic, caPrivate)
	if err != nil {
		return err
	}
	err = caCrt.ToFile(caCrtFilePath)
	if err != nil {
		return err
	}

	// server
	serverPrivate := &certificate.RSAPrivate{}
	err = serverPrivate.Create(2048)
	if err != nil {
		return err
	}
	serverPublic, err := serverPrivate.Public()
	if err != nil {
		return err
	}
	serverTemplate := &certificate.CrtTemplate{
		Organization:       "server",
		OrganizationalUnit: serverCrtOU,
		Hosts: []string{
			"127.0.0.1",
			"localhost",
		},
	}
	template, err = serverTemplate.Template()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	serverCrt := &certificate.CrtPfx{}
	err = serverCrt.Create(template, caCrt.Certificate(), serverPublic, caPrivate)
	if err != nil {
		return err
	}
	err = serverCrt.ToFile(serverCrtFilePath, caCrt, serverPrivate, serverCrtPassword)
	if err != nil {
		return err
	}

	// client
	clientPrivate := &certificate.RSAPrivate{}
	err = clientPrivate.Create(2048)
	if err != nil {
		return err
	}
	clientPublic, err := clientPrivate.Public()
	if err != nil {
		return err
	}
	clientTemplate := &certificate.CrtTemplate{
		Organization:       "client",
		OrganizationalUnit: clientCrtOU,
	}
	template, err = clientTemplate.Template()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	clientCrt := &certificate.CrtPfx{}
	err = clientCrt.Create(template, caCrt.Certificate(), clientPublic, caPrivate)
	if err != nil {
		return err
	}
	err = clientCrt.ToFile(clientCrtFilePath, caCrt, clientPrivate, clientCrtPassword)
	if err != nil {
		return err
	}

	return nil
}

func deleteCertificates() error {
	folder := crtFileFolder()
	return os.RemoveAll(folder)
}

func getCrtFilePath(name string) string {
	return filepath.Join(crtFileFolder(), name)
}

func crtFileFolder() string {
	_, file, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(file), "crts")
}
