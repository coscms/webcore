package s3manager_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/admpub/godotenv"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/s3manager/s3client"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func TestStat(t *testing.T) {
	return
	projectDir := filepath.Join(echo.Wd(), `../../../`)
	envFile := filepath.Join(projectDir, `.env`)
	err := godotenv.Overload(envFile)
	if err != nil {
		panic(err)
	}
	cfg := &dbschema.NgingCloudStorage{
		Key:      os.Getenv(`S3_KEY`),
		Secret:   os.Getenv(`S3_SECRET`),
		Secure:   `Y`,
		Region:   os.Getenv(`S3_REGION`),
		Bucket:   os.Getenv(`S3_BUCKET`),
		Endpoint: os.Getenv(`S3_ENDPOINT`),
		Baseurl:  os.Getenv(`S3_BASEURL`),
	}
	mgr := s3client.New(cfg, 1024000)
	exists, err := mgr.Exists(context.Background(), ``)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, false, exists)
}
