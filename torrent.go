// Copyright (c) 2020-2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ia

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/anacrolix/torrent"
)

// DownloadTorrents downloads the named Internet Archive items via
// torrent.
func DownloadTorrents(ids []string, dir string) error {
	conf := torrent.NewDefaultClientConfig()
	conf.DataDir = dir
	c, err := torrent.NewClient(conf)
	if err != nil {
		return err
	}

	for i, id := range ids {
		url := "https://archive.org/download/" + id + "/" + id + "_archive.torrent"
		fmt.Printf("(%d/%d) Adding %s\n", i+1, len(ids), id)
		filename := filepath.Join(dir, path.Base(url))
		if err := DownloadFile(url, filename); err != nil {
			return err
		}

		t, err := c.AddTorrentFromFile(filename)
		if err != nil {
			return err
		}
		t.DownloadAll()
		if i%15 == 14 {
			c.WaitAll()
		}
	}
	c.WaitAll()
	return nil
}
