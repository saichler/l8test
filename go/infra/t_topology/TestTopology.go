package t_topology

import (
	"fmt"
	. "github.com/saichler/l8test/go/infra/t_resources"
	. "github.com/saichler/l8test/go/infra/t_service"
	. "github.com/saichler/l8types/go/ifs"
	"github.com/saichler/layer8/go/overlay/health"
	"github.com/saichler/layer8/go/overlay/protocol"
	. "github.com/saichler/layer8/go/overlay/vnet"
	. "github.com/saichler/layer8/go/overlay/vnic"
	"strconv"
	"sync"
)

type TestTopology struct {
	vnets       map[string]*VNet
	vnetsOrder  []*VNet
	vnics       map[string]IVNic
	handlers    map[string]*TestServiceHandler
	trHandlers  map[string]*TestServiceTransactionHandler
	repHandlers map[string]*TestServiceReplicationHandler
	mtx         *sync.RWMutex
}

func NewTestTopology(vnicCountPervNet int, vnetPorts []int, level LogLevel) *TestTopology {
	protocol.Discovery_Enabled = false
	this := &TestTopology{}
	this.vnets = make(map[string]*VNet)
	this.vnics = make(map[string]IVNic)
	this.handlers = make(map[string]*TestServiceHandler)
	this.trHandlers = make(map[string]*TestServiceTransactionHandler)
	this.repHandlers = make(map[string]*TestServiceReplicationHandler)
	this.vnetsOrder = make([]*VNet, 0)
	this.mtx = &sync.RWMutex{}

	for _, vNetPort := range vnetPorts {
		_vnet := createVnet(vNetPort, level)
		this.vnets[_vnet.Resources().SysConfig().LocalAlias] = _vnet
		this.vnetsOrder = append(this.vnetsOrder, _vnet)
	}

	for i := 0; i < len(this.vnetsOrder)-1; i++ {
		for j := i + 1; j < len(this.vnetsOrder); j++ {
			connectVnets(this.vnetsOrder[i], this.vnetsOrder[j])
		}
	}

	if !WaitForCondition(this.areVnetsConnected, 2, nil, "Vnets are not ready and connected") {
		panic("Vnets are not ready and connects ")
	}

	for _, vNetPort := range vnetPorts {
		for i := 0; i < vnicCountPervNet; i++ {
			if i == vnicCountPervNet-1 {
				_vnic, _, _, _ := createVnic(vNetPort, i+1, -1, level)
				this.vnics[_vnic.Resources().SysConfig().LocalAlias] = _vnic
			} else {
				_vnic, handler, trHandler, repHandler := createVnic(vNetPort, i+1, 0, level)
				this.vnics[_vnic.Resources().SysConfig().LocalAlias] = _vnic
				this.handlers[_vnic.Resources().SysConfig().LocalAlias] = handler
				this.trHandlers[_vnic.Resources().SysConfig().LocalAlias] = trHandler
				this.repHandlers[_vnic.Resources().SysConfig().LocalAlias] = repHandler
			}
		}
	}

	if !WaitForCondition(this.areVnetsHaveAllVnics, 3, nil, "Vnet are not ready 2") {
		for _, vnet := range this.vnets {
			hc := health.Health(vnet.Resources())
			fmt.Println(vnet.Resources().SysConfig().LocalAlias, " ", vnet.ExternalCount(), vnet.LocalCount(), len(hc.All()))
		}
		panic("Vnet are not ready 2")
	}

	if !WaitForCondition(this.areVnicsReady, 5, nil, "Vnics are not ready!") {
		vnicName := ""
		vnicSum := 0
		for vnetNum := 1; vnetNum <= 3; vnetNum++ {
			for vnicNum := 1; vnicNum <= 4; vnicNum++ {
				nic := this.VnicByVnetNum(vnetNum, vnicNum)
				hc := health.Health(nic.Resources())
				hp := hc.All()
				if len(hp) != 15 {
					vnicName = nic.Resources().SysConfig().LocalAlias
					vnicSum = len(hp)
					break
				}
			}
			if vnicName != "" {
				break
			}
		}
		panic("Vnics are not ready, vnic " + vnicName + " has only " + strconv.Itoa(vnicSum) + " instead of 15")
	}
	return this
}

func (this *TestTopology) areVnicsReady() bool {
	for vnetNum := 1; vnetNum <= 3; vnetNum++ {
		for vnicNum := 1; vnicNum <= 4; vnicNum++ {
			nic := this.VnicByVnetNum(vnetNum, vnicNum)
			hc := health.Health(nic.Resources())
			hp := hc.All()
			if len(hp) != 15 {
				return false
			}
		}
	}
	return true
}

func (this *TestTopology) areVnetsConnected() bool {
	for _, vnet := range this.vnets {
		if vnet.ExternalCount() != 2 {
			return false
		}
	}
	return true
}

func (this *TestTopology) areVnetsHaveAllVnics() bool {
	for _, vnet := range this.vnets {
		if vnet.ExternalCount() != 2 {
			return false
		}
		if vnet.LocalCount() != 4 {
			return false
		}
		hc := health.Health(vnet.Resources())
		if len(hc.All()) != 15 {
			return false
		}
	}
	return true
}

func (this *TestTopology) Vnet(vnetPort int) *VNet {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, -1)
	return this.vnets[alias]
}

func (this *TestTopology) VnicByPort(vnetPort, vnicNum int) IVNic {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, vnicNum)
	return this.vnics[alias]
}

func (this *TestTopology) VnicByVnetNum(vnetNum, vnicNum int) IVNic {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	vnetPort := int(this.vnetsOrder[vnetNum-1].Resources().SysConfig().VnetPort)
	alias := AliasOf(vnetPort, vnicNum)
	return this.vnics[alias]
}

func (this *TestTopology) HandlerByPort(vnetPort, vnicNum int) *TestServiceHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	alias := AliasOf(vnetPort, vnicNum)
	return this.handlers[alias]
}

func (this *TestTopology) HandlerByVnetNum(vnetNum, vnicNum int) *TestServiceHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	vnetPort := int(this.vnetsOrder[vnetNum-1].Resources().SysConfig().VnetPort)
	alias := AliasOf(vnetPort, vnicNum)
	return this.handlers[alias]
}

func (this *TestTopology) TrHandlerByVnetNum(vnetNum, vnicNum int) *TestServiceTransactionHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	vnetPort := int(this.vnetsOrder[vnetNum-1].Resources().SysConfig().VnetPort)
	alias := AliasOf(vnetPort, vnicNum)
	return this.trHandlers[alias]
}

func (this *TestTopology) RepHandlerByVnetNum(vnetNum, vnicNum int) *TestServiceReplicationHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	vnetPort := int(this.vnetsOrder[vnetNum-1].Resources().SysConfig().VnetPort)
	alias := AliasOf(vnetPort, vnicNum)
	return this.repHandlers[alias]
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
	for _, _handler := range this.trHandlers {
		_handler.Reset()
	}
	for _, _handler := range this.repHandlers {
		_handler.Reset()
	}
}

func (this *TestTopology) AllHandlers() []*TestServiceHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]*TestServiceHandler, 0)
	for _, h := range this.handlers {
		result = append(result, h)
	}
	return result
}

func (this *TestTopology) AllTrHandlers() []*TestServiceTransactionHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]*TestServiceTransactionHandler, 0)
	for _, h := range this.trHandlers {
		result = append(result, h)
	}
	return result
}

func (this *TestTopology) AllRepHandlers() []*TestServiceReplicationHandler {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]*TestServiceReplicationHandler, 0)
	for _, h := range this.repHandlers {
		result = append(result, h)
	}
	return result
}

func (this *TestTopology) AllVnics() []IVNic {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	result := make([]IVNic, 0)
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

func (this *TestTopology) ReActivateTestService(nic IVNic) {
	h, err := nic.Resources().Services().Activate(ServiceType, ServiceName, 0, nic.Resources(), nil,
		nic.Resources().SysConfig().LocalAlias)
	if err != nil {
		panic(err)
	}
	handler := h.(*TestServiceHandler)
	this.handlers[nic.Resources().SysConfig().LocalAlias] = handler
}
