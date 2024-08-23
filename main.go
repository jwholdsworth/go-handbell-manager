package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/gousb"
)

const VENDOR_ID = gousb.ID(4094)
const PRODUCT_ID = gousb.ID(4104)
const (
	Handstroke float32 = -0.3
	Backstroke float32 = 0
)

var keys = map[int]string{
	1: "J",
	2: "F",
}

type ButtonPress struct {
	First  bool
	Second bool
}

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == VENDOR_ID && desc.Product == PRODUCT_ID
	})

	if err != nil {
		log.Panic("Could not find any devices")
	}

	log.Printf("Detected %d devices", len(devices))

	var wg sync.WaitGroup
	wg.Add(len(devices))

	for i, d := range devices {
		defer d.Close()
		go loadController(d, i)
	}

	wg.Wait()
}

func loadController(device *gousb.Device, controllerNumber int) {
	cfg, err := device.Config(1)
	if err != nil {
		log.Fatalf("Error getting configuration for controller %d: %v", controllerNumber, err)
	}
	defer cfg.Close()

	device.SetAutoDetach(true)

	intf, err := cfg.Interface(0, 0)
	if err != nil {
		log.Fatalf("Error reading interface for controller %d: %v", controllerNumber, err)
	}
	defer intf.Close()

	epIn, err := intf.InEndpoint(1)
	if err != nil {
		log.Fatalf("Error reading endpoint for controller %d: %v", controllerNumber, err)
	}

	buf := make([]byte, 10*epIn.Desc.MaxPacketSize)

	buttonPressed := ButtonPress{
		First:  false,
		Second: false,
	}

	for {
		readBytes, err := epIn.Read(buf)
		if err != nil {
			fmt.Println("Read returned an error:", err)
		}
		if readBytes == 0 {
			log.Fatalf("IN endpoint 6 returned 0 bytes of data.")
		}

		input := buf[3]

		// the 4th byte is the one we're interested in. The bits are organised as follows:
		// 128: 2nd button pressed
		// 64: 1st button pressed
		// 32-0: position and accellerometer
		// 18 is level, ~30 is vertical. Lower than 18 means a fast swing down. Higher than ~30 is a fast swing up.
		// fmt.Printf("%08b\t%d", input, input)

		fmt.Println(input)

		handleButtonPress(controllerNumber, input, &buttonPressed)

		// fmt.Printf("%06b\t%d\t", buf[3], buf[3])

		// fmt.Println()

		// fmt.Println(buf[0], "\t", buf[1], "\t", buf[2], "\t", buf[3], "\t", buf[4])
	}

}

func handleButtonPress(controller int, input byte, buttonPressed *ButtonPress) {
	if (input>>7)&1 == 1 {
		if !buttonPressed.First {
			fmt.Printf("Button 1 pressed on controller %d", controller)
			// button has just been pressed - do stuff!
		}
		buttonPressed.First = true
	} else {
		buttonPressed.First = false
	}

	if (input>>6)&1 == 1 {
		if !buttonPressed.Second {
			fmt.Printf("Button 2 pressed on controller %d", controller)
			// button has just been pressed - do stuff!
		}

		buttonPressed.Second = true
	} else {
		buttonPressed.Second = false
	}

	// clear bits 6&7
	input = byte(clearBit(input, 6))
	input = byte(clearBit(input, 7))
}

func clearBit(input byte, pos uint) byte {
	var mask byte = ^(1 << pos)
	input &= mask
	return input
}
