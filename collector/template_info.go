// Copyright 2015 The Prometheus Authors
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

//go:build !nomeminfo
// +build !nomeminfo

package collector

func (c *templateCollector) getTemplateInfo() (map[string]float64, error) {

	// to manipulate files see example with meminfo in meminfo_linux.go
	var (
		templateInfo = map[string]float64{}
	)

	templateInfo["tempate_info_param_total"] = 4.0
	templateInfo["tempate_info_param_1"] = 1.0
	templateInfo["tempate_info_param_2"] = 2.0
	templateInfo["tempate_info_param_3"] = 3.0
	templateInfo["tempate_info_param_4"] = 4.0

	return templateInfo, nil
}
