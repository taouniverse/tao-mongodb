// Copyright 2022 huija
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
	"github.com/taouniverse/tao"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// ConfigKey for this repo
const ConfigKey = "mongodb"

// Config implements tao.Config
// declare the configuration you want & define some default values
type Config struct {
	Host      string   `json:"host"`
	Port      int      `json:"port"`
	User      string   `json:"user"`
	Password  string   `json:"password"`
	RunAfters []string `json:"run_after,omitempty"`
}

var defaultMongodb = &Config{
	Host:      "localhost",
	Port:      27017,
	User:      "tao",
	Password:  "123456qwe",
	RunAfters: []string{},
}

// Name of Config
func (m *Config) Name() string {
	return ConfigKey
}

// ValidSelf with some default values
func (m *Config) ValidSelf() {
	if m.Host == "" {
		m.Host = defaultMongodb.Host
	}
	if m.Port == 0 {
		m.Port = defaultMongodb.Port
	}
	if m.User == "" {
		m.User = defaultMongodb.User
	}
	if m.Password == "" {
		m.Password = defaultMongodb.Password
	}
	if m.RunAfters == nil {
		m.RunAfters = defaultMongodb.RunAfters
	}
}

// ToTask transform itself to Task
func (m *Config) ToTask() tao.Task {
	return tao.NewTask(
		ConfigKey,
		func(ctx context.Context, param tao.Parameter) (tao.Parameter, error) {
			// non-block check
			select {
			case <-ctx.Done():
				return param, tao.NewError(tao.ContextCanceled, "%s: context has been canceled", ConfigKey)
			default:
			}
			// JOB code run after RunAfters, you can just do nothing here
			err := Client.Ping(ctx, readpref.Primary())
			return param, err
		},
		tao.SetClose(func() error {
			ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFunc()
			return Client.Disconnect(ctx)
		}))
}

// RunAfter defines pre task names
func (m *Config) RunAfter() []string {
	return m.RunAfters
}
