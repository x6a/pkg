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

package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
	"x6a.dev/pkg/errors"
)

var EtcdClient *clientv3.Client
var Close = make(chan struct{})

func init() {
	go func() {
		<-Close
		EtcdClient.Close()
	}()
}

func NewClient(etcdEndpoint, etcdUsername, etcdPassword, etcdCACertB64 string) (*clientv3.Client, error) {
	var certs *x509.CertPool
	var etcdCfg clientv3.Config

	if len(etcdCACertB64) > 0 {
		blob, err := base64.URLEncoding.DecodeString(etcdCACertB64)
		if err != nil {
			return nil, errors.Wrapf(err, "[%v] function base64.URLEncoding.DecodeString(etcdCACertB64)", errors.Trace())
		}
		certs = x509.NewCertPool()
		if ok := certs.AppendCertsFromPEM(blob); !ok {
			return nil, errors.Wrapf(err, "[%v] function certs.AppendCertsFromPEM(blob)", errors.Trace())
		}

		etcdCfg = clientv3.Config{
			Endpoints: []string{
				etcdEndpoint,
			},
			Username:          etcdUsername,
			Password:          etcdPassword,
			DialTimeout:       10 * time.Second,
			DialKeepAliveTime: 10 * time.Second,
			TLS: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				MaxVersion:         tls.VersionTLS12,
				ClientCAs:          certs,
				InsecureSkipVerify: false,
			},
		}
	} else {
		etcdCfg = clientv3.Config{
			Endpoints: []string{
				etcdEndpoint,
			},
			Username:          etcdUsername,
			Password:          etcdPassword,
			DialTimeout:       10 * time.Second,
			DialKeepAliveTime: 10 * time.Second,
		}
	}

	c, err := clientv3.New(etcdCfg)
	for i := 0; err != nil && i < 10; i++ {
		fmt.Println("WARNING: unable to connect to etcd, retrying in 3s..")
		time.Sleep(3 * time.Second)
		c, err = clientv3.New(etcdCfg)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "[%v] function clientv3.New(etcdCfg): unable to etcd db", errors.Trace())
	}

	return c, nil
}

func etcdErrorHandler(err error) {
	switch err {
	case context.Canceled:
		log.Fatalf("ctx is canceled by another routine: %v", err)
	case context.DeadlineExceeded:
		log.Fatalf("ctx is attached with a deadline is exceeded: %v", err)
	case rpctypes.ErrEmptyKey:
		log.Fatalf("client-side error: %v", err)
	default:
		log.Fatalf("bad cluster endpoints, which are not etcd servers: %v", err)
	}
}

func Put(k, v string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := EtcdClient.Put(ctx, k, v)
	if err != nil {
		etcdErrorHandler(err)
		return errors.Wrapf(err, "[%v] function EtcdClient.Put(ctx, k, v)", errors.Trace())
	}

	return nil
}

func Get(k string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := EtcdClient.Get(ctx, k)
	if err != nil {
		etcdErrorHandler(err)
		return "", errors.Wrapf(err, "[%v] function EtcdClient.Get(ctx, k)", errors.Trace())
	}

	return string(resp.Kvs[0].Value), nil
}

func Delete(k string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := EtcdClient.Delete(ctx, k)
	if err != nil {
		etcdErrorHandler(err)
		return errors.Wrapf(err, "[%v] function EtcdClient.Delete(ctx, k)", errors.Trace())
	}

	return nil
}
