package main

import (
	"github.com/saichler/l8test/go/infra/t_service"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/reflect/go/reflect/introspecting"
)

var Plugin ifs.IPlugin = &TestPlugin{}

type TestPlugin struct {
}

func (this *TestPlugin) Install(vnic ifs.IVNic) error {
	vnic.Resources().Logger().Info("Registering Test Elements on ", vnic.Resources().SysConfig().LocalAlias)
	vnic.Resources().Registry().Register(&testtypes.TestProto{})
	node, err := vnic.Resources().Introspector().Inspect(&testtypes.TestProto{})
	if err != nil {
		return err
	}
	introspecting.AddPrimaryKeyDecorator(node, "MyString")

	vnic.Resources().Services().RegisterServiceHandlerType(&t_service.TestServiceHandler{})
	vnic.Resources().Services().RegisterServiceHandlerType(&t_service.TestServiceTransactionHandler{})
	vnic.Resources().Services().RegisterServiceHandlerType(&t_service.TestServiceReplicationHandler{})

	_, err = vnic.Resources().Services().Activate(t_service.ServiceType, t_service.ServiceName, 0, vnic.Resources(), nil, "plugin")
	if err != nil {
		return err
	}

	_, err = vnic.Resources().Services().Activate(t_service.ServiceTrType, t_service.ServiceName, 1, vnic.Resources(), nil, "plugin")
	if err != nil {
		return err
	}

	_, err = vnic.Resources().Services().Activate(t_service.ServiceRepType, t_service.ServiceName, 2, vnic.Resources(), vnic, "plugin")

	return nil
}
