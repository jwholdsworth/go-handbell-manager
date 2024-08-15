package main

import (
	"log"

	"github.com/google/gousb"
)

var VENDOR_ID = gousb.ID(4094)
var PRODUCT_ID = gousb.ID(4104)

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == VENDOR_ID && desc.Product == PRODUCT_ID
	})

	for _, device := range devices {
		defer device.Close()
	}

	if err != nil {
		log.Fatalf("Error opening devices: %v", err)
	}

	if len(devices) == 0 {
		log.Fatal("No devices found")
	}

	// Pick the first device found.
	device := devices[0]

	// Switch the configuration to #1.
	cfg, err := device.Config(1)
	if err != nil {
		log.Fatalf("%s.Config(1): %v", device, err)
	}
	defer cfg.Close()

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		log.Fatalf("%s.Interface(0, 0): %v", cfg, err)
	}
	defer intf.Close()

	epIn, err := intf.InEndpoint(1)
	if err != nil {
		log.Fatalf("%s.InEndpoint(1): %v", intf, err)
	}

	// Buffer large enough for 10 USB packets from endpoint 6.
	buf := make([]byte, 10*epIn.Desc.MaxPacketSize)
	// total := 0
	// // Repeat the read/write cycle 10 times.
	// for i := 0; i < 10; i++ {
	// 	// readBytes might be smaller than the buffer size. readBytes might be greater than zero even if err is not nil.
	// 	readBytes, err := epIn.Read(buf)
	// 	if err != nil {
	// 		fmt.Println("Read returned an error:", err)
	// 	}
	// 	if readBytes == 0 {
	// 		log.Fatalf("IN endpoint 1 returned 0 bytes of data.")
	// 	}
	// 	total += readBytes
	// }
	// fmt.Printf("Total number of bytes copied: %d\n", total)

	inputStream, err := epIn.NewStream(epIn.Desc.MaxPacketSize, 10)
	if err != nil {
		log.Fatalf("Error reading the stream: %v", err)
	}
	input, _ := inputStream.Read(buf)
	log.Printf("Bytes read: %v", input)
}
