// SPDX-FileCopyrightText: 2024 Buoyant Inc.
// SPDX-License-Identifier: Apache-2.0
//
// Copyright 2022-2024 Buoyant Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.  You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package faces

import (
	"fmt"
	"net/http"
	"net/http/httputils"
	"time"

	"github.com/BuoyantIO/faces-demo/v2/pkg/utils"
)

type ColorServer struct {
	BaseServer
	color string
}

func NewColorServer(serverName string) *ColorServer {
	srv := &ColorServer{
		BaseServer: BaseServer{
			Name: serverName,
		},
	}

	srv.SetupFromEnvironment()
	// srv.SetUpdater(srv.updater)

	srv.RegisterNormal("/", srv.colorGetHandler)

	return srv
}

func (srv *ColorServer) SetupFromEnvironment() {
	srv.BaseServer.SetupFromEnvironment()

	colorName := utils.StringFromEnv("COLOR", "blue")
	srv.color = Colors.Lookup(colorName)

	fmt.Printf("%s %s: color %s => %s\n", time.Now().Format(time.RFC3339), srv.Name, colorName, srv.color)
}

func (srv *ColorServer) colorGetHandler(r *http.Request, rstat *BaseRequestStatus) *BaseServerResponse {
	res, err := httputils.DumpRequest(r, true)
	if err != nil {
		fmt.Printf("Error dumping request: %v\n", err)
	}
	fmt.Printf("%s %s: %s\n", time.Now().Format(time.RFC3339), srv.Name, res)

	// The only error we need to handle here is the internal rate limiter.
	if rstat.ratelimited {
		errstr := fmt.Sprintf("Rate limited (%.1f RPS > max %.1f RPS)", srv.CurrentRate(), srv.maxRate)

		return &BaseServerResponse{
			StatusCode: http.StatusTooManyRequests,
			Data: map[string]interface{}{
				"color":  Colors.Lookup(Defaults["color-ratelimit"]),
				"rate":   fmt.Sprintf("%.1f RPS", srv.CurrentRate()),
				"errors": []string{errstr},
			},
		}
	}

	return &BaseServerResponse{
		StatusCode: http.StatusOK,
		Data: map[string]interface{}{
			"color":  srv.color,
			"rate":   fmt.Sprintf("%.1f RPS", srv.CurrentRate()),
			"errors": []string{},
		},
	}
}
