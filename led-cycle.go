package main

import (
  "os"
  "os/signal"
  "time"
  "fmt"
  "github.com/ddrager/go-pi-blaster"
)

const redPin = 17
const greenPin = 22
const bluePin = 24

const jump = 5


var a []int64
var b piblaster.Blaster

func setRGB(red int64, green int64, blue int64) {
  b.Apply(redPin, float64(red)/255.0)  
  b.Apply(greenPin, float64(green)/255.0)
  b.Apply(bluePin, float64(blue)/255.0)
}

func main() {
  fmt.Printf("Running\n")
  var red int64 = 255
  var green int64 = 0
  var blue int64 = 0
  a := []int64{redPin, greenPin, bluePin}

  b.Start(a)

  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)


  // main loop
  go func() {
    timer := time.Tick(time.Millisecond * 1)
    for _ = range timer {
      if red >= 255 && blue <= 0 && green < 255 {
        green = green + jump
      } else if green >= 255 && blue <= 0 && red > 0 {
        red = red - jump
      } else if red <= 0 && green >= 255 && blue < 255 {
        blue = blue + jump
      } else if red <= 0 && blue >= 255 && green > 0 {
        green = green - jump
      } else if green <= 0 && blue >= 255 && red < 255 {
        red = red + jump
      } else if red >= 255 && green <= 0 && blue > 0 {
        blue = blue - jump
      }
      setRGB(red, green, blue)
      //fmt.Printf("Red: %d, Green: %d, Blue: %d\n", red, green, blue)
    }
  }()

  //time.Sleep(time.Millisecond * 15000)

    for sig := range c {
      fmt.Printf("captured %v, exiting..", sig)
      //b.DumpCurrent();
      b.Reset();
      os.Exit(1)
    }


}
