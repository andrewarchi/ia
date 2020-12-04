// Copyright (c) 2020 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ia

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Search queries the Internet Archive for the identifiers of all
// matching items.
func Search(query string) ([]string, error) {
	url := fmt.Sprintf("https://archive.org/services/search/v1/scrape?q=%s&count=10000", query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ia: http status %s", resp.Status)
	}
	defer resp.Body.Close()

	type iaItem struct {
		Identifier string `json:"identifier"`
	}
	var items struct {
		Items []iaItem `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	ids := make([]string, len(items.Items))
	for i, item := range items.Items {
		ids[i] = item.Identifier
	}
	return ids, nil
}
