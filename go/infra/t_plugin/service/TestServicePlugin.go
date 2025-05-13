package main

import (
	"github.com/saichler/l8test/go/infra/t_service"
	"github.com/saichler/l8types/go/ifs"
)

var Plugin ifs.IPlugin = &TestServicePlugin{}

type TestServicePlugin struct {
}

func (this *TestServicePlugin) Install(vnic ifs.IVNic) error {
	vnic.Resources().Logger().Info("#2 Registering Test Services on ", vnic.Resources().SysConfig().LocalAlias)
	vnic.Resources().Services().RegisterServiceHandlerType(&t_service.TestServiceHandler{})
	vnic.Resources().Services().RegisterServiceHandlerType(&t_service.TestServiceTransactionHandler{})
	vnic.Resources().Services().RegisterServiceHandlerType(&t_service.TestServiceReplicationHandler{})

	_, err := vnic.Resources().Services().Activate(t_service.ServiceType, t_service.ServiceName, 0, vnic.Resources(), nil, "plugin")
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
