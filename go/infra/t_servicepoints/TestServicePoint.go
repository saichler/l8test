package t_servicepoints

import (
	"errors"
	. "github.com/saichler/l8test/go/infra/t_resources"
	. "github.com/saichler/serializer/go/serialize/object"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/testtypes"
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

func (this *TestServicePointBase) Activate(serviceName string, serviceArea uint16,
	r common.IResources, l common.IServicePointCacheListener, args ...interface{}) error {
	this.name = args[0].(string)
	return nil
}

func (this *TestServicePointBase) DeActivate() error {
	return nil
}

func (this *TestServicePointBase) Post(pb common.IElements, resourcs common.IResources) common.IElements {
	Log.Debug("Post -", this.name, "- Test callback")
	this.postNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Post - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) Put(pb common.IElements, resourcs common.IResources) common.IElements {
	Log.Debug("Put -", this.name, "- Test callback")
	this.putNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Put - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) Patch(pb common.IElements, resourcs common.IResources) common.IElements {
	Log.Debug("Patch -", this.name, "- Test callback")
	this.patchNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Patch - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) Delete(pb common.IElements, resourcs common.IResources) common.IElements {
	Log.Debug("Delete -", this.name, "- Test callback")
	this.deleteNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Delete - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) GetCopy(pb common.IElements, resourcs common.IResources) common.IElements {
	Log.Debug("Get -", this.name, "- Test callback")
	this.getNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("GetCopy - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) Get(pb common.IElements, resourcs common.IResources) common.IElements {
	Log.Debug("Get -", this.name, "- Test callback")
	this.getNumber.Add(1)
	var err error
	if this.errorMode {
		err = errors.New("Get - TestServicePointBase Error")
	}
	return New(err, pb.Element())
}
func (this *TestServicePointBase) Failed(pb common.IElements, resourcs common.IResources, info common.IMessage) common.IElements {
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
func (this *TestServicePointBase) ServiceModel() common.IElements {
	return New(nil, &testtypes.TestProto{})
}

type TestServicePointHandler struct {
	TestServicePointBase
}

func (this *TestServicePointHandler) TransactionMethod() common.ITransactionMethod {
	return nil
}

type TestServicePointTransactionHandler struct {
	TestServicePointBase
}

func (this *TestServicePointTransactionHandler) TransactionMethod() common.ITransactionMethod {
	return this
}

func (this *TestServicePointTransactionHandler) Replication() bool {
	return false
}
func (this *TestServicePointTransactionHandler) ReplicationCount() int {
	return 0
}
func (this *TestServicePointTransactionHandler) KeyOf(elements common.IElements) string {
	return ""
}

type TestServicePointReplicationHandler struct {
	TestServicePointBase
}

func (this *TestServicePointReplicationHandler) TransactionMethod() common.ITransactionMethod {
	return this
}

func (this *TestServicePointReplicationHandler) Replication() bool {
	return true
}
func (this *TestServicePointReplicationHandler) ReplicationCount() int {
	return 2
}
func (this *TestServicePointReplicationHandler) KeyOf(elements common.IElements, resources common.IResources) string {
	pb := elements.Element().(*testtypes.TestProto)
	return pb.MyString
}
