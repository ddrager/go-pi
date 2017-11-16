Go-pi
=====

This repository is a playground for my [Go](https://golang.org/)-[Raspberry Pi](http://www.raspberrypi.org/) activities. I use Raspbian.

LED Strip Control
=================

One of the items I have attached to my Pi via GPIO is this [LED RGB Strip](http://www.amazon.com/gp/product/B00JA8NUZU/ref=as_li_tl?ie=UTF8&camp=1789&creative=390957&creativeASIN=B00JA8NUZU&linkCode=as2&tag=theblackhol0a-20&linkId=DODKDIB3FLTKVW3D). I followed @popoklopsi's tutorial [here](http://popoklopsi.github.io/RaspberryPi-LedStrip/) for wiring. The circuit looks like this:

![Wiring diagram for Raspberry Pi LEDs](https://raw.githubusercontent.com/ddrager/go-pi/master/static/img/rgb-pi-led.png)

*Image credit [Mitchtech](http://mitchtech.net/raspberry-pi-pwm-rgb-led-strip/raspi_rgb_led/)*

Installation & Run
==================

* [Install Go](https://golang.org/dl/) on the Raspberry pI.
* Install [pi-blaster](https://github.com/sarfata/pi-blaster), a C driver for PWM management on the Raspberry Pi
* Fetch: `go get -v github.com/ddrager/go-pi` Error is normal here. 
* Run: `go run web-led.go` to run web server.

Other resources
===============

Other useful resources if you are interested in getting into this type of thing:

* [RBG-pi server](https://github.com/ryupold/rgb-pi) (c/python)
