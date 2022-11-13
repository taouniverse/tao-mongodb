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
	"fmt"
	"github.com/taouniverse/tao"
	"time"

	// Load the required dependencies.
	// An error occurs when there was no package in the root directory.
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
import _ "github.com/taouniverse/tao-mongodb"
*/

// M config of mongodb
var M = new(Config)

func init() {
	err := tao.Register(ConfigKey, M, setup)
	if err != nil {
		panic(err.Error())
	}
}

// Client of mongodb
var Client *mongo.Client

// setup unit with the global config 'M'
// execute when init tao universe
func setup() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	Client, err = mongo.Connect(ctx,
		options.Client().
			ApplyURI(fmt.Sprintf("mongodb://%s:%d", M.Host, M.Port)).
			SetAuth(options.Credential{
				Username: M.User,
				Password: M.Password,
			}),
	)
	return err
}
