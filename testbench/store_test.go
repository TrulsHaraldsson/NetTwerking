package main

import (
  "fmt"
  "testing"
  "../golang"
  "time"
)
/*
Connection pattern : nodes {8401,8402,8403} all connect to 8400.
*/

func TestStoreOnce(t *testing.T){
  //start := d7024e.StartNode(8400, "none", "2111111400000000000000000000000000000000")
  A := d7024e.NewKademlia(8400, "2111111400000000000000000000000000000000")
  A.Start(8400)
  time.Sleep(50 * time.Millisecond)

  //d7024e.StartNode(8401, "localhost:8400", "2111111400000000000000000000000000000001")
  B := d7024e.NewKademlia(8401, "2111111400000000000000000000000000000001")
  B.Start(8401)
  B.Ping("localhost:8400")
  time.Sleep(50 * time.Millisecond)

  C := d7024e.NewKademlia(8402, "2111111400000000000000000000000000000002")
  C.Start(8402)
  C.Ping("localhost:8400")
  time.Sleep(50 * time.Millisecond)

  D := d7024e.NewKademlia(8403, "2111111400000000000000000000000000000003")
  D.Start(8403)
  D.Ping("localhost:8400")
  time.Sleep(50 * time.Millisecond)

  fmt.Println("All nodes connected")
  contact := A.SendStoreMessage(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), []byte("Test store"))
  if string(contact) != "stored"{
    t.Error("Value not stored!")
  }else{
    fmt.Println("Complete store.")
  }
}

func TestMultiStore(t *testing.T){
  A := d7024e.NewKademlia(8410, "2111111400000000000000000000000000000010")
  A.Start(8410)
  time.Sleep(50 * time.Millisecond)

  B := d7024e.NewKademlia(8411, "2111111400000000000000000000000000000011")
  B.Start(8411)
  B.Ping("localhost:8410")
  time.Sleep(50 * time.Millisecond)

  contact1 := A.SendStoreMessage(d7024e.NewKademliaID("2111111400000000000000000000000000000010"), []byte("First"))
  if string(contact1) != "stored"{
    t.Error("Value not stored!", contact1)
  }else{
    fmt.Println("First store complete")
    contact2 := A.SendStoreMessage(d7024e.NewKademliaID("2111111400000000000000000000000000000010"), []byte("Second"))
    if string(contact2) != "stored"{
      t.Error("Value not stored!", contact2)
    }else{
      fmt.Println("Second store complete.")
      contact3 := A.SendStoreMessage(d7024e.NewKademliaID("2111111400000000000000000000000000000010"), []byte("Third"))
      if string(contact3) != "stored"{
        t.Error("Value not stored!", contact3)
      }else{
        fmt.Println("Third store complete.")
      }
    }
  }
}
