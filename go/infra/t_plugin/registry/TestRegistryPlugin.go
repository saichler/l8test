package main

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
)

var Plugin ifs.IPlugin = &TestRegistryPlugin{}

type TestRegistryPlugin struct{}

func (this *TestRegistryPlugin) Install(vnic ifs.IVNic) error {
	vnic.Resources().Logger().Info("#2 Registering Test Elements on ", vnic.Resources().SysConfig().LocalAlias)
	vnic.Resources().Introspector().Clean("TestProto")
	vnic.Resources().Registry().UnRegister("TestProto")
	vnic.Resources().Registry().UnRegister("TestProtoSub")
	vnic.Resources().Registry().UnRegister("TestProtoSubSub")
	vnic.Resources().Registry().UnRegister("TestProtoList")

	vnic.Resources().Introspector().Decorators().AddPrimaryKeyDecorator(&testtypes.TestProto{}, "MyString")
	vnic.Resources().Registry().Register(&testtypes.TestProtoList{})
	return nil
}
