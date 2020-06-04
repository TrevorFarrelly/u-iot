package exp

import (
	"fmt"
)

type StatusVar struct {
  name string
  t string
  val interface{}
  c chan interface{}
}

func AddStatusVar(name, t string) (*interface{}, error) {
  var n *StatusVar
  switch t {
    case "bool":
      n = &StatusVar{
        name: name,
        t: t,
        val: false,
        c: make(chan interface{}),
      }
    case "int":
      n = &StatusVar{
        name: name,
        t: t,
        val: 0,
        c: make(chan interface{}),
      }
    case "string":
      n = &StatusVar{
        name: name,
        t: t,
        val: "",
        c: make(chan interface{}),
      }
    case "float64":
      n = &StatusVar{
        name: name,
        t: t,
        val: 0.0,
        c: make(chan interface{}),
      }
    default:
      return nil, fmt.Errorf("unsupported type: %s", t)
  }
  return &n.val, nil
}
