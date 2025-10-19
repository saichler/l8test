package t_topology

import (
	"time"

	"github.com/saichler/l8bus/go/overlay/vnet"
	"github.com/saichler/l8bus/go/overlay/vnic"
	"github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/l8test/go/infra/t_service"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
)

func createVnet(vnetPort int, level ifs.LogLevel) *vnet.VNet {
	_resources, _ := t_resources.CreateResources(vnetPort, -1, level)
	_vnet := vnet.NewVNet(_resources)
	_vnet.Start()
	return _vnet
}

func createVnic(vnetPort int, vnicNum int, serviceArea int32, level ifs.LogLevel) (ifs.IVNic, *t_service.TestServiceHandler, *t_service.TestServiceTransactionHandler, *t_service.TestServiceReplicationHandler) {
	_resources, alias := t_resources.CreateResources(vnetPort, vnicNum, level)
	var handler *t_service.TestServiceHandler
	var handlerTr *t_service.TestServiceTransactionHandler
	var handlerRep *t_service.TestServiceReplicationHandler

	if serviceArea != -1 {
		_resources.Registry().Register(&testtypes.TestProtoList{})
		_resources.Registry().Register(&testtypes.TestProto{})
	}

	_vnic := vnic.NewVirtualNetworkInterface(_resources, nil)
	_vnic.Resources().SysConfig().KeepAliveIntervalSeconds = 30
	_vnic.Start()

	if serviceArea != -1 {
		_vnic.WaitForConnection()

		sla := ifs.NewServiceLevelAgreement(&t_service.TestServiceHandler{}, t_service.ServiceName, 0, false, nil)
		sla.SetArgs(alias)
		h, err := _resources.Services().Activate(sla, _vnic)
		if err != nil {
			panic(err)
		}
		handler = h.(*t_service.TestServiceHandler)

		sla = ifs.NewServiceLevelAgreement(&t_service.TestServiceTransactionHandler{}, t_service.ServiceName, 1, true, nil)
		sla.SetArgs(alias)
		hTr, err := _resources.Services().Activate(sla, _vnic)
		if err != nil {
			panic(err)
		}
		handlerTr = hTr.(*t_service.TestServiceTransactionHandler)

		sla = ifs.NewServiceLevelAgreement(&t_service.TestServiceReplicationHandler{}, t_service.ServiceName, 2, true, nil)
		sla.SetArgs(alias)
		hRep, err := _resources.Services().Activate(sla, _vnic)
		if err != nil {
			panic(err)
		}
		handlerRep = hRep.(*t_service.TestServiceReplicationHandler)
	}

	return _vnic, handler, handlerTr, handlerRep
}

func connectVnets(vnet1, vnet2 *vnet.VNet) {
	vnet1.ConnectNetworks("127.0.0.1", vnet2.Resources().SysConfig().VnetPort)
}

func Sleep() {
	time.Sleep(time.Millisecond * 100)
}
