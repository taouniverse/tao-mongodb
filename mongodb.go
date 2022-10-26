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
	"github.com/taouniverse/tao"
	// Load the required dependencies.
	// An error occurs when there was no package in the root directory.
	_ "go.mongodb.org/mongo-driver/mongo"
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

// TODO setup unit with the global config 'M'
// execute when init tao universe
func setup() error {
	return nil
}