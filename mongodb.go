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
	"fmt"
	"time"

	"github.com/taouniverse/tao"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
import _ "github.com/taouniverse/tao-mongodb"
*/

// M is the global config instance for tao-mongodb
var M = &Config{}

// Factory is the global factory instance for managing mongo.Client
var Factory *tao.BaseFactory[*mongo.Client]

func init() {
	var err error
	Factory, err = tao.Register(ConfigKey, M, NewMongoDB)
	if err != nil {
		panic(err.Error())
	}
}

// NewMongoDB creates a new MongoDB client for factory pattern
func NewMongoDB(name string, config InstanceConfig) (*mongo.Client, func() error, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx,
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port)).
			SetAuth(options.Credential{
				Username: config.User,
				Password: config.Password,
			}),
	)
	if err != nil {
		return nil, nil, tao.NewErrorWrapped("mongodb: fail to create mongo client", err)
	}

	closer := func() error {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()
		return client.Disconnect(ctx)
	}

	return client, closer, nil
}

// Client returns the default mongo client instance
func Client() (*mongo.Client, error) {
	return Factory.Get(M.GetDefaultInstanceName())
}

// GetClient returns the mongo client instance by name
func GetClient(name string) (*mongo.Client, error) {
	return Factory.Get(name)
}
