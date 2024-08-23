# go-handbell-manager

This is a very low-fi hacked together version of Graham John's excellent Handbell Manager, designed to run on Linux.

Right now, it's a very basic application that reads the raw USB data from the Action XL Motion Controllers and sends a corresponding keypress to whichever application is in the foreground, should you move it fast/far enough. That being said, it works well enough for a couple of hours hackery on a Friday night!

## Todo

* Write some tests
* Sort out the button presses so you can start/stop/go like you can in Graham's Windows version
* Target just the Abel process so it doesn't start typing "F" and "J" all over the place
* Plus numerous other things for sure!
* Document how to not need to run it as root (probably some apparmour nonesense in Ubuntu)
* Document any external libraries that might need to be installed (e.g. libusb)
