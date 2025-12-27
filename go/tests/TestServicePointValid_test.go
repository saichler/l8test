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

package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/saichler/l8bus/go/overlay/health"
	"github.com/saichler/l8bus/go/overlay/protocol"
	"github.com/saichler/l8test/go/infra/t_resources"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	tear()
}

func TestServiceValid(t *testing.T) {

	t_resources.CreateTestModelInstance(1)

	for vnetNum := 1; vnetNum <= 3; vnetNum++ {
		for vnicNum := 1; vnicNum <= 4; vnicNum++ {
			nic := topo.VnicByVnetNum(vnetNum, vnicNum)
			hp := health.All(nic.Resources())
			if len(hp) != 15 {
				t_resources.Log.Fail(t, "Expected ", nic.Resources().SysConfig().LocalAlias,
					" to have 15 heath points, but it has ", len(hp))
				return
			}
		}
	}

	fmt.Println("Before Total - ", protocol.MsgLog.Total())
	time.Sleep(time.Second * 10)
	protocol.MsgLog.Print()
}
