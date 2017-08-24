package bftkvst

import (
	"bytes"
	"encoding/binary"

	"github.com/google/keytransparency/core/monitor/storage"
	localst "github.com/google/keytransparency/impl/monitor/storage"
	ktpb "github.com/google/keytransparency/core/proto/keytransparency_v1_types"
	"github.com/google/trillian"

	"github.com/yahoo/bftkv/api"
	"github.com/yahoo/bftkv/protocol"
)

type bftkvStorage struct {
	storage *localst.Storage
	client *protocol.Client
}

func New(path string) storage.Storage {
	storage := localst.New()
	client, err := api.OpenClient(path)
	if err != nil {
		client = nil
	}
	return &bftkvStorage{storage, client}
}

func (s *bftkvStorage) Set(epoch int64,
	seenNanos int64,
	smr *trillian.SignedMapRoot,
	response *ktpb.GetMutationsResponse,
	errorList []error) error {

	err := s.storage.Set(epoch, seenNanos, smr, response, errorList)
	if err != nil {
		return err
	}
	if s.client != nil {
		var buf bytes.Buffer
		binary.Write(&buf, binary.BigEndian, epoch)
		key := buf.Bytes()
		if _, err := s.client.Read(key); err == nil {
			return storage.ErrAlreadyStored
		}
		err = s.client.Write(key, []byte(response.String()))
	}
	return err
}

func (s *bftkvStorage) Get(epoch int64) (*storage.MonitoringResult, error) {
	return s.storage.Get(epoch)
}

func (s *bftkvStorage) LatestEpoch() int64 {
	return s.storage.LatestEpoch()
}
