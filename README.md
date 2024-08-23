# go-handbell-manager

This is a very low-fi hacked together version of Graham John's excellent Handbell Manager, designed to run on Linux.

Right now, it's a very basic application that reads the raw USB data from the Action XL Motion Controllers and sends a corresponding keypress to whichever application is in the foreground, should you move it fast/far enough. That being said, it works well enough for a couple of hours hackery on a Friday night!

## Installation instructions (Ubuntu)

1. To install it, copy the contents below into a text file, and change "myname" to be what your primary user group is (this is likely to be the same as your username).
1. Then copy it from your text editor into a terminal and run it - this will download it, and update your Udev permissions to allow the program to talk to the controllers over USB.
1. Plug in your Action XL Motion Controllers - if you already have, unplug them and plug them back in to allow Udev to find them.
1. Run `./go-handbell-manager`

```
wget -O go-handbell-manager https://github.com/jwholdsworth/go-handbell-manager/releases/download/v1.0.0/go-handbell-manager
chmod +x go-handbell-manager
echo 'SUBSYSTEMS=="usb", ATTRS{idVendor}=="0ffe", ATTRS{idProduct}=="1008", GROUP="myname", MODE="0666"' | sudo tee -a /etc/udev/rules.d/50-action-xl-motion-controllers.rules
sudo udevadm control --reload
```

## Todo

* Write some tests
* Sort out the button presses so you can start/stop/go like you can in Graham's Windows version
* Target just the Abel process so it doesn't start typing "F" and "J" all over the place
* Document any external libraries that might need to be installed (e.g. libusb)
