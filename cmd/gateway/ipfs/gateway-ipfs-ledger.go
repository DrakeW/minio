package ipfs

import (
	"context"

	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/namespace"
	badger "github.com/ipfs/go-ds-badger"
)

// ledgerStore keeps track of the mapping between
// s3 object to its CID stored on IPFS
type ledgerStore struct {
	ds datastore.Datastore
}

func newLedgerStore(path string) (*ledgerStore, error) {
	store, err := badger.NewDatastore(path, &badger.DefaultOptions)
	if err != nil {
		return nil, err
	}
	return &ledgerStore{
		ds: namespace.Wrap(store, datastore.NewKey("ledger")),
	}, nil
}

// TODO: implement this
func (l *ledgerStore) getAllBuckets(ctx context.Context) ([]ledgerStoreBucket, error) {
	return []ledgerStoreBucket{}, nil
}

// TODO: implement this
func (l *ledgerStore) getBucket(ctx context.Context, bucketName string) (ledgerStoreBucket, error) {
	return ledgerStoreBucket{}, nil
}

type ledgerStoreBucket struct {
	entries map[objectName]ipfsCID
}

type objectName string // TODO: some datastore key type?
type ipfsCID string    // TODO: some ipfs multiformat type?

func (b *ledgerStoreBucket) getObjectCID(name objectName) (ipfsCID, error) {
	return "xxx", nil
}
