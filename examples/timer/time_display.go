package main

import (
  "fmt"
  "time"
  "sync"
  "github.com/leesper/tao"
)

func main() {
  wg := &sync.WaitGroup{}
  wheel := tao.NewTimingWheel()
  timerId := wheel.AddTimer(
    time.Now().Add(2 * time.Second),
    500 * time.Millisecond,
    tao.NewOnTimeOut(nil, func(t time.Time) { fmt.Printf("TIME OUT AT %s\n", t) }))
  fmt.Printf("Add timer %d\n", timerId)

  timerId = wheel.AddTimer(
    time.Now().Add(3 * time.Second),
    50 * time.Millisecond,
    tao.NewOnTimeOut(nil, func(t time.Time) { fmt.Printf("CANCEL ME IF YOU CAN\n") }))
  fmt.Printf("Add another timer %d, now we have %d timers\n", timerId, wheel.Size())

  wg.Add(1)
  go func() {
    for i := 0; i < 20; i++ {
      select {
      case timeout := <-wheel.TimeOutChan:
        timeout.Callback(time.Now())
      }
    }
    wg.Done()
  }()
  wg.Wait()

  wheel.CancelTimer(timerId)
  fmt.Printf("Cancel timer %d, now we have %d timers\n", timerId, wheel.Size())

  wg.Add(1)
  go func() {
    for i := 0; i < 10; i++ {
      select {
      case timeout := <-wheel.TimeOutChan:
        timeout.Callback(time.Now())
      }
    }
    wg.Done()
  }()
  wg.Wait()
  wheel.Stop()
}
