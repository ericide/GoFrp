package slave

import (
	"time"
)

type DataObject struct {
	Pre        *[]byte
	Data       *[]byte
	DataLength int64
}

func Start(remoteCmdHost string, localHost string) {

	for {
		connectRemoteController(remoteCmdHost, localHost)
		time.Sleep(2 * time.Second)
	}
}

func connectRemoteController(remoteCmdHost string, localHost string) {

	tunnelClient := TunnelClient{
		TunnelServerAddress: remoteCmdHost,
		ExtAddress:          localHost,
	}
	tunnelClient.Start()
}
