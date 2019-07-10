package elasticsearch

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/x6a/pkg/errors"
)

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
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: 10*time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS12,
				RootCAs:    certs,
			},
		},
	}

	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, errors.Wrapf(err, "[%v] function elasticsearch.NewClient(esCfg)", errors.Trace())
	}

	return es, nil
}

func ESGet(es *elasticsearch.Client, index, objID string, srcIncludes ...string) ([]byte, error) {
	resp, err := es.Get(
		index,
		objID,
		es.Get.WithPretty(),
		es.Get.WithSourceIncludes(srcIncludes...),
	)
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
