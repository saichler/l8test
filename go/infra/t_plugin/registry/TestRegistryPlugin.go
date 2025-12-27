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
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
)

var Plugin ifs.IPlugin = &TestRegistryPlugin{}

type TestRegistryPlugin struct{}

func (this *TestRegistryPlugin) Install(vnic ifs.IVNic) error {
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
