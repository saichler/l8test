package t_topology

import (
	"github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/l8test/go/infra/t_servicepoints"
	"github.com/saichler/layer8/go/overlay/vnet"
	"github.com/saichler/layer8/go/overlay/vnic"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/testtypes"
	"time"
)

func createVnet(vnetPort int) *vnet.VNet {
	_resources, _ := t_resources.CreateResources(vnetPort, -1)
	_vnet := vnet.NewVNet(_resources)
	_vnet.Start()
	return _vnet
}

func createVnic(vnetPort int, vnicNum int, serviceArea int32) (common.IVirtualNetworkInterface, *t_servicepoints.TestServicePointHandler, *t_servicepoints.TestServicePointTransactionHandler, *t_servicepoints.TestServicePointReplicationHandler) {
	_resources, alias := t_resources.CreateResources(vnetPort, vnicNum)
	var handler *t_servicepoints.TestServicePointHandler
	var handlerTr *t_servicepoints.TestServicePointTransactionHandler
	var handlerRep *t_servicepoints.TestServicePointReplicationHandler

	if serviceArea != -1 {
		_resources.Registry().Register(&testtypes.TestProto{})
		_resources.ServicePoints().AddServicePointType(&t_servicepoints.TestServicePointHandler{})
		_resources.ServicePoints().AddServicePointType(&t_servicepoints.TestServicePointTransactionHandler{})
		_resources.ServicePoints().AddServicePointType(&t_servicepoints.TestServicePointReplicationHandler{})

		h, err := _resources.ServicePoints().Activate(t_servicepoints.ServicePointType, t_servicepoints.ServiceName, 0, _resources, nil, alias)
		if err != nil {
			panic(err)
		}
		handler = h.(*t_servicepoints.TestServicePointHandler)

		hTr, err := _resources.ServicePoints().Activate(t_servicepoints.ServicePointTrType, t_servicepoints.ServiceName, 1, _resources, nil, alias)
		if err != nil {
			panic(err)
		}
		handlerTr = hTr.(*t_servicepoints.TestServicePointTransactionHandler)
	}
	_vnic := vnic.NewVirtualNetworkInterface(_resources, nil)
	_vnic.Resources().SysConfig().KeepAliveIntervalSeconds = 30
	_vnic.Start()

	if serviceArea != -1 {
		_vnic.WaitForConnection()
		hRep, err := _resources.ServicePoints().Activate(t_servicepoints.ServicePointRepType, t_servicepoints.ServiceName, 2, _resources, _vnic, alias)
		if err != nil {
			panic(err)
		}
		handlerRep = hRep.(*t_servicepoints.TestServicePointReplicationHandler)
	}

	return _vnic, handler, handlerTr, handlerRep
}

func connectVnets(vnet1, vnet2 *vnet.VNet) {
	vnet1.ConnectNetworks("127.0.0.1", vnet2.Resources().SysConfig().VnetPort)
}

func Sleep() {
	time.Sleep(time.Millisecond * 100)
}
