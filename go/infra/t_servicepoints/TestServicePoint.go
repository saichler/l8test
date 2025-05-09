package t_servicepoints

import (
	"errors"
	. "github.com/saichler/l8test/go/infra/t_resources"
	. "github.com/saichler/l8srlz/go/serialize//object"
	"github.com/saichler/l8services/go/services/dcache"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"sync/atomic"
)

type TestServicePointBase struct {
	name         string
	postNumber   atomic.Int32
	putNumber    atomic.Int32
	patchNumber  atomic.Int32
	deleteNumber atomic.Int32
	getNumber    atomic.Int32
	failedNumber atomic.Int32
	errorMode    bool
}

const (
	ServiceName         = "Tests"
	ServicePointType    = "TestServicePointHandler"
	ServicePointTrType  = "TestServicePointTransactionHandler"
	ServicePointRepType = "TestServicePointReplicationHandler"
)

func (this *TestServicePointHandler) Activate(serviceName string, serviceArea uint16,
	r ifs.IResources, l ifs.IServiceCacheListener, args ...interface{}) error {
	this.name = args[0].(string)
	return nil
}

func (this *TestServicePointTransactionHandler) Activate(serviceName string, serviceArea uint16,
	r ifs.IResources, l ifs.IServiceCacheListener, args ...interface{}) error {
	this.name = args[0].(string)
	return nil
}

func (this *TestServicePointReplicationHandler) Activate(serviceName string, serviceArea uint16,
	r ifs.IResources, l ifs.IServiceCacheListener, args ...interface{}) error {
	this.name = args[0].(string)
	this.cache = dcache.NewDistributedCache(serviceName, serviceArea, "TestProto", r.SysConfig().LocalUuid, l, r)
	return nil
}
func (this *TestServicePointBase) DeActivate() error {
	return nil
}

func (this *TestServicePointBase) Post(pb ifs.IElements, resourcs ifs.IResources) ifs.IElements {
	Log.Debug("Post -", this.name, "- Test callback")
	this.postNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Post - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) Put(pb ifs.IElements, resourcs ifs.IResources) ifs.IElements {
	Log.Debug("Put -", this.name, "- Test callback")
	this.putNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Put - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) Patch(pb ifs.IElements, resourcs ifs.IResources) ifs.IElements {
	Log.Debug("Patch -", this.name, "- Test callback")
	this.patchNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Patch - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) Delete(pb ifs.IElements, resourcs ifs.IResources) ifs.IElements {
	Log.Debug("Delete -", this.name, "- Test callback")
	this.deleteNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Delete - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) GetCopy(pb ifs.IElements, resourcs ifs.IResources) ifs.IElements {
	Log.Debug("Get -", this.name, "- Test callback")
	this.getNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("GetCopy - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) Get(pb ifs.IElements, resourcs ifs.IResources) ifs.IElements {
	Log.Debug("Get -", this.name, "- Test callback")
	this.getNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Get - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) Failed(pb ifs.IElements, resourcs ifs.IResources, info ifs.IMessage) ifs.IElements {
	dest := "n/a"
	msg := "n/a"
	if info != nil {
		dest = info.Source()
		msg = info.FailMessage()
	}
	Log.Debug("Failed -", this.name, " to ", dest, "- Test callback")
	Log.Debug("Failed Reason is ", msg)
	this.failedNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Failed - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) EndPoint() string {
	return "/Tests"
}
func (this *TestServicePointBase) ServiceName() string {
	return ServiceName
}
func (this *TestServicePointBase) ServiceModel() ifs.IElements {
	return New(nil, &testtypes.TestProto{})
}

type TestServicePointHandler struct {
	TestServicePointBase
}

func (this *TestServicePointHandler) TransactionMethod() ifs.ITransactionMethod {
	return nil
}

type TestServicePointTransactionHandler struct {
	TestServicePointBase
}

func (this *TestServicePointTransactionHandler) TransactionMethod() ifs.ITransactionMethod {
	return this
}

func (this *TestServicePointTransactionHandler) Replication() bool {
	return false
}
func (this *TestServicePointTransactionHandler) ReplicationCount() int {
	return 0
}
func (this *TestServicePointTransactionHandler) KeyOf(elements ifs.IElements, resources ifs.IResources) string {
	return ""
}

type TestServicePointReplicationHandler struct {
	TestServicePointBase
	cache ifs.IDistributedCache
}

func (this *TestServicePointReplicationHandler) TransactionMethod() ifs.ITransactionMethod {
	return this
}

func (this *TestServicePointReplicationHandler) Replication() bool {
	return true
}
func (this *TestServicePointReplicationHandler) ReplicationCount() int {
	return 2
}
func (this *TestServicePointReplicationHandler) KeyOf(elements ifs.IElements, resources ifs.IResources) string {
	pb := elements.Element().(*testtypes.TestProto)
	return pb.MyString
}
