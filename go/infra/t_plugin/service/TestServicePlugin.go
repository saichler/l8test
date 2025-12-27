/*
 * © 2025 Sharon Aicler (saichler@gmail.com)
 *
 * Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"

	"github.com/saichler/l8test/go/infra/t_service"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
)

var Plugin ifs.IPlugin = &TestServicePlugin{}

type TestServicePlugin struct {
}

func (this TestServicePlugin) InstallRegistry(vnic ifs.IVNic) error {
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

func (this *TestServicePlugin) Install(vnic ifs.IVNic) error {
	this.InstallRegistry(vnic)

	vnic.Resources().Logger().Info("#2 Registering Test Services on ", vnic.Resources().SysConfig().LocalAlias)
	vnic.Resources().Services().RegisterServiceHandlerType(&t_service.TestServiceHandler{})
	vnic.Resources().Services().RegisterServiceHandlerType(&t_service.TestServiceTransactionHandler{})
	vnic.Resources().Services().RegisterServiceHandlerType(&t_service.TestServiceReplicationHandler{})

	sla := ifs.NewServiceLevelAgreement(&t_service.TestServiceHandler{}, t_service.ServiceName, 0, false, nil)
	sla.SetArgs("plugin")
	_, err := vnic.Resources().Services().Activate(sla, vnic)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return err
	}

	sla = ifs.NewServiceLevelAgreement(&t_service.TestServiceTransactionHandler{}, t_service.ServiceName, 1, true, nil)
	sla.SetArgs("plugin")
	_, err = vnic.Resources().Services().Activate(sla, vnic)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return err
	}

	sla = ifs.NewServiceLevelAgreement(&t_service.TestServiceReplicationHandler{}, t_service.ServiceName, 2, true, nil)
	sla.SetArgs("plugin")
	_, err = vnic.Resources().Services().Activate(sla, vnic)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return err
	}
	return nil
}
