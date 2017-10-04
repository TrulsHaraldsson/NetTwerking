package d7024e

/*
This entire file content will be implemented into Kademlia later on, with these timers, messages will be sent periodically for storage and pings etc.
*/

//Containers in golang : https://www.youtube.com/watch?v=HPuvDm8IC-4

import (
  "fmt"
  "time"
  "math/rand"
)

type Timers struct{}

/*
* startTimer starts the timers that are used to periodically send pings to nodes.
*/
func (timers *Timers) startTimer(){
  ticker := time.NewTicker(time.Millisecond * 1000)
    go func() {
        for t := range ticker.C {
            fmt.Println("Tick")
            //Will change to kademliaSendPing
            timers.interval(t)
        }
    }()
    time.Sleep(time.Second * 10) //In final version, there will be no sleep (or atleast not one that stops within an hour or so).
    ticker.Stop()
    now := time.Now()
    fmt.Println("Time stopped at : ", now.String())
}

func (timers *Timers) showTimer(tick time.Time){
  fmt.Println("Tick-Tack", tick.String() ,"\n")
  time.Sleep(time.Millisecond * 100)
  diff := time.Since(tick)
  fmt.Println("Time difference : ", diff.String())
}

/*
* interval will be created in storage to mimic a timer until files are removed in RAM.
* The function used in storage for this is : MoveToMemory
* (ALTERNATIVE IDEA : With the timers in Kademlia, let storage listen to incoming "pings" stating file is active, if a timer runs out, move the unattractive file to Memory from RAM.)
*/
func (timers *Timers) interval(tick time.Time){
  random := rand.Intn(100)
  time.Sleep(time.Millisecond * time.Duration(random)) //Only to mimic lag.
  tickTack := time.Since(tick)
  fmt.Println("time : ", tickTack.String())
  if tickTack < time.Duration(time.Millisecond * 80) {
    fmt.Println("-Tack is within interval!\n")
    //DO NOTHING
  }else{
    fmt.Println("-Tack too late!\n")
    //Remove files from RAM
  }
}
