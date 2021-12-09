package main

import (
	"fmt"
	"time"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

const (
	spiClockSpeed  physic.Frequency = physic.KiloHertz * 10
	spiBusDevPath  string           = "/dev/spidev0"
	spiDevPath     string           = spiBusDevPath + ".0"
	spiMode        spi.Mode         = spi.Mode2
	spiBitsPerWord int              = 8
)

func init() {
	_, err := host.Init()

	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Starting")
	//dataIn := make([]byte, len(dataOut), len(dataOut))
	//trxFunc, _, err := setupConnections(spiDevPath)
	ticker := time.NewTicker(time.Millisecond * 100)
	exitedChan := make(chan bool)

	// if err != nil {
	// 	panic(err)
	// }

	go runLoop(exitedChan, ticker.C)

	time.Sleep(time.Second * 10)
	exitedChan <- true
	time.Sleep(time.Second * 10)
}

func runLoop(exitedChan chan bool, trigger <-chan time.Time) {
	for {
		select {
		case _ = <-exitedChan:
			{
				fmt.Println("Exiting")
				return
			}
		case _ = <-trigger:
			{
				fmt.Println("Sending Data")
			}
		}
		//trxFunc(dataOut, dataIn)
	}
}

func setupConnections(spiDevPath string) (trxFunc func(w, r []byte) error, clsFunc func() error, err error) {
	var initErr error

	var spiPort spi.PortCloser
	var spiConn spi.Conn

	//Open port and connections
	spiPort, initErr = spireg.Open(spiDevPath)
	if initErr == nil {
		spiConn, initErr = spiPort.Connect(spiClockSpeed, spiMode, spiBitsPerWord)

		if initErr == nil {
			clsFunc = func() error {
				return spiPort.Close()
			}
			trxFunc = spiConn.Tx
		}
	}

	return trxFunc, clsFunc, initErr
}
