// Copyright (c) 2020-2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ia

import (
	"fmt"

	"github.com/andrewarchi/browser/jsonutil"
)

// Search queries the Internet Archive for the identifiers of all
// matching items.
func Search(query string) ([]string, error) {
	url := "https://archive.org/services/search/v1/scrape?q=" + query + "&count=10000"
	resp, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type Scrape struct {
		Items []struct {
			Identifier string `json:"identifier"`
		} `json:"items"`
		Count int `json:"count"`
		Total int `json:"total"`
		// TODO fields for error response
	}
	var items Scrape
	if err := jsonutil.Decode(resp.Body, &items); err != nil {
		return nil, err
	}

	// TODO handle paging
	if items.Count != items.Total {
		return nil, fmt.Errorf("ia: queried %d of %d items", items.Count, items.Total)
	}

	ids := make([]string, len(items.Items))
	for i, item := range items.Items {
		ids[i] = item.Identifier
	}
	return ids, nil
}
