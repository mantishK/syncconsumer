# Syncconsumer
## synchronous consumer for Go.

Often times we come across a usecase where we want to append the data to a slice from asynchronous go routines.
This tiny lib can help you do exactly that.

## download
```
go get github.com/mantishK/syncconsumer@latest
```

## usage

```go
import "github.com/mantishK/syncconsumer"

func main() {
  	all := make([]int, 0)
	s := syncconsumer.New(
		func(data int) {
			all = append(all, data)
		})

	wg := sync.WaitGroup{}

	// send data
	wg.Add(1)
	go func() {
		s.Publish(1, 2, 3)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		s.Publish(8, 9, 10)
		wg.Done()
	}()
	wg.Wait()
	s.WaitNClose()

	fmt.Println(all)
}
```
