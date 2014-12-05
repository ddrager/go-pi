package main

import (
  "os"
  "os/signal"
  "time"
  "fmt"
  "github.com/ddrager/go-pi-blaster"

  "net/http"
  "github.com/zenazn/goji"
  "github.com/zenazn/goji/web"
)

const (
  redPin = 17
  greenPin = 22
  bluePin = 24

  jump = 17

  staticFilesLocation = "/home/pi/gocode/src/github.com/ddrager/go-pi/static"
)

var (
  a []int64
  b piblaster.Blaster
  control int = 0
)

func setRGB(red int64, green int64, blue int64) {
  b.Apply(redPin, float64(red)/255.0)  
  b.Apply(greenPin, float64(green)/255.0)
  b.Apply(bluePin, float64(blue)/255.0)
}

// cycle colors
func cycle(c web.C, w http.ResponseWriter, r *http.Request) {

  // set initial colors
  var red int64 = 255
  var green int64 = 0
  var blue int64 = 0

  control = 1

  // main loop
  go func() {
    timer := time.Tick(time.Millisecond * 30)
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

      if control == 0 { break }

      //fmt.Printf("Red: %d, Green: %d, Blue: %d\n", red, green, blue)
    }   
  }() 
}

func clear(c web.C, w http.ResponseWriter, r *http.Request) {
  stopAndClear()
  setRGB(0, 0, 0)
}

func stopAndClear() {
  control = 0
  // allow running processes to stop
  time.Sleep(100 * time.Millisecond)
  b.Reset();
}

func menu(c web.C, w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func main() {
  fmt.Printf("Running\n")

  // set initial pins and start interface
  a := []int64{redPin, greenPin, bluePin}
  b.Start(a)

  control = 1

  // reset pins when program ends
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)

  fmt.Printf("Launching web server...\n")
  
  // web functionality
  goji.Get("/hello/:name", menu)
  goji.Get("/cycle", cycle)
  goji.Get("/clear", clear)
  

  goji.Handle("/*", http.FileServer(http.Dir(staticFilesLocation)))
  goji.Serve()

  // capture control c, probably a better way to do this
  //time.Sleep(time.Millisecond * 15000)
  for sig := range c {
    fmt.Printf("captured %v, exiting..", sig)
    //b.DumpCurrent()
    stopAndClear()
    os.Exit(1)
  }


}
