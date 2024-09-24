package main

import (
	"log"
	"sync"

	"github.com/google/gousb"
	"github.com/micmonay/keybd_event"
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
	Button1: keybd_event.VK_F9,        // start
	Button2: keybd_event.VK_G,         // go
	Button3: keybd_event.VK_A,         // bob
	Button4: keybd_event.VK_SEMICOLON, // single
}

var keys = map[int]int{
	1: keybd_event.VK_J,
	2: keybd_event.VK_F,
}

var keyboard keybd_event.KeyBonding

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

	keyboard, err = keybd_event.NewKeyBonding()
	if err != nil {
		log.Fatalf("There was a problem simulting key presses. The error was: %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(devices))

	for i, device := range devices {
		go loadController(device, i+1)
	}

	wg.Wait()
}

func loadController(device *gousb.Device, controllerNumber int) {
	defer closeDevice(device, controllerNumber)
	configuration, err := device.Config(1)
	if err != nil {
		log.Fatalf("Error getting configuration for controller %d: %v", controllerNumber, err)
	}
	defer closeConfiguration(configuration, controllerNumber)

	device.SetAutoDetach(true)

	intf, err := configuration.Interface(0, 0)
	if err != nil {
		log.Fatalf("Error reading interface for controller %d: %v", controllerNumber, err)
	}
	defer closeInterface(intf, controllerNumber)

	endpoint, err := intf.InEndpoint(1)
	if err != nil {
		log.Fatalf("Error reading endpoint for controller %d: %v", controllerNumber, err)
	}

	buffer := make([]byte, 10*endpoint.Desc.MaxPacketSize)

	buttonPressed := ButtonPress{
		First:  false,
		Second: false,
	}

	var lastStroke Stroke = Backstroke

	for {
		bytes, err := endpoint.Read(buffer)
		if err != nil {
			log.Fatalf("Unable to read from controller %d. The error was: %s", controllerNumber, err)
		}
		if bytes == 0 {
			log.Fatalf("Received 0 bytes from controller %d", controllerNumber)
		}

		// the 4th byte is the one we're interested in. The bits are organised as follows:
		// 128: 2nd button pressed
		// 64: 1st button pressed
		// 32-0: position and accellerometer
		// 18 is level, ~30 is vertical. Lower than 18 means a fast swing down. Higher than ~30 is a fast swing up.
		input := buffer[3]

		input = handleButtonPress(controllerNumber, input, &buttonPressed)

		if lastStroke == Handstroke && input < 18 {
			// ring the backstroke
			log.Printf("Backstroke rung by controller %d", controllerNumber)
			lastStroke = Backstroke
			sendKeyPress(keys[controllerNumber])
		}

		if lastStroke == Backstroke && input > 30 {
			// ring the handstroke
			log.Printf("Handstroke rung by controller %d", controllerNumber)
			lastStroke = Handstroke
			sendKeyPress(keys[controllerNumber])
		}
	}
}

func handleButtonPress(controller int, input byte, buttonPressed *ButtonPress) byte {
	if (input>>7)&1 == 1 {
		if !buttonPressed.First {
			log.Printf("Button 1 pressed on controller %d", controller)
			if controller == 1 {
				sendKeyPress(buttons.Button1)
			}
			if controller == 2 {
				sendKeyPress(buttons.Button3)
			}
		}
		buttonPressed.First = true
	} else {
		buttonPressed.First = false
	}

	if (input>>6)&1 == 1 {
		if !buttonPressed.Second {
			log.Printf("Button 2 pressed on controller %d", controller)
			if controller == 1 {
				sendKeyPress(buttons.Button2)
			}
			if controller == 2 {
				sendKeyPress(buttons.Button4)
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

func sendKeyPress(key int) {
	keyboard.SetKeys(key)
	keyboard.Press()
	keyboard.Release()
	keyboard.Clear()
}

func closeConfiguration(configuration *gousb.Config, controller int) {
	configuration.Close()
	log.Printf("Closing configuration for controller %d", controller)
}

func closeInterface(intf *gousb.Interface, controller int) {
	intf.Close()
	log.Printf("Closing interface for controller %d", controller)
}

func closeDevice(device *gousb.Device, controller int) {
	device.Close()
	log.Printf("Closing device for controller %d", controller)
}
