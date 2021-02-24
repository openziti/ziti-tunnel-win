package service

import (
	"bufio"
	"encoding/json"
	"io"
	"net"
	"strings"
	"time"

	"github.com/Microsoft/go-winio"
	"github.com/openziti/desktop-edge-win/service/ziti-tunnel/dto"
)

type fetchFromRTS func([]string, *dto.TunnelStatus, string) dto.Response

func sendMessagetoPipe(ipcPipeConn net.Conn, commandMsg *dto.CommandMsg, args []string) error {
	writer := bufio.NewWriter(ipcPipeConn)
	enc := json.NewEncoder(writer)

	err := enc.Encode(commandMsg)
	if err != nil {
		log.Error("could not encode or writer list identities message, %v", err)
		return err
	}

	log.Debug("Message sent to ipc pipe")

	writer.Flush()

	return nil
}

func readMessageFromPipe(ipcPipeConn net.Conn, readDone chan struct{}, fn fetchFromRTS, args []string, flag string) {
	for {

		reader := bufio.NewReader(ipcPipeConn)
		msg, err := reader.ReadString('\n')

		if err != nil {
			log.Error(err)
			return
		}

		if len(args) == 0 {
			args = append(args, "all")
		}

		dec := json.NewDecoder(strings.NewReader(msg))
		var tunnelStatus dto.ZitiTunnelStatus
		if err := dec.Decode(&tunnelStatus); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		if tunnelStatus.Status != nil {
			responseMsg := fn(args, tunnelStatus.Status, flag)
			if responseMsg.Code == SUCCESS {
				log.Info(responseMsg.Message)
				log.Info(responseMsg.Payload)
			} else {
				log.Info(responseMsg.Message)
				log.Info(responseMsg.Error)
			}
		} else {
			log.Errorf("Ziti tunnel retuned nil status")
		}

		readDone <- struct{}{}
		break

	}
	return
}

//GetIdentities is to fetch identities through cmdline
func GetIdentities(args []string, flag string) {
	getDataFromIpcPipe(&GET_STATUS, GetIdentitiesFromRTS, args, flag)
}

//GetServices is to fetch services through cmdline
func GetServices(args []string, flag string) {
	getDataFromIpcPipe(&GET_STATUS, GetServicesFromRTS, args, flag)
}

func getDataFromIpcPipe(commandMsg *dto.CommandMsg, fn fetchFromRTS, args []string, flag string) {
	log.Infof("fetching identities through cmdline...%s", args)

	log.Debug("Connecting to pipe")
	timeout := 2000 * time.Millisecond
	ipcPipeConn, err := winio.DialPipe(ipcPipeName(), &timeout)
	defer closeConn(ipcPipeConn)

	if err != nil {
		log.Errorf("Connection to ipc pipe is not established, %v", err)
		return
	}
	readDone := make(chan struct{})
	defer close(readDone) // ensure that goroutine exits

	go readMessageFromPipe(ipcPipeConn, readDone, fn, args, flag)

	err = sendMessagetoPipe(ipcPipeConn, commandMsg, args)
	if err != nil {
		log.Errorf("Message is not sent to ipc pipe, %v", err)
		return
	}

	log.Debugf("Connection to ipc pipe is established - %s and remote address %s", ipcPipeConn.LocalAddr().String(), ipcPipeConn.RemoteAddr().String())

	<-readDone
	log.Debug("read finished normally")
}