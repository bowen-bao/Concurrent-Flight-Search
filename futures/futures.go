package futures

type Future struct{
	c  chan interface{}
}

// Create new future
func NewFuture() *Future{
	c := make(chan interface{})
	return &Future{c}
}

//Send the result of p to the Future's channel
func (f *Future) Set(p interface{}){
	go func(){
		for {
			f.c <- p
		}
	}()
}

//Retrieve the value of the Future's channel
func (f *Future) Get() interface{} {
	return <- f.c
}
