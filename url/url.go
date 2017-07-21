package url

import "fmt"

type UrlTestGo struct {
}

func (U *UrlTestGo) SetValue() (err error) {

	for i := 0; i < 100; i++ {
		a := 2
		b := 3
		c1 := a + b
		fmt.Println(c1)
	}
	return nil
}

func (U *UrlTestGo) Add (a int32, b int32) (r int32, err error) {
	println(a+b)
	return a+b,nil
}

func (U *UrlTestGo) Add100(a int32,b int32)(r int32,e error)  {
    return a+b+100,nil
}
