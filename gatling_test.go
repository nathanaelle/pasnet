package pasnet	// import "github.com/nathanaelle/pasnet"

import (
	"errors"
	"testing"
)



func Test_Gatling(t *testing.T) {
	g1	:= gatling(

		bullet(func(fd int, b bool)error{
			t.Logf("%d + %+v",fd, b)
			return nil
		}, true),

		bullet(func(fd int, b bool)error{
			t.Logf("%d + %+v",fd, b)
			return nil
		}, false),

		bullet(func(fd int, b string)error{
			t.Logf("%d + %+v",fd, b)
			return nil
		}, "hello world"),
	)

	if err	:= g1(0); err != nil {
		t.Errorf("g1 %s", err)
	}

	if err	:= g1(10); err != nil {
		t.Errorf("g1 %s", err)
	}

	g2	:= gatling(

		bullet(func(fd int, b bool)error{
			t.Logf("%d + %+v",fd, b)
			return nil
		}, true),

		bullet(func(fd int, b bool)error{
			t.Logf("%d + %+v",fd, b)
			return errors.New("normal error")
		}, false),

		bullet(func(fd int, b string)error{
			t.Logf("%d + %+v",fd, b)
			return errors.New("impossible error")
		}, "not seen"),
	)

	if err	:= g2(0); err.Error() != "normal error" {
		t.Errorf("g1 %+v", err)
	}

	if err	:= g2(10); err.Error() != "normal error" {
		t.Errorf("g1 %+v", err)
	}

}
