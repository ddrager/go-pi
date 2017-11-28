package main

import (
  "os"
  "os/signal"
  "time"
  "fmt"
  "strconv"
  "github.com/ddrager/go-pi-blaster"
  "net/http"
  "github.com/zenazn/goji"
  "github.com/zenazn/goji/web"
)

const (
  redPin = 17
  greenPin = 22
  bluePin = 24

  jump = 5 // must be factor of 255
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
  stopAndClear()

  delay, _ := strconv.ParseInt(c.URLParams["delay"], 10, 64)

  if delay == 0 {
    delay = 30;
  }

  // set initial colors
  var red int64 = 255
  var green int64 = 0
  var blue int64 = 0

  control = 1

  fmt.Printf("Cycling with delay %d\n", delay)
  fmt.Fprintf(w, "Cycling with a delay of %d milliseconds", delay)

  // main loop
  go func() {
    timer := time.Tick(time.Millisecond * time.Duration(delay))
    for _ = range timer {
      if control == 0 { break }

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
    fmt.Printf("Cycle timer broken")
  }()
}

func christmas(c web.C, w http.ResponseWriter, r *http.Request) {

  stopAndClear()

  control = 1

  fmt.Printf("Turning on Christmas mode\n")

  delay := 2000
  color := 0 // current color

  // set initial color - before timer starts
  setRGB(0, 255, 0)

  // main loop
  go func() {
    timer := time.Tick(time.Millisecond * time.Duration(delay))
    for _ = range timer {

      if control == 0 { break }

      if color == 0 { 
        setRGB(255, 0, 0); 
        color = 1 
      } else {
        setRGB(0, 255, 0);
        color = 0
      }

    }
    fmt.Printf("Christmas timer broken")
  }()
}

// cycle colors
func christmas_fade(c web.C, w http.ResponseWriter, r *http.Request) {
  stopAndClear()

  delay, _ := strconv.ParseInt(c.URLParams["delay"], 10, 64)

  if delay == 0 {
    delay = 30;
  }

  // set initial colors
  var red int64 = 0
  var green int64 = 0
  var blue int64 = 0
  var direction int64 = 0
  var color int64 = 0

  control = 1

  fmt.Printf("Cycling with delay %d\n", delay)
  fmt.Fprintf(w, "Cycling with a delay of %d milliseconds", delay)

  // main loop
  go func() {
    timer := time.Tick(time.Millisecond * time.Duration(delay))
    for _ = range timer {

      if control == 0 { break }

      // make the jump in color
      if direction == 0 {
        if color == 0 { 
          red = red + jump
        } else { green = green + jump }
      } else if direction == 1 {
        if color == 0 { 
          red = red - jump 
        } else { green = green - jump }
      }
      // set the lights
      setRGB(red, green, blue)

      // switch colors & direction if at bottom of range
      if color == 0 && direction == 1 && red == 0 { color = 1; direction = 0;  
      } else if color == 1 && direction == 1 && green == 0 { color = 0; direction = 0; }

      // if at top of range, switch direction and keep color
      if (red >= 255 || green >= 255) && direction == 0 {
        direction = 1
      }

      if control == 0 { break }

      // fmt.Printf("Red: %d, Green: %d, Blue: %d\n", red, green, blue)
    }
    fmt.Printf("Cycle timer broken")
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
  setRGB(0, 0, 0)
  b.Reset();
}

func menu(c web.C, w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func websetRGB(c web.C, w http.ResponseWriter, r *http.Request) {
  stopAndClear()
  red, _ := strconv.ParseInt(c.URLParams["red"], 10, 64)
  green, _ := strconv.ParseInt(c.URLParams["green"], 10, 64)
  blue, _ := strconv.ParseInt(c.URLParams["blue"], 10, 64)
  setRGB(red, green, blue)
}

func main() {
  fmt.Printf("Running web server\n")

  // set initial pins and start interface
  a := []int64{redPin, greenPin, bluePin}
  b.Start(a)

  control = 1

  // reset pins when program ends
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)

  pwd, _ := os.Getwd()
  staticFilesLocation := pwd + "/static"

  fmt.Printf("Launching web server...\n")
  fmt.Printf("Loading static files from %s\n", staticFilesLocation)

  // web functionality
  goji.Get("/hello/:name", menu)
  goji.Get("/cycle", cycle)
  goji.Get("/cycle/:delay", cycle)
  goji.Get("/christmas", christmas)
  goji.Get("/christmas_fade", christmas_fade)
  goji.Get("/clear", clear)
  goji.Get("/setrgb/:red/:green/:blue", websetRGB)


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
