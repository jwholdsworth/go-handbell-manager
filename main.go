package main

import (
	"fmt"
	"log"
	"sync"

	"git.tcp.direct/kayos/sendkeys"
	"github.com/google/gousb"
)

const (
	PRODUCT_ID = gousb.ID(4104)
	VENDOR_ID  = gousb.ID(4094)
)

const (
	Handstroke Stroke = iota
	Backstroke Stroke = iota
)

// button map with standard Abel defaults
var buttons = KeyMap{
	Button1: "Z", // start - this needs to be F9
	Button2: "G", // go
	Button3: "A", // bob
	Button4: ";", // single
}

var keys = map[int]string{
	1: "J",
	2: "F",
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

	for i := 0; i < len(devices); i++ {
		defer devices[i].Close()
		go loadController(devices[i], i+1)
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

	key, err := sendkeys.NewKBWrapWithOptions(sendkeys.Noisy)
	if err != nil {
		log.Panic(err)
	}

	var lastStroke Stroke = Backstroke

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

		input = handleButtonPress(controllerNumber, input, &buttonPressed, key)

		if lastStroke == Handstroke && input < 18 {
			// ring the backstroke
			fmt.Printf("Backstroke rung by controller %d\n", controllerNumber)
			lastStroke = Backstroke
			key.Type(keys[controllerNumber])
		}

		if lastStroke == Backstroke && input > 30 {
			// ring the handstroke
			fmt.Printf("Handstroke rung by controller %d\n", controllerNumber)
			lastStroke = Handstroke
			key.Type(keys[controllerNumber])
		}

	}

}

func handleButtonPress(controller int, input byte, buttonPressed *ButtonPress, key *sendkeys.KBWrap) byte {
	if (input>>7)&1 == 1 {
		if !buttonPressed.First {
			fmt.Printf("Button 1 pressed on controller %d\n", controller)
			if controller == 1 {
				key.Type(buttons.Button1)
			}
			if controller == 2 {
				key.Type(buttons.Button3)
			}
		}
		buttonPressed.First = true
	} else {
		buttonPressed.First = false
	}

	if (input>>6)&1 == 1 {
		if !buttonPressed.Second {
			fmt.Printf("Button 2 pressed on controller %d\n", controller)
			if controller == 1 {
				key.Type(buttons.Button2)
			}
			if controller == 2 {
				key.Type(buttons.Button4)
			}
		}

		buttonPressed.Second = true
	} else {
		buttonPressed.Second = false
	}

	// clear bits 6&7
	input = byte(clearBit(input, 6))
	input = byte(clearBit(input, 7))

	return input
}

func clearBit(input byte, pos uint) byte {
	var mask byte = ^(1 << pos)
	input &= mask
	return input
}
