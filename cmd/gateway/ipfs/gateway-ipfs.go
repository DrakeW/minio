package ipfs

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/ipfs/ipfs-cluster/api/rest/client"
	"github.com/minio/cli"
	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/pkg/auth"
	ma "github.com/multiformats/go-multiaddr"
)

const (
	ipfsBackend = "ipfs"
)

func init() {
	const ipfsGatewayTemplate = `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} {{if .VisibleFlags}}[FLAGS]{{end}} [API_ADDR]
{{if .VisibleFlags}}
FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}
API_ADDR:
  IPFS Cluster API Address. If none the gateway will be communicating with the public IPFS network.

EXAMPLES:
  1. Start minio gateway server for IPFS on custom IPFS Cluster REST API endpoint.
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_ACCESS_KEY{{.AssignmentOperator}}ipfsclusterusername
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_SECRET_KEY{{.AssignmentOperator}}ipfsclusterpassword
     {{.Prompt}} {{.HelpName}} /p2p/mycustomp2paddress OR myhost:myport

  2. Start minio gateway server for IPFS on custom IPFS Cluster REST API endpoint with edge caching enabled.
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_ACCESS_KEY{{.AssignmentOperator}}ipfsclusterusername
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_SECRET_KEY{{.AssignmentOperator}}ipfsclusterpassword
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_DRIVES{{.AssignmentOperator}}"/mnt/drive1,/mnt/drive2,/mnt/drive3,/mnt/drive4"
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_EXCLUDE{{.AssignmentOperator}}"bucket1/*,*.png"
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_QUOTA{{.AssignmentOperator}}90
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_AFTER{{.AssignmentOperator}}3
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_WATERMARK_LOW{{.AssignmentOperator}}75
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_WATERMARK_HIGH{{.AssignmentOperator}}85
     {{.Prompt}} {{.EnvVarSetCommand}} MINIO_AZURE_CHUNK_SIZE_MB {{.AssignmentOperator}}25
     {{.Prompt}} {{.HelpName}} /p2p/mycustomp2paddress OR myhost:myport
`
	minio.RegisterGatewayCommand(cli.Command{
		Name:               ipfsBackend,
		Usage:              "IPFS Storage",
		Action:             ipfsGatewayMain,
		CustomHelpTemplate: ipfsGatewayTemplate,
		HideHelpCommand:    true,
	})
}

func ipfsGatewayMain(ctx *cli.Context) {
}

// Ipfs implements Gateway
type Ipfs struct {
	apiAddr string
}

// Name implements Name method of Gateway interface
func (g *Ipfs) Name() string {
	return ipfsBackend
}

// NewGatewayLayer implements NewGatewayLayer method of Gateway interface
// TODO: fill this method in
func (g *Ipfs) NewGatewayLayer(creds auth.Credentials) (minio.ObjectLayer, error) {
	username := creds.AccessKey
	password := creds.SecretKey
	endpoint := g.apiAddr
	// TODO: path should be taken as input
	object, err := NewIpfsObjectLayer(username, password, endpoint, "~/.test-ipfs-ds")
	if err != nil {
		return nil, err
	}
	return object, nil
}

// Production implements Production method of Gateway interface
func (g *Ipfs) Production() bool {
	return false
}

// IpfsObjectLayer implements ObjectLayer
type IpfsObjectLayer struct {
	ledger *ledgerStore // ledger store handles the the mapping between buckets/objects to CIDs
	client *ipfsClient  // client handles the interaction with IPFS
}

// NewIpfsObjectLayer initializes an IpfsObjectLayer
func NewIpfsObjectLayer(
	username, password, endpoint,
	dPath string,
) (*IpfsObjectLayer, error) {
	ipfsClient, err := newIpfsClient(username, password, endpoint)
	if err != nil {
		return nil, err
	}
	ledger, err := newLedgerStore(dPath)
	if err != nil {
		return nil, err
	}
	return &IpfsObjectLayer{
		ledger: ledger,
		client: ipfsClient,
	}, nil
}

// START - implements the ObjectLayer interface

func (obj *IpfsObjectLayer) GetBucketInfo(ctx context.Context, bucket string) (bucketInfo minio.BucketInfo, err error) {
}

func (obj *IpfsObjectLayer) ListBuckets(ctx context.Context) (buckets []minio.BucketInfo, err error) {}

func (obj *IpfsObjectLayer) DeleteBucket(ctx context.Context, bucket string, forceDelete bool) error {}

func (obj *IpfsObjectLayer) ListObjects(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
}
func (obj *IpfsObjectLayer) ListObjectsV2(ctx context.Context, bucket, prefix, continuationToken, delimiter string, maxKeys int, fetchOwner bool, startAfter string) (result minio.ListObjectsV2Info, err error) {
}
func (obj *IpfsObjectLayer) GetObjectNInfo(ctx context.Context, bucket, object string, rs *minio.HTTPRangeSpec, h http.Header, lockType minio.LockType, opts minio.ObjectOptions) (reader *minio.GetObjectReader, err error) {
}
func (obj *IpfsObjectLayer) GetObject(ctx context.Context, bucket, object string, startOffset int64, length int64, writer io.Writer, etag string, opts minio.ObjectOptions) (err error) {
}
func (obj *IpfsObjectLayer) GetObjectInfo(ctx context.Context, bucket, object string, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
}
func (obj *IpfsObjectLayer) PutObject(ctx context.Context, bucket, object string, data *minio.PutObjReader, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
}
func (obj *IpfsObjectLayer) CopyObject(ctx context.Context, srcBucket, srcObject, destBucket, destObject string, srcInfo minio.ObjectInfo, srcOpts, dstOpts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
}
func (obj *IpfsObjectLayer) DeleteObject(ctx context.Context, bucket, object string, opts minio.ObjectOptions) (minio.ObjectInfo, error) {
}
func (obj *IpfsObjectLayer) DeleteObjects(ctx context.Context, bucket string, objects []minio.ObjectToDelete, opts minio.ObjectOptions) ([]minio.DeletedObject, []error) {
}

// END - implements the ObjectLayer interface

// ipfsClient handles the read/write operation on IPFS
type ipfsClient struct {
	c client.Client
}

func newIpfsClient(username, password, endpoint string) (*ipfsClient, error) {
	// build client config
	cfg := &client.Config{}
	if isEndpointP2p(endpoint) {
		apiAddr, err := ma.NewMultiaddr(endpoint)
		if err != nil {
			return nil, err
		}
		cfg.APIAddr = apiAddr
	} else {
		hostAndPort := strings.Split(endpoint, ":")
		cfg.Host = hostAndPort[0]
		cfg.Port = hostAndPort[1]
	}
	cfg.Username = username
	cfg.Password = password
	cfg.LogLevel = "DEBUG" // TODO: temporary value here, make it a command line argument in the future

	client, err := client.NewDefaultClient(cfg)
	if err != nil {
		return &ipfsClient{
			c: client,
		}, nil
	}
	return nil, err
}

func isEndpointP2p(endpoint string) bool {
	if strings.HasPrefix("/ipfs", endpoint) ||
		strings.HasPrefix("/p2p", endpoint) ||
		strings.HasPrefix("/dnsaddr", endpoint) {
		return true
	}
	return false
}

// TODO: fill this in
func (i *ipfsClient) addFile(f io.Reader) (ipfsCID, error) {
	return "", nil
}

// TODO: fill this in
func (i *ipfsClient) getFile(cid ipfsCID) ([]byte, error) {
	return []byte{}, nil
}
