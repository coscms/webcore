package cloudbackup

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webcore/model"
	smb "github.com/hirochachacha/go-smb2"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func init() {
	Register(model.StorageEngineSMB, newStorageSMB, smbForms, `SMB`)
}

func newStorageSMB(ctx echo.Context, cfg dbschema.NgingCloudBackup) (Storager, error) {
	if len(cfg.StorageConfig) == 0 {
		return nil, ErrEmptyConfig
	}
	conf := echo.H{}
	err := json.Unmarshal([]byte(cfg.StorageConfig), &conf)
	if err != nil {
		return nil, err
	}
	password := common.Crypto().Decode(conf.String(`password`))
	return NewStorageSMB(conf.String(`addr`), conf.String(`username`), password, conf.String(`sharename`)), nil
}

var smbForms = []Form{
	{Type: `text`, Label: `主机地址`, Name: `storageConfig.addr`, Required: true, Placeholder: `<IP或域名>:<端口>`},
	{Type: `text`, Label: `用户名`, Name: `storageConfig.username`, Required: true},
	{Type: `password`, Label: `密码`, Name: `storageConfig.password`, Required: true},
	{Type: `text`, Label: `共享名称`, Name: `storageConfig.sharename`, Required: true},
}

func NewStorageSMB(addr, username, password, sharename string) Storager {
	return &StorageSMB{addr: addr, username: username, password: password, sharename: sharename}
}

type StorageSMB struct {
	addr      string // host:port
	username  string
	password  string
	sharename string
	conn      net.Conn
	session   *smb.Session
	share     *smb.Share
	prog      notice.Progressor
}

func (s *StorageSMB) Connect() (err error) {
	if !strings.Contains(s.addr, `:`) {
		s.addr += `:445`
	}
	s.conn, err = net.Dial("tcp", s.addr)
	if err != nil {
		return
	}

	d := &smb.Dialer{
		Initiator: &smb.NTLMInitiator{
			User:     s.username,
			Password: s.password,
		},
	}

	s.session, err = d.Dial(s.conn)
	if err != nil {
		s.Close()
		return
	}

	s.share, err = s.session.Mount(s.sharename)
	if err != nil {
		s.Close()
		return
	}
	return
}

func (s *StorageSMB) Put(ctx context.Context, reader io.Reader, ppath string, size int64) (err error) {
	s.share.MkdirAll(path.Dir(ppath), os.ModePerm)
	var fp *smb.File
	fp, err = s.share.Create(ppath)
	if err != nil {
		return
	}
	defer fp.Close()
	_, err = io.Copy(fp, reader)
	return
}

func (s *StorageSMB) Download(ctx context.Context, ppath string, w io.Writer) error {
	resp, err := s.share.Open(ppath)
	if err != nil {
		return err
	}
	defer resp.Close()
	if s.prog != nil {
		stat, err := resp.Stat()
		if err != nil {
			return err
		}
		s.prog.Add(stat.Size())
		w = s.prog.ProxyWriter(w)
		defer s.prog.Reset()
	}
	_, err = io.Copy(w, resp)
	return err
}

func (s *StorageSMB) SetProgressor(prog notice.Progressor) {
	s.prog = prog
}

func (s *StorageSMB) Restore(ctx context.Context, ppath string, destpath string, callback func(from, to string)) error {
	resp, err := s.share.Open(ppath)
	if err != nil {
		return err
	}
	stat, err := resp.Stat()
	if err != nil {
		resp.Close()
		return err
	}
	if !stat.IsDir() {
		resp.Close()
		if callback != nil {
			callback(ppath, destpath)
		}
		return DownloadFile(s, ctx, ppath, destpath)
	}
	dirs, err := resp.Readdir(-1)
	resp.Close()
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		spath := path.Join(ppath, dir.Name())
		dest := filepath.Join(destpath, dir.Name())
		if dir.IsDir() {
			err = com.MkdirAll(dest, os.ModePerm)
			if err == nil {
				err = s.Restore(ctx, spath, dest, callback)
			}
		} else {
			if callback != nil {
				callback(spath, dest)
			}
			err = DownloadFile(s, ctx, spath, dest)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StorageSMB) RemoveDir(ctx context.Context, ppath string) error {
	return s.share.RemoveAll(ppath)
}

func (s *StorageSMB) Remove(ctx context.Context, ppath string) error {
	return s.share.Remove(ppath)
}

func (s *StorageSMB) Close() (err error) {
	if s.share != nil {
		err = s.share.Umount()
	}
	if s.session != nil {
		err = s.session.Logoff()
	}
	if s.conn != nil {
		err = s.conn.Close()
	}
	return
}
