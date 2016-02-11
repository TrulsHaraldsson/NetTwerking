package d7024e

import (
  "testing"
  "time"
  "fmt"
//  "reflect"
)

func TestStartTimer(t *testing.T){
  timers := Timers{}
  timers.startTimer()
}

func TestTimer(t *testing.T){
  timers := Timers{}
  ticker := time.NewTicker(time.Millisecond * 500)
    go func() {
        for t := range ticker.C {
            fmt.Println("Tick\n")
            timers.showTimer(t)
        }
    }()

  time.Sleep(time.Millisecond * 2600)
  ticker.Stop()
  now := time.Now()
  fmt.Println("Time stopped at : ", now.String())
}

func TestInterval(t *testing.T){
  timers := Timers{}
  ticker := time.NewTicker(time.Millisecond * 500)
    go func() {
        for t := range ticker.C {
            timers.interval(t)
        }
    }()

  time.Sleep(time.Millisecond * 2600)
  ticker.Stop()
  now := time.Now()
  fmt.Println("Time stopped at : ", now.String())
}
