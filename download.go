// Copyright (c) 2020-2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ia

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const TimestampFormat = "20060102150405"

func PageURL(url, timestamp string) string {
	return "https://web.archive.org/web/" + timestamp + "id_/" + url
}

func DownloadFile(url, filename string) error {
	if _, err := os.Stat(filename); err == nil {
		// Skip existing
		return nil
	}

	resp, err := httpGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}

	if mod := resp.Header.Get("Last-Modified"); mod != "" {
		mt, err := time.Parse(time.RFC1123, mod)
		if err != nil {
			return err
		}
		return os.Chtimes(filename, mt, mt)
	}
	return nil
}

func httpGet(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("ia: http status %s", resp.Status)
	}
	return resp, nil
}
