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
	"context"
	"time"

	"github.com/taouniverse/tao"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ConfigKey for this repo
const ConfigKey = "mongodb"

// InstanceConfig 单实例配置
type InstanceConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
}

// Config 总配置，实现 tao.MultiConfig 接口
type Config struct {
	tao.BaseMultiConfig[InstanceConfig]
	RunAfters []string `json:"run_after,omitempty" yaml:"run_after,omitempty"`
}

var defaultInstance = &InstanceConfig{
	Host:     "localhost",
	Port:     27017,
	User:     "tao",
	Password: "123456qwe",
}

// Name of Config
func (m *Config) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (m *Config) ValidSelf() {
	for name, instance := range m.Instances {
		if instance.Host == "" {
			instance.Host = defaultInstance.Host
		}
		if instance.Port == 0 {
			instance.Port = defaultInstance.Port
		}
		if instance.User == "" {
			instance.User = defaultInstance.User
		}
		if instance.Password == "" {
			instance.Password = defaultInstance.Password
		}
		m.Instances[name] = instance
	}
	if m.RunAfters == nil {
		m.RunAfters = []string{}
	}
}

// ToTask transform itself to Task
func (m *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			for name := range m.Instances {
				client, err := Factory.Get(name)
				if err != nil {
					return param, err
				}
				err = client.Ping(ctx, readpref.Primary())
				if err != nil {
					return param, err
				}
			}
			return param, nil
		},
		tao.SetClose(func() error {
			_, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFunc()
			return Factory.CloseAll()
		}))
}

// RunAfter defines pre task names
func (m *Config) RunAfter() []string {
	return m.RunAfters
}
