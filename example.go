package main

import (
  "fmt"
  "github.com/ddrager/go-pi-blaster"
)



func main() {
  fmt.Printf("Running\n")
  a := []int64{17, 22, 24}

  var b piblaster.Blaster


  b.Start(a)
  b.Apply(17, 1.1);

  b.Apply(22, -7.0);

  b.Apply(24, .56);

  b.DumpCurrent();
}
