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

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

/*
/// REGEX C#
logRegexDict = new Dictionary<String, String>();
logRegexDict["RoundTripTime"] = @"Tx thread info: round trip time [(]ms[)] =\s+(\d*)";
logRegexDict["Bandwidth"] = @"Tx thread info: bw limit = (\d*)(, plateau = )*(\d+\.*\d*)*, avg tx = (\d+\.*\d*), avg rx = (\d+\.*\d*) [(](kbit|KBytes)/s[)]";
logRegexDict["Packets"] = @"Stat frms: R=(\d*)/(\d*)/(\d*)\s+T=(\d*)/(\d*)/(\d*)\s+[(]A/I/O[)] Loss=(\d+\.*\d*)%/(\d+\.*\d*)% [(]R/T[)]";
logRegexDict["Imaging"] = @"log [(](Soft|Tera2800)IPC[)]: (tbl \d* )?fps (\d+\.*\d*)( quality (\d*))?";
logRegexDict["Closed"] = @"Session closed remotely!";
*/

var (
	regexList = []*regexp.Regexp{
		regexp.MustCompile(`Tx thread info: round trip time [(]ms[)] =\s+(\d*), variance =\s+(\d*), rto =\s+(\d*), last =\s+(\d*), max =\s+(\d*)`),
		regexp.MustCompile(`Tx thread info: bw limit = (\d*), plateau = *(\d+\.*\d*)*, avg tx = (\d+\.*\d*), avg rx = (\d+\.*\d*) [(](?:kbit|KBytes)/s[)]`),
		regexp.MustCompile(`Stat frms: R=(\d*)/(\d*)/(\d*)\s+T=(\d*)/(\d*)/(\d*)\s+[(]A/I/O[)] Loss=(\d+\.*\d*)%/(\d+\.*\d*)% [(]R/T[)]`),
		regexp.MustCompile(`log [(]SoftIPC[)]: tbl (\d*) fps (\d+\.*\d*) quality (\d*)`),
		regexp.MustCompile(`log [(]SoftIPC[)]:  bits/pixel -\s+(\d+\.*\d*), bits/sec -\s+(\d+\.*\d*), MPix/sec -\s+(\d+\.*\d*)`),
		// regexp.MustCompile(`ubs-BW-decr: Decrease (loss) loss=0.001 current[kbit/s]=155.1188, active[kbit/s]=2555.1782 -> 534.9705, adjust factor=2.00%, floor[kbit/s]=104.0000`),
	}
	namesParam = [][]string{
		[]string{"round_trip_time", "rtt_variance", "rtt_rto", "rtt_last", "rtt_max"},
		[]string{"bw_limit", "bw_plateau", "avg_tx", "avg_rx"},
		[]string{"pkg_received_image", "pkg_received_audio", "pkg_received_others", "pkg_transf_image", "pkg_transf_audio", "pkg_transf_others", "pkg_loss_received", "pkg_loss_transf"},
		[]string{"tbl", "fps", "image_quality"},
		[]string{"bits_per_pixel", "bits_per_sec", "MPix_per_sec"},
	}
	regSessionClosed = regexp.MustCompile(`Session (closed) remotely!`)
)

func (p *pcoipCollector) getPcoipInfo() (map[string]float64, error) {
	file, err := os.Open(pcoipFilePath("server.1000.log"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parsePcoipInfo(file)
}

func parsePcoipInfo(r io.Reader) (map[string]float64, error) {
	var (
		pcoipInfo = map[string]float64{}
		scanner   = bufio.NewScanner(r)
	)

	first := true

	for scanner.Scan() {
		line := scanner.Text()

		if first {
			stringDate := line[:24]
			date, err := time.Parse("2006-01-02T15:04:05.000Z", stringDate)
			if err != nil {
				fmt.Printf("err: %s\n", err.Error())
			}
			pcoipInfo["session_duration"] = time.Now().Sub(date).Seconds()
			first = false
		}

		if len(regSessionClosed.FindStringSubmatch(line)) != 0 {
			return map[string]float64{}, nil
		}

		for i, regex := range regexList {

			match := regex.FindStringSubmatch(line)
			if len(match) == 0 {
				continue // if any regex find match
			}

			// iterate over names and link with each matched in vector match
			for j := range namesParam[i] {

				// prevent wrong match size
				if len(namesParam[i]) != len(match)-1 {
					fmt.Printf("The matches has not the same size of names dict - names size: :%d\n matchs:%d\n", len(namesParam[i]), len(match)-1)
					fmt.Printf("matchs: %v", match)
					continue
				}

				name := namesParam[i][j]
				value, err := strconv.ParseFloat(match[j+1], 64)
				if err != nil {
					fmt.Printf("err: %s\n", err.Error())
					continue
				}

				pcoipInfo[name] = value
			}
			break // if matched dont need test others regex
		}

	}
	return pcoipInfo, scanner.Err()
}
