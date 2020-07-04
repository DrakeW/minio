package ipfs

import (
	"io"

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
	apiAddr ma.Multiaddr
}

func (g *Ipfs) Name() string {
	return ipfsBackend
}

func (g *Ipfs) NewGatewayLayer(creds auth.Credentials) (minio.ObjectLayer, error) {
	return nil, nil
}

func (g *Ipfs) Production() bool {
	return false
}

// IpfsObject implements ObjectLayer
type IpfsObject struct {
	ledger *ledgerStore // ledger store handles the the mapping between buckets/objects to CIDs
	client *ipfsClient  // client handles the interaction with IPFS
}

// TODO: figure out its signature
func NewIpfsObject() *IpfsObject {}

// ipfsClient handles the read/write operation on IPFS
type ipfsClient struct {
	c client.Client
}

func newIpfsClient(cfg *client.Config) *ipfsClient {
	client := client.NewDefaultClient(cfg)
	return &ipfsClient{
		c: client,
	}
}

func (i *ipfsClient) addFile(f io.Reader) (ipfsCID, error) {}

func (i *ipfsClient) getFile(cid ipfsCID) ([]bytes, error) {}
