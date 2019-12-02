import "sync/atomic"

type Counter struct{
	value int
}

func NewCounter() Counter {
	ch := make(chan int,1)

	go func(){
		for{
			switch(){
				case
				default:
					value ++
					ch<-value
			}
		}

	}()

}

func (c *Counter) increment() int32 {
return atomic.AddInt32((*int32)(c), 1)
}

func (c *count32) value() int32 {
return atomic.LoadInt32((*int32)(c))
}