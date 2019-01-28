package p2p

import (
	maddr "github.com/multiformats/go-multiaddr"
	log "github.com/sirupsen/logrus"
)

var defaultBootstrapAddrStrings = []string{
	// "/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
	// "/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	// "/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
	// "/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
	// "/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
	"/ip4/192.168.1.67/tcp/4001/ipfs/QmRLEhyR23LxcWntQXHR9LGUowZohL7mAeVQgxZTw7TD1d", // Local node
}

func getBootstrapAddrs() []maddr.Multiaddr {
	bootstrapAddrs := []maddr.Multiaddr{}
	for _, bootstrapAddr := range defaultBootstrapAddrStrings {
		addr, err := maddr.NewMultiaddr(bootstrapAddr)
		if err != nil {
			log.Info(err)
		}
		bootstrapAddrs = append(bootstrapAddrs, addr)
	}
	return bootstrapAddrs
}
