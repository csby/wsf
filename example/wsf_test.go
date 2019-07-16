package example

import (
	"fmt"
	"github.com/csby/wsf/types"
	"testing"
	"time"
)

func Test_wsf(t *testing.T) {
	fmt.Println("certificate folder path:", crtFileFolder())
	err := createCertificates()
	if err != nil {
		t.Fatal(err)
	}
	defer deleteCertificates()

	server := &Server{}
	err = server.Run(func(server types.Server) {
		time.Sleep(time.Second * 4)

		client := &Client{}
		err := client.Run()
		select {
		case <-time.After(time.Second * 30):
		case <-err:

		}
		server.Shutdown()
	})
	if err != nil {
		t.Fatal(err)
	}
}
