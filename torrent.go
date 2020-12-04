// Copyright (c) 2020 Andrew Archibald
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
		url := fmt.Sprintf("https://archive.org/download/%s/%s_archive.torrent", id, id)
		fmt.Printf("(%d/%d) Adding %s\n", i, len(ids), id)
		filename := filepath.Join(dir, path.Base(url))
		if err := DownloadFile(url, filename); err != nil {
			return err
		}
		t, err := c.AddTorrentFromFile(filename)
		if err != nil {
			return err
		}
		t.DownloadAll()
	}
	fmt.Println(c.WaitAll())
	return nil
}
