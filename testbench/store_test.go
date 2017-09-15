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
  start := d7024e.StartNode(8400, "none", "2111111400000000000000000000000000000000")
  time.Sleep(50 * time.Millisecond)
  d7024e.StartNode(8401, "localhost:8400", "2111111400000000000000000000000000000001")
  d7024e.StartNode(8402, "localhost:8400", "2111111400000000000000000000000000000002")
  d7024e.StartNode(8403, "localhost:8400", "2111111400000000000000000000000000000003")
  fmt.Println("All nodes connected", start)
  contact := start.SendStoreMessage(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), []byte("Test store"))
  if string(contact) != "stored"{
    t.Error("Value not stored!", contact)
  }else{
    fmt.Println("Complete store.")
  }
}

func TestMultiStore(t *testing.T){
  start := d7024e.StartNode(8410, "none", "3111111400000000000000000000000000000000")
  time.Sleep(100 * time.Millisecond) // Let system react properly.
  d7024e.StartNode(8411, "localhost:8410", "3111111400000000000000000000000000000001")
  contact1 := start.SendStoreMessage(d7024e.NewKademliaID("3111111400000000000000000000000000000000"), []byte("First"))
  if string(contact1) != "stored"{
    t.Error("Value not stored!", contact1)
  }else{
    fmt.Println("First store complete")
    contact2 := start.SendStoreMessage(d7024e.NewKademliaID("3111111400000000000000000000000000000000"), []byte("Second"))
    if string(contact2) != "stored"{
      t.Error("Value not stored!", contact2)
    }else{
      fmt.Println("Second store complete.")
      contact3 := start.SendStoreMessage(d7024e.NewKademliaID("3111111400000000000000000000000000000000"), []byte("Third"))
      if string(contact3) != "stored"{
        t.Error("Value not stored!", contact3)
      }else{
        fmt.Println("Third store complete.")
      }
    }
  }
}
