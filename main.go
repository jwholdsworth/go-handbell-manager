// log a description of events when pressing button #1 or moving hat#1.
// 10sec timeout.
package main

import (
	"log"
	"time"

	"github.com/google/gousb"
	. "github.com/splace/joysticks"
)

var VENDOR_ID = gousb.ID(4094)

var PRODUCT_ID = gousb.ID(4104)

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()
	devices, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == VENDOR_ID && desc.Product == PRODUCT_ID
	})

	if err != nil {
		log.Panic("Could not find any devices")
	}

	for i, d := range devices {
		defer d.Close()
		loadController(i + 1)
	}
	// try connecting to specific controller.
	// the index is system assigned, typically it increments on each new controller added.
	// indexes remain fixed for a given controller, if/when other controller(s) are removed.

	log.Println("Timeout in 10 secs.")
	time.Sleep(time.Second * 30)
	log.Println("Shutting down due to timeout.")
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

	// handle event channels
	go func() {
		for {
			select {
			case <-b1press:
				log.Println("button #1 pressed")
			case <-b2press:
				log.Println("button #2 pressed")
			case v := <-v1move:
				vpos := v.(AxisEvent)
				log.Printf("Controller %d moved to %f", controller, vpos.V)
			}
		}
	}()
}
