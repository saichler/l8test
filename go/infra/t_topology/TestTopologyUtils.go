package t_topology

import (
	"github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/l8test/go/infra/t_servicepoints"
	"github.com/saichler/layer8/go/overlay/vnet"
	"github.com/saichler/layer8/go/overlay/vnic"
	"github.com/saichler/types/go/common"
	"time"
)

func createVnet(vnetPort int) *vnet.VNet {
	_resources, _ := t_resources.CreateResources(vnetPort, -1)
	_vnet := vnet.NewVNet(_resources)
	_vnet.Start()
	return _vnet
}

func createVnic(vnetPort int, vnicNum int, serviceArea int32) (common.IVirtualNetworkInterface, *t_servicepoints.TestServicePointHandler) {
	_resources, alias := t_resources.CreateResources(vnetPort, vnicNum)
	var handler *t_servicepoints.TestServicePointHandler
	if serviceArea != -1 {
		handler = t_servicepoints.NewTestServicePointHandler(alias)
		err := _resources.ServicePoints().RegisterServicePoint(handler, uint16(serviceArea))
		if err != nil {
			panic(err)
		}
	}
	_vnic := vnic.NewVirtualNetworkInterface(_resources, nil)
	_vnic.Resources().SysConfig().KeepAliveIntervalSeconds = 30
	_vnic.Start()
	return _vnic, handler
}

func connectVnets(vnet1, vnet2 *vnet.VNet) {
	vnet1.ConnectNetworks("127.0.0.1", vnet2.Resources().SysConfig().VnetPort)
}

func Sleep() {
	time.Sleep(time.Millisecond * 100)
}
