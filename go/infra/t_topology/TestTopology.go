package t_topology

import (
	. "github.com/saichler/l8test/go/infra/t_resources"
	. "github.com/saichler/l8test/go/infra/t_servicepoints"
	. "github.com/saichler/layer8/go/overlay/vnet"
	. "github.com/saichler/layer8/go/overlay/vnic"
	. "github.com/saichler/types/go/common"
	"sync"
)

type TestTopology struct {
	vnets      map[string]*VNet
	vnetsOrder []*VNet
	vnics      map[string]IVirtualNetworkInterface
	handlers   map[string]*TestServicePointHandler
	mtx        *sync.RWMutex
}

func NewTestTopology(vnicCountPervNet int, vnetPorts ...int) *TestTopology {
	this := &TestTopology{}
	this.vnets = make(map[string]*VNet)
	this.vnics = make(map[string]IVirtualNetworkInterface)
	this.handlers = make(map[string]*TestServicePointHandler)
	this.vnetsOrder = make([]*VNet, 0)
	this.mtx = &sync.RWMutex{}

	for _, vNetPort := range vnetPorts {
		_vnet := createVnet(vNetPort)
		this.vnets[_vnet.Resources().SysConfig().LocalAlias] = _vnet
		this.vnetsOrder = append(this.vnetsOrder, _vnet)
	}
	Sleep()

	for _, vNetPort := range vnetPorts {
		for i := 0; i < vnicCountPervNet; i++ {
			if i == vnicCountPervNet-1 {
				_vnic, _ := createVnic(vNetPort, i+1, -1)
				this.vnics[_vnic.Resources().SysConfig().LocalAlias] = _vnic
			} else {
				_vnic, handler := createVnic(vNetPort, i+1, 0)
				this.vnics[_vnic.Resources().SysConfig().LocalAlias] = _vnic
				this.handlers[_vnic.Resources().SysConfig().LocalAlias] = handler
			}
			Sleep()
		}
	}

	Sleep()
	Sleep()
	for i := 0; i < len(this.vnetsOrder)-1; i++ {
		for j := i + 1; j < len(this.vnetsOrder); j++ {
			connectVnets(this.vnetsOrder[i], this.vnetsOrder[j])
			Sleep()
			Sleep()
		}
	}
	return this
}

func (this *TestTopology) Vnet(vnetPort int) *VNet {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, -1)
	return this.vnets[alias]
}

func (this *TestTopology) VnicByPort(vnetPort, vnicNum int) IVirtualNetworkInterface {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, vnicNum)
	return this.vnics[alias]
}

func (this *TestTopology) VnicByVnetNum(vnetNum, vnicNum int) IVirtualNetworkInterface {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	vnetPort := int(this.vnetsOrder[vnetNum-1].Resources().SysConfig().VnetPort)
	alias := AliasOf(vnetPort, vnicNum)
	return this.vnics[alias]
}

func (this *TestTopology) HandlerByPort(vnetPort, vnicNum int) *TestServicePointHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, vnicNum)
	return this.handlers[alias]
}

func (this *TestTopology) HandlerByVnetNum(vnetNum, vnicNum int) *TestServicePointHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	vnetPort := int(this.vnetsOrder[vnetNum-1].Resources().SysConfig().VnetPort)
	alias := AliasOf(vnetPort, vnicNum)
	return this.handlers[alias]
}

func (this *TestTopology) Shutdown() {
	for _, _vnic := range this.vnics {
		_vnic.Shutdown()
	}
	for _, _vnet := range this.vnets {
		_vnet.Shutdown()
	}
}

func (this *TestTopology) ResetHandlers() {
	for _, _handler := range this.handlers {
		_handler.Reset()
	}
}

func (this *TestTopology) AllHandlers() []*TestServicePointHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]*TestServicePointHandler, 0)
	for _, h := range this.handlers {
		result = append(result, h)
	}
	return result
}

func (this *TestTopology) AllVnics() []IVirtualNetworkInterface {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]IVirtualNetworkInterface, 0)
	for _, vnic := range this.vnics {
		result = append(result, vnic)
	}
	return result
}

func (this *TestTopology) RenewVnic(alias string) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	nic, ok := this.vnics[alias]
	if ok {
		nic.Shutdown()
		delete(this.vnics, alias)
		r := nic.Resources()
		r.SysConfig().LocalUuid = ""
		r.SysConfig().RemoteUuid = ""
		nic = NewVirtualNetworkInterface(nic.Resources(), nil)
		nic.Start()
		this.vnics[alias] = nic
	} else {
		Log.Error("Unable to find vnic ", alias)
	}
}

func (this *TestTopology) SetLogLevel(lvl LogLevel) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	for _, net := range this.vnets {
		net.Resources().Logger().SetLogLevel(lvl)
	}
	for _, nic := range this.vnics {
		nic.Resources().Logger().SetLogLevel(lvl)
	}
}

func (this *TestTopology) ReActivateTestService(nic IVirtualNetworkInterface) {
	h, err := nic.Resources().ServicePoints().Activate(ServicePointType, ServiceName, 0, nic.Resources(), nil,
		nic.Resources().SysConfig().LocalAlias)
	if err != nil {
		panic(err)
	}
	handler := h.(*TestServicePointHandler)
	this.handlers[nic.Resources().SysConfig().LocalAlias] = handler
}
