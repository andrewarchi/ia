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
	return DownloadFileChecked(url, filename, nil)
}

func DownloadFileChecked(url, filename string, sha1Sum []byte) error {
	// Skip existing
	if _, err := os.Stat(filename); err == nil {
		if sha1Sum != nil {
			return ValidateFile(filename, nil, sha1Sum, nil)
		}
		// TODO check ETag, if it is a checksum
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
	var r io.Reader = resp.Body
	if sha1Sum != nil {
		r = NewReadValidator(r, url, nil, sha1Sum, nil)
	}
	if _, err := io.Copy(f, r); err != nil {
		return err
	}

	mod := resp.Header.Get("Last-Modified")
	if mod == "" {
		mod = resp.Header.Get("X-Archive-Orig-Last-Modified")
	}
	if mod != "" {
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
