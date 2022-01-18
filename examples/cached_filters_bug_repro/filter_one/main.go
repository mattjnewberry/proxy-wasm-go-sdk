// Copyright 2020-2021 Tetrate
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

package main

import (
	"math/rand"
	"time"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const tickMilliseconds uint32 = 10000

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &filterOne{}
}

type filterOne struct {
	types.DefaultPluginContext
	contextID uint32
	config []byte
	guid uint64
}

// Override types.DefaultPluginContext.
func (ctx *filterOne) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	configData, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("failed to get plugin config: %v", err)
		return types.OnPluginStartStatusFailed
	}
	if err := proxywasm.SetTickPeriodMilliSeconds(tickMilliseconds); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
	}
	rand.Seed(time.Now().UnixNano())
	ctx.config = configData
	ctx.guid =  rand.Uint64()
	t := time.Now().UnixNano()
	proxywasm.LogInfof("Filter 1 OnPluginStart: config: %s, guid: %d, time: %d", ctx.config, ctx.guid, t)
	return types.OnPluginStartStatusOK
}

// Override types.DefaultPluginContext.
func (*filterOne) NewHttpContext(uint32) types.HttpContext { return &types.DefaultHttpContext{} }

// Override types.DefaultPluginContext.
func (ctx *filterOne) OnTick() {
	t := time.Now().UnixNano()
	proxywasm.LogInfof("Filter 1 OnTick: config: %s, guid: %d, time: %d", ctx.config, ctx.guid, t)
}