// log a description of events when pressing button #1 or moving hat#1.
// 10sec timeout.
package main

import (
	"log"
	"sync"

	"git.tcp.direct/kayos/sendkeys"
	"github.com/google/gousb"
	. "github.com/splace/joysticks"
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

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == VENDOR_ID && desc.Product == PRODUCT_ID
	})

	if err != nil {
		log.Panic("Could not find any devices")
	}

	var wg sync.WaitGroup
	wg.Add(len(devices))

	for i, d := range devices {
		defer d.Close()
		go loadController(i + 1)
	}
	wg.Wait()
}

func loadController(controller int) {
	device := Connect(controller)

	// using Connect allows a device to be interrogated
	log.Printf("Action XL Controller %d: Buttons:%d, Hats:%d\n", controller, len(device.Buttons), len(device.HatAxes)/2)

	// get/assign channels for specific events
	b1press := device.OnClose(1)
	b2press := device.OnClose(2)
	v1move := device.OnPanX(1)

	// start feeding OS events onto the event channels.
	go device.ParcelOutEvents()

	var lastStroke = Backstroke

	key, err := sendkeys.NewKBWrapWithOptions(sendkeys.Noisy)
	if err != nil {
		log.Panic(err)
	}

	for {
		select {
		case <-b1press:
			log.Println("button #1 pressed")
		case <-b2press:
			log.Println("button #2 pressed")
		case v := <-v1move:
			vpos := v.(AxisEvent)
			if vpos.V >= float32(Backstroke) && lastStroke == Handstroke {
				// backstroke rung
				lastStroke = Backstroke
				log.Printf("Controller %d backstroke rung", controller)

				// send keyboard signal
				key.Type(keys[controller])
			}
			if vpos.V <= float32(Handstroke) && lastStroke == Backstroke {
				// handstroke rung
				lastStroke = Handstroke
				log.Printf("Controller %d handstroke rung", controller)

				// send keyboard signal
				key.Type(keys[controller])
			}
		}
	}
}
