package redact_test

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/ucarion/redact"
)

type demo struct {
	Password     string
	SecretNumber int
	SecretArray  []string
	SecretMap    map[string]interface{}
	Demo2        demo2
	Demo2Array   []demo2
	Demo2Map     map[string]*demo2
	ChanFoo      chan int
}

type demo2 struct {
	NestedPassword string
	Demo3          demo3
}

type demo3 struct {
	DeeplyNestedPassword string
}

func TestRedact(t *testing.T) {
	cases := []struct {
		path []string
		in   demo
		out  demo
	}{
		{
			[]string{"Password"},
			demo{Password: "foo"},
			demo{Password: ""},
		},
		{
			[]string{"SecretNumber"},
			demo{SecretNumber: 42},
			demo{SecretNumber: 0},
		},
		{
			[]string{"SecretArray"},
			demo{SecretArray: []string{"foo"}},
			demo{SecretArray: nil},
		},
		{
			[]string{"SecretMap"},
			demo{SecretMap: map[string]interface{}{"foo": ""}},
			demo{SecretMap: nil},
		},
		{
			[]string{"Demo2", "NestedPassword"},
			demo{Password: "foo", Demo2: demo2{NestedPassword: "foo"}},
			demo{Password: "foo", Demo2: demo2{NestedPassword: ""}},
		},
		{
			[]string{"Demo2", "Demo3", "DeeplyNestedPassword"},
			demo{Password: "foo", Demo2: demo2{NestedPassword: "foo", Demo3: demo3{DeeplyNestedPassword: "foo"}}},
			demo{Password: "foo", Demo2: demo2{NestedPassword: "foo", Demo3: demo3{DeeplyNestedPassword: ""}}},
		},
		{
			[]string{"Demo2Array", "NestedPassword"},
			demo{Demo2Array: []demo2{demo2{NestedPassword: "foo"}, demo2{NestedPassword: "foo"}}},
			demo{Demo2Array: []demo2{demo2{NestedPassword: ""}, demo2{NestedPassword: ""}}},
		},
		{
			[]string{"Demo2Map", "foo", "NestedPassword"},
			demo{Demo2Map: map[string]*demo2{"foo": &demo2{NestedPassword: "foo"}, "bar": &demo2{NestedPassword: "foo"}}},
			demo{Demo2Map: map[string]*demo2{"foo": &demo2{NestedPassword: ""}, "bar": &demo2{NestedPassword: "foo"}}},
		},
		{
			[]string{"Demo2Map", "foo"},
			demo{Demo2Map: map[string]*demo2{"foo": &demo2{NestedPassword: "foo"}}},
			demo{Demo2Map: map[string]*demo2{"foo": nil}},
		},
	}

	for i, tt := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			redact.Redact(tt.path, &tt.in)
			if !reflect.DeepEqual(tt.in, tt.out) {
				t.Errorf("path %+v: want %+v, got %+v", tt.path, tt.out, tt.in)
			}
		})
	}
}

func ExampleRedact() {
	type user struct {
		Name     string
		Password string
	}

	u := user{Name: "John", Password: "letmein"}
	redact.Redact([]string{"Password"}, &u)

	users := []user{
		user{Name: "John", Password: "letmein"},
		user{Name: "Mary", Password: "123456"},
	}
	redact.Redact([]string{"Password"}, &users)

	usersMap := map[string]*user{
		"a": &user{Name: "John", Password: "letmein"},
		"b": &user{Name: "Mary", Password: "123456"},
	}
	redact.Redact([]string{"a", "Password"}, &usersMap)

	fmt.Println(u)
	fmt.Println(users)
	fmt.Println(usersMap["a"], usersMap["b"])

	// Output:
	// {John }
	// [{John } {Mary }]
	// &{John } &{Mary 123456}
}
