package main

import (
	"fmt"
	"os"

	"github.com/dkaman/reg"
)

type Nozzle interface {
	Login(username, password string) (bool, error)
}

type NozzleDriver interface {
	New(opts map[string]string) (Nozzle, error)
}

type oktaDriver struct{}

func (*oktaDriver) New(opts map[string]string) (Nozzle, error) {
	return &okta{}, nil
}

type okta struct{}

func (*okta) Login(username, password string) (bool, error) {
	return true, nil
}

type adfsDriver struct{}

func (*adfsDriver) New(opts map[string]string) (Nozzle, error) {
	return &adfs{}, nil
}

type adfs struct{}

func (*adfs) Login(username, password string) (bool, error) {
	return false, nil
}

func init() {
	reg.Register[NozzleDriver]("okta", &oktaDriver{})
	reg.Register[NozzleDriver]("adfs", &adfsDriver{})
}

func main() {
	fmt.Printf("types: %v\n", reg.Types())
	fmt.Printf("drivers: %v\n", reg.Drivers[NozzleDriver]())

	name := "okta"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	fmt.Printf("opening %s\n", name)
	d, err := reg.Open[NozzleDriver](name)
	if err != nil {
		panic(err)
	}
	conf := map[string]string{"foo": "bar"}
	noz, err := d.New(conf)
	if err != nil {
		panic(err)
	}
	ok, err := noz.Login("username", "password")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%t\n", ok)
}
