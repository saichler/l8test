package main

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/reflect/go/reflect/introspecting"
)

var Plugin ifs.IPlugin = &TestRegistryPlugin{}

type TestRegistryPlugin struct{}

func (this TestRegistryPlugin) Install(vnic ifs.IVNic) error {
	vnic.Resources().Logger().Info("#2 Registering Test Elements on ", vnic.Resources().SysConfig().LocalAlias)
	vnic.Resources().Introspector().Clean("TestProto")
	vnic.Resources().Registry().UnRegister("TestProto")
	vnic.Resources().Registry().UnRegister("TestProtoSub")
	vnic.Resources().Registry().UnRegister("TestProtoSubSub")
	node, err := vnic.Resources().Introspector().Inspect(&testtypes.TestProto{})
	if err != nil {
		return err
	}
	introspecting.AddPrimaryKeyDecorator(node, "MyString")
	return nil
}
