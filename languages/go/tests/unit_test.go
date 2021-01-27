package oso_test

import (
	"fmt"
	"reflect"
	"testing"

	oso "github.com/osohq/go-oso"
)

func TestSomething(t *testing.T) {
	fmt.Println("Hello, World!")

}

func TestNewOso(t *testing.T) {
	if o, err := oso.NewOso(); err != nil {
		t.Fatalf("Failed to set up Oso: %v", err)
	} else if reflect.TypeOf(o) != reflect.TypeOf(oso.Oso{}) {
		t.Fatalf("Expected type oso.Oso, got: %v", reflect.TypeOf(o))
	}
}

func TestLoadFile(t *testing.T) {
	var o oso.Oso
	var err error
	if o, err = oso.NewOso(); err != nil {
		t.Fatalf("Failed to set up Oso: %v", err)
	}

	if err = o.LoadFile("test.polar"); err != nil {
		t.Error(err.Error())
	}

	if err = o.LoadFile("test.polar"); err == nil {
		t.Error("Failed to error on loading duplicate file")
	}

	if err = o.LoadFile("test.txt"); err == nil {
		t.Error("Failed to error on loading non-polar file (.txt)")
	}

	if err = o.LoadFile("fake.polar"); err == nil {
		t.Error("Failed to error on loading non-existent file")
	}
}

func TestLoadString(t *testing.T) {
	var o oso.Oso
	var err error
	if o, err = oso.NewOso(); err != nil {
		t.Fatalf("Failed to set up Oso: %v", err)
	}

	if err = o.LoadString("f(1);"); err != nil {
		t.Error(err.Error())
	}

}

func TestClearRules(t *testing.T) {

}

func TestQueryStr(t *testing.T) {
	var o oso.Oso
	var err error
	if o, err = oso.NewOso(); err != nil {
		t.Fatalf("Failed to set up Oso: %v", err)
	}

	o.LoadString("f(1);")
	results, errors := o.QueryStr("f(x)")

	if err = <-errors; err != nil {
		t.Error(err.Error())
	} else {
		var got []map[string]interface{}
		expected := map[string]interface{}{"x": int64(1)}
		for elem := range results {
			got = append(got, elem)
		}
		if len(got) > 1 {
			t.Errorf("Received too many results: %v", got)
		} else if !reflect.DeepEqual(got[0], expected) {
			t.Errorf("Expected: %v, got: %v", expected, got[0])
		}
	}

	o.LoadString("g(x) if x.Fake();")
	results, errors = o.QueryStr("g(1)")

	if err = <-errors; err == nil {
		t.Error("Expected Polar runtime error, got none")
	}
}

func TestQueryRule(t *testing.T) {
	var o oso.Oso
	var err error
	if o, err = oso.NewOso(); err != nil {
		t.Fatalf("Failed to set up Oso: %v", err)
	}

	o.LoadString("f(1, 2);")
	results, errors := o.QueryRule("f", oso.ValueVariable("x"), oso.ValueVariable("y"))

	if err = <-errors; err != nil {
		t.Error(err.Error())
	} else {
		var got []map[string]interface{}
		expected := map[string]interface{}{"x": int64(1), "y": int64(2)}
		for elem := range results {
			got = append(got, elem)
		}
		if len(got) != 1 {
			t.Errorf("Received incorrect number of results: %v", got)
		} else if !reflect.DeepEqual(got[0], expected) {
			t.Errorf("Expected: %v, got: %v", expected, got[0])
		}
	}

	o.LoadString("g(x) if x.Fake();")
	results, errors = o.QueryRule("g", 1)

	if err = <-errors; err == nil {
		t.Error("Expected Polar runtime error, got none")
	}

}

func TestIsAllowed(t *testing.T) {
	var o oso.Oso
	var err error
	if o, err = oso.NewOso(); err != nil {
		t.Fatalf("Failed to set up Oso: %v", err)
	}

	o.LoadString("allow(\"foo\", \"bar\", \"baz\");")
	if a, e := o.IsAllowed("foo", "bar", "baz"); e != nil {
		t.Error(e.Error())
	} else if !a {
		t.Error("IsAllowed returned false, expected true")
	}

	if a, e := o.IsAllowed("foo", "baz", "bar"); e != nil {
		t.Error(e.Error())
	} else if a {
		t.Error("IsAllowed returned true, expected false")
	}

}
