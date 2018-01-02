package cluster

type Config struct {
	SignerCount  int
	SyncerCount  int
	NodeDataPath string
	KeyStorePath string
	GethPath     string
}
