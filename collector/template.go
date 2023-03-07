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

//go:build (darwin || linux || openbsd || netbsd) && !nomeminfo
// +build darwin linux openbsd netbsd
// +build !nomeminfo

package collector

import (
	"fmt"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	templateSubsystem = "memory"
)

type templateCollector struct {
	logger log.Logger
}

func init() {
	registerCollector("templateInfo", defaultEnabled, NewTemplateCollector)
}

// NewTemplateCollector returns a new Collector exposing memory stats.
func NewTemplateCollector(logger log.Logger) (Collector, error) {
	return &templateCollector{logger}, nil
}

// Update calls (*templateCollector).gettemplate to get the platform specific
// memory metrics.
func (t *templateCollector) Update(ch chan<- prometheus.Metric) error {
	var metricType prometheus.ValueType
	templateInfo, err := t.getTemplateInfo()
	if err != nil {
		return fmt.Errorf("couldn't get templateinfo: %w", err)
	}
	level.Debug(t.logger).Log("msg", "Set node_mem", "templateInfo", templateInfo)

	for k, v := range templateInfo {

		if strings.HasSuffix(k, "_total") {
			metricType = prometheus.CounterValue
		} else {
			metricType = prometheus.GaugeValue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, templateSubsystem, k),
				fmt.Sprintf("Template information field %s.", k),
				nil, nil,
			),
			metricType, v,
		)
	}
	return nil
}
