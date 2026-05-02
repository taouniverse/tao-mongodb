// Copyright 2021-2026 huija
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/taouniverse/tao"
)

func TestConfig(t *testing.T) {
	m := &Config{
		BaseMultiConfig: tao.BaseMultiConfig[InstanceConfig]{
			Instances: []tao.Instance[InstanceConfig]{
				{Name: "default", Cfg: InstanceConfig{}},
			},
		},
	}
	m.ValidSelf()

	instance, _ := m.GetInstanceByName("default")
	assert.Equal(t, "localhost", instance.Host)
	assert.Equal(t, 27017, instance.Port)
	assert.Equal(t, "tao", instance.User)

	t.Log(m.RunAfter())
	t.Log(m.ToTask())
}

func TestConfigDefaultInstanceName(t *testing.T) {
	c := &Config{}
	assert.Equal(t, "default", c.GetDefaultInstanceName())

	c.Default = "master"
	assert.Equal(t, "master", c.GetDefaultInstanceName())
}
