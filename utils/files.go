// Copyright (C) 2019 <x6a@7n.io>
//
// pkg is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// pkg is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with pkg. If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"x6a.dev/pkg/errors"
)

// FileToB64 read and convert a file to base64
func FileToB64(file string) (string, error) {
	var blob []byte

	if _, err := os.Stat(file); err == nil {
		blob, err = ioutil.ReadFile(file)
		if err != nil {
			return "", errors.Wrapf(err, "[%v] function ioutil.ReadFile(file)", errors.Trace())
		}
	} else if os.IsNotExist(err) {
		fmt.Printf("file %v not found", file)
		return "", errors.Wrapf(err, "[%v] file %v not found", errors.Trace(), file)
	} else {
		return "", errors.Wrapf(err, "[%v] file stat error", errors.Trace())
	}

	return base64.URLEncoding.EncodeToString(blob), nil
}
