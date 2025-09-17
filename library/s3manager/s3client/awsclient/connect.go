/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package awsclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/coscms/webcore/dbschema"
)

func Connect(c context.Context, m *dbschema.NgingCloudStorage, bucketName string) (client *AWSClient, err error) {
	cfg, err := NewConfig(c, m)
	if err != nil {
		return nil, err
	}
	return &AWSClient{Client: s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.EndpointOptions.DisableHTTPS = m.Secure != `Y`
	}), bucketName: bucketName}, nil
}

func NewConfig(c context.Context, m *dbschema.NgingCloudStorage) (aws.Config, error) {
	opts := []func(*config.LoadOptions) error{ //[]config.LoadOptionsFunc{
		config.WithRegion(m.Region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     m.Key,
				SecretAccessKey: m.Secret,
				//Source: "provider",
			},
		}),
	}
	if len(m.Endpoint) > 0 {
		scheme := `https://`
		if m.Secure != `Y` {
			scheme = `http://`
		}
		opts = append(opts, config.WithBaseEndpoint(scheme+m.Endpoint))
	}
	return config.LoadDefaultConfig(c, opts...)
}
