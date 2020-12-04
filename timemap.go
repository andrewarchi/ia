// Copyright (c) 2020 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ia

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type TimeMap struct {
	URLKey     string
	Timestamp  string
	Original   string
	MIMEType   string
	StatusCode string
	Digest     string
	Redirect   string
	RobotFlags string
	Length     string
	Offset     string
	Filename   string
}

func GetTimeMap(url string) ([]TimeMap, error) {
	res, err := http.Get("https://web.archive.org/web/timemap/json/" + url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var lines [][]string
	if err := json.NewDecoder(res.Body).Decode(&lines); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errors.New("ia: empty response")
	}
	timemap := make([]TimeMap, len(lines)-1)
	for i, line := range lines[1:] {
		if len(line) != 11 {
			return nil, fmt.Errorf("ia: time map entry %d not length 11", i)
		}
		timemap[i] = TimeMap{line[0], line[1], line[2], line[3], line[4],
			line[5], line[6], line[7], line[8], line[9], line[10]}
	}
	return timemap, nil
}
