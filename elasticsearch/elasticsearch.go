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

package elasticsearch

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"

	"x6a.dev/pkg/errors"
)

var ESClient *elasticsearch.Client
var Close = make(chan struct{})

func init() {
	go func() {
		<-Close
	}()
}

func ESConnect(esURL, esUsername, esPassword, esCACertB64 string) (*elasticsearch.Client, error) {
	var certs *x509.CertPool

	if len(esCACertB64) > 0 {
		blob, err := base64.URLEncoding.DecodeString(esCACertB64)
		if err != nil {
			return nil, errors.Wrapf(err, "[%v] function base64.URLEncoding.DecodeString(esCACertB64)", errors.Trace())
		}
		certs = x509.NewCertPool()
		if ok := certs.AppendCertsFromPEM(blob); !ok {
			return nil, errors.Wrapf(err, "[%v] function certs.AppendCertsFromPEM(blob)", errors.Trace())
		}
	}

	esCfg := elasticsearch.Config{
		Addresses: []string{
			esURL,
		},
		Username: esUsername,
		Password: esPassword,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: 10 * time.Second,
			DialContext:           (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS12,
				RootCAs:    certs,
			},
		},
	}

	es, err := elasticsearch.NewClient(esCfg)
	for i := 0; err != nil && i < 10; i++ {
		log.Println("WARNING: unable to connect to elasticsearch db, retrying in 3s..")
		time.Sleep(3 * time.Second)
		es, err = elasticsearch.NewClient(esCfg)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "[%v] function elasticsearch.NewClient(esCfg): unable to connect to elasticsearch db", errors.Trace())
	}

	return es, nil
}

func ESGet(es *elasticsearch.Client, index, objID string, srcIncludes ...string) ([]byte, error) {
	var err error
	var resp *esapi.Response

	if len(srcIncludes) > 0 {
		resp, err = es.Get(
			index,
			objID,
			es.Get.WithPretty(),
			es.Get.WithSourceIncludes(srcIncludes...),
		)
	} else {
		resp, err = es.Get(
			index,
			objID,
			es.Get.WithPretty(),
		)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "[%v] function es.Get()", errors.Trace())
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, errors.Errorf("[%s] Error getting document ID=%s", resp.Status(), objID)
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
	if err != nil {
		return nil, errors.Wrapf(err, "[%v] function ioutil.ReadAll()", errors.Trace())
	}

	return body, nil
}

func ESIndex(es *elasticsearch.Client, index, objID string, obj []byte) error {
	resp, err := es.Index(
		index,
		bytes.NewReader(obj),
		es.Index.WithDocumentID(objID),
		es.Index.WithPretty(),
	)
	if err != nil {
		return errors.Wrapf(err, "[%v] function es.Index()", errors.Trace())
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return errors.Errorf("[%s] Error indexing document ID=%s", resp.Status(), objID)
	}

	return nil
}
