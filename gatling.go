package pasnet	// import "github.com/nathanaelle/pasnet"

import	(
	"time"
	"fmt"
	"reflect"
)



type	o_bullet	struct {
	F	func(int, interface {}) error
	A	interface{}
}

type	o_gatling []o_bullet


type	gatlingerror struct {
	expected	reflect.Type
	given		reflect.Type
}

func (ge *gatlingerror) Error() string {
	return fmt.Sprintf("incompatible types : Expected [ %s ] Given [ %s ]", ge.expected, ge.given )
}

// hand made bullet
func bullet_hm(f interface{}, v interface{}) (ret func(int)error, err error)  {
	V_ret	:= reflect.ValueOf(&ret).Elem()
	T_fd	:= reflect.TypeOf(int(0))
	V_f	:= reflect.ValueOf(f)
	V_v	:= reflect.ValueOf(v)
	T_exp	:= reflect.FuncOf([]reflect.Type{ T_fd, V_v.Type() }, []reflect.Type{ reflect.ValueOf(&err).Elem().Type() }, false )

	if V_f.Type() != T_exp {
		return nil, &gatlingerror{ T_exp, V_f.Type() }
	}

	V_ret.Set(reflect.MakeFunc(V_ret.Type(), func(v []reflect.Value) []reflect.Value {
		return V_f.Call([]reflect.Value{ v[0], V_v })
	}))

	return ret,nil
}


func bullet(f interface{}, v interface{}) (func(int)error) {
	ret,err := bullet_hm(f, v)
	if err != nil {
		panic(err)
	}

	return	ret
}


func gatling(bullets ...func(int)error) func(int)error {
	return func(fd int) error {
		for _,b := range bullets {
			err := b(fd)
			if err != nil {
				return err
			}
		}

		return nil
	}
}



func gatling_run(g o_gatling, fd int) (err error) {
	for _,b := range g {
		err	= b.F(fd, b.A)
		if err != nil {
			return err
		}
	}
	return nil
}

func bullet_bool(f func(int,bool)error) (func(int,interface{})error)  {
	return func(fd int, i interface{}) error {
		v, ok := i.(bool)
		if ok {
			return f(fd, v)
		}
		panic("WTF Wrong Type for bool !!!")
	}
}

func bullet_nil(f func(int)error) (func(int,interface{})error)  {
	return func(fd int, _ interface{}) error {
		return f(fd)
	}
}



func bullet_int(f func(int,int)error) (func(int,interface{})error)  {
	return func(fd int, i interface{}) error {
		v, ok := i.(int)
		if ok {
			return f(fd, v)
		}
		panic("WTF Wrong Type for int !!!")
	}
}

func bullet_duration(f func(int,time.Duration)error) (func(int,interface{})error)  {
	return func(fd int, i interface{}) error {
		v, ok := i.(time.Duration)
		if ok {
			return f(fd, v)
		}
		panic("WTF Wrong Type for Duration !!!")
	}
}
