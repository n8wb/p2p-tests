package main

import (
	"os"
	"log"
	"fmt"
	"sync"
	"os/signal"
	"strings"
	//p2pd "github.com/libp2p/go-libp2p-daemon"
	ma "github.com/multiformats/go-multiaddr"
	peer "github.com/libp2p/go-libp2p-peer"
	ps "github.com/libp2p/go-libp2p-peerstore"
	// "syscall"
)

func chanwait() {
	var end_waiter sync.WaitGroup
	end_waiter.Add(1)
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
	}()
	end_waiter.Wait()
}

func runServer() {
	server := new(Server)
	defer server.Close()
	// handleSignals()
	err := server.Start()
	if err != nil {
		log.Println(err)
		return
	}
}

func runClient() {
	client := new(Client)
	defer client.Close()
	err := client.Init()
	if err != nil {
		log.Println(err)
		return
	}
	for {
		response, err := client.ClientExec("peepeepoopookak")
		if err != nil {
			log.Println(err)
		}
		log.Println(response)
	}
}

// func handleSignals() {
// 	signals := make(chan os.Signal, 1)

// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
// 	<-signals
// 	log.Println("signal received")
// }



func CreatePeerInfos(peers []string) ([]ps.PeerInfo,error) {
	out := []ps.PeerInfo{}
	for _,rawPeer := range peers {
		pidSock := strings.Split(rawPeer,"@")
		socket := pidSock[1]
		rawPid := pidSock[0]
		ipPort := strings.Split(socket,":")
		port := ipPort[1]
		ip := ipPort[0]

		mAddr,err := ma.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s",ip,port))
		if err != nil {
			return nil,err
		}
		pid,err := peer.IDB58Decode(rawPid)
		if err != nil {
			return nil,err
		}

		out = append(out,ps.PeerInfo{
			ID:pid,
			Addrs: []ma.Multiaddr{mAddr},
		})
	}
	return out,nil
}