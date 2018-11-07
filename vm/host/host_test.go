package host

import (
	"testing"

	"errors"

	"time"

	. "github.com/golang/mock/gomock"
	"github.com/iost-official/go-iost/vm/database"
)

func watchTime(f func()) time.Duration {
	ta := time.Now()
	f()
	tb := time.Now().Sub(ta)
	return tb
}

func sliceEqual(a, b []string) bool {
	if len(a) == len(b) {
		for i, s := range a {
			if s != b[i] {
				return false
			}
		}
		return true
	}
	return false
}

func myinit(t *testing.T, ctx *Context) (*database.MockIMultiValue, Host) {
	mockCtrl := NewController(t)
	defer mockCtrl.Finish()
	db := database.NewMockIMultiValue(mockCtrl)
	bdb := database.NewVisitor(100, db)

	//monitor := Monitor{}

	host := NewHost(ctx, bdb, nil, nil)
	return db, *host
}

func TestHost_Put(t *testing.T) {

	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Put(Any(), Any(), Any()).Do(func(a, b, c string) {
		if a != "state" || b != "b-contractName-hello" || c != "sworld" {
			t.Fatal(a, b, c)
		}
	})

	mock.EXPECT().Get("state", "b-contractName-hello").Return("", errors.New("not found"))

	host.Put("hello", "world")
	if host.cost["contractName"].Data != 24 {
		t.Fatal(host.cost)
	}
}

func TestHost_Put2(t *testing.T) {

	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Put(Any(), Any(), Any()).AnyTimes().Do(func(a, b, c string) {
		t.Log("put: ", a, b, c)
	})

	mock.EXPECT().Get("state", "b-contractName-hello").Return("sa", nil)

	host.Put("hello", "world")
	if host.cost["contractName"].Data != 4 {
		t.Fatal(host.cost)
	}
}

func TestHost_PutUserSpace(t *testing.T) {

	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Put(Any(), Any(), Any()).AnyTimes().Do(func(a, b, c string) {
		t.Log("put: ", a, b, c)
	})

	mock.EXPECT().Get("state", "b-contractName@abc-hello").Return("sa", nil)

	host.Put("hello", "world", "abc")
	if host.cost["abc"].Data != 4 {
		t.Fatal(host.cost)
	}

	v, _ := host.Get("hello", "abc")
	if v.(string) != "world" {
		t.Fatal(v)
	}
}

func TestHost_Del(t *testing.T) {
	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Put(Any(), Any(), Any()).AnyTimes().Do(func(a, b, c string) {
		t.Log("put: ", a, b, c)
	})

	mock.EXPECT().Get("state", "b-contractName-hello").Return("sworld", nil)
	mock.EXPECT().Get("state", "b-contractName@abc-hello").Return("sworld", nil)

	host.Del("hello")
	if host.cost["contractName"].Data != -24 {
		t.Fatal(host.cost)
	}
	host.Del("hello", "abc")
	if host.cost["contractName"].Data != -24 {
		t.Fatal(host.cost)
	}
}

func TestHost_Get(t *testing.T) {
	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Get(Any(), Any()).DoAndReturn(func(a, b string) (string, error) {
		if a != "state" || b != "b-contractName-hello" {
			t.Fatal(a, b)
		}
		return "sworld", nil
	})

	ans, _ := host.Get("hello")
	if ans != "world" {
		t.Fatal(ans)
	}
}

func TestHost_MapPut(t *testing.T) {

	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Put("state", "m-contractName-hello-1", Any()).Do(func(a, b, c string) {
		if a != "state" || b != "m-contractName-hello-1" || c != "sworld" {
			t.Fatal(a, b, c)
		}
	})
	mock.EXPECT().Put("state", "m-contractName-hello", Any()).Do(func(a, b, c string) {
		if c != "@1" {
			t.Fatal(c)
		}
	})
	mock.EXPECT().Has("state", "m-contractName-hello-1").Return(false, nil)
	mock.EXPECT().Get("state", "m-contractName-hello").Return("", errors.New("not found"))
	mock.EXPECT().Get("state", "m-contractName-hello-1").Return("", errors.New("not found"))

	tr := watchTime(func() {
		host.MapPut("hello", "1", "world")
	})
	if tr > time.Millisecond {
		t.Log("to slow")
	}

	if host.cost["contractName"].Data != 26 {
		t.Fatal(host.cost)
	}
}

func TestHost_MapPut_Owner(t *testing.T) {

	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Put("state", "m-contractName@abc-hello-1", Any()).Do(func(a, b, c string) {
		if c != "sworld" {
			t.Fatal(c)
		}
	})
	mock.EXPECT().Put("state", "m-contractName@abc-hello", Any()).Do(func(a, b, c string) {
		if c != "@1" {
			t.Fatal(c)
		}
	})
	mock.EXPECT().Has("state", "m-contractName@abc-hello-1").Return(false, nil)
	mock.EXPECT().Get("state", "m-contractName@abc-hello").Return("", errors.New("not found"))
	mock.EXPECT().Get("state", "m-contractName@abc-hello-1").Return("", errors.New("not found"))

	tr := watchTime(func() {
		host.MapPut("hello", "1", "world", "abc")
	})
	if tr > time.Millisecond {
		t.Log("to slow")
	}

	if host.cost["abc"].Data != 30 {
		t.Fatal(host.cost)
	}
}

func TestHost_MapGet(t *testing.T) {

	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Get(Any(), Any()).DoAndReturn(func(a, b string) (string, error) {
		if a != "state" || b != "m-contractName-hello-1" {
			t.Fatal(a, b)
		}
		return "sworld", nil
	})

	ans, _ := host.MapGet("hello", "1")
	if ans != "world" {
		t.Fatal(ans)
	}
}

func TestHost_MapGet_Owner(t *testing.T) {

	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Get(Any(), Any()).DoAndReturn(func(a, b string) (string, error) {
		if a != "state" || b != "m-contractName@abc-hello-1" {
			t.Fatal(a, b)
		}
		return "sworld", nil
	})

	ans, _ := host.MapGet("hello", "1", "abc")
	if ans != "world" {
		t.Fatal(ans)
	}
}

func TestHost_MapKeys(t *testing.T) {

	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Get("state", "m-contractName-hello").Return("@a@b@c", nil)

	ans, _ := host.MapKeys("hello")
	if !sliceEqual(ans, []string{"a", "b", "c"}) {
		t.Fatal(ans)
	}
}

func TestHost_MapKeys_Owner(t *testing.T) {

	ctx := NewContext(nil)
	ctx.Set("commit", "abc")
	ctx.Set("contract_name", "contractName")

	mock, host := myinit(t, ctx)

	mock.EXPECT().Get("state", "m-contractName@abc-hello").Return("@a@b@c", nil)

	ans, _ := host.MapKeys("hello", "abc")
	if !sliceEqual(ans, []string{"a", "b", "c"}) {
		t.Fatal(ans)
	}
}

func TestHost_BlockInfo(t *testing.T) {

}

func TestTeller_Transfer(t *testing.T) {
	ctx := NewContext(nil)
	ctx.Set("contract_name", "contractName")
	ctx.Set("auth_list", map[string]int{"hello": 1, "b": 0})

	mock, host := myinit(t, ctx)

	var (
		ihello = int64(1000)
		iworld = int64(0)
	)

	mock.EXPECT().Get(Any(), Any()).AnyTimes().DoAndReturn(func(table string, key string) (string, error) {
		switch key {
		case "i-hello":
			return database.MustMarshal(ihello), nil
		case "i-world":
			return database.MustMarshal(iworld), nil
		}
		return database.MustMarshal(nil), nil
	})

	mock.EXPECT().Put(Any(), Any(), Any()).AnyTimes().DoAndReturn(func(a, b, c string) error {
		t.Log("put:", a, b, database.MustUnmarshal(c))
		switch b {
		case "i-hello":
			ihello = database.MustUnmarshal(c).(int64)
		case "i-world":
			iworld = database.MustUnmarshal(c).(int64)
		}

		return nil
	})

	host.Transfer("hello", "world", "3")
	host.Transfer("hello", "world", "3")

}
