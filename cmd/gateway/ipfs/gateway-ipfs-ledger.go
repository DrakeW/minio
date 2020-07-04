package ipfs

import (
	"context"

	"github.com/ipfs/go-datastore"
)

// ledgerStore keeps track of the mapping between
// s3 object to its CID stored on IPFS
type ledgerStore struct {
	ds      *datastore.Datastore
	buckets map[string]ledgerStoreBucket // a local cache of buckets
}

func newLedgerStore() *ledgerStore {}

func (l *ledgerStore) getAllBuckets(ctx context.Context) ([]ledgerStoreBucket, error) {
	return []ledgerStoreBucket{}, nil
}

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
