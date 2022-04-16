package util

type Glimit struct {
	Num int
	C   chan struct{}
}

func NewG(num int) *Glimit {
	return &Glimit{
		Num: num,
		C:   make(chan struct{}, num),
	}
}

func (g *Glimit) Run(f func(args interface{}), args interface{}) {
	g.C <- struct{}{}
	go func(args interface{}) {
		f(args)
		<-g.C
	}(args)
}
