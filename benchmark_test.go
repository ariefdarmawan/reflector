package reflector_test

import (
	"testing"
	"time"
	"github.com/ariefdarmawan/reflector"
)

func BenchmarkSet(b *testing.B) {
	bCount := b.N
	for i := 0; i < bCount; i++ {
		data := new(obj)
		err := reflector.From(data).
			Set("ID", "Obj1").
			Set("Name", "Obj1 Name").
			Set("Int", 10).
			Set("Dec", float64(20.30)).
			Set("Date", time.Now()).
			Flush()
		if err!=nil {
			b.Fatalf("flush error: %s", err.Error())
		}
		if data.Dec!=20.30 {
			b.Fatalf("dec is not 20.30")
		}
	}
}


func BenchmarkGet(b *testing.B) {
	data := new(obj)
	data.Dec = 20.30
	data.ID = "Obj1"
		
	bCount := b.N
	for i := 0; i < bCount; i++ {
		dec, err := reflector.From(data).Get("Dec")
		if err!=nil {
			b.Fatalf("flush error: %s", err.Error())
		}
		if dec !=float64(20.30) {
			b.Fatalf("dec is not 20.30")
		}
	}
}

func BenchmarkCopy(b *testing.B) {
		obj1 := &struct {
			ID   string
			Name string
			F64  float64
			Sub  struct {
				Random string
			}
		}{"ID01", "Test Saja", 0.85, struct{ Random string }{"Random01"}}

		obj2 := &struct {
			ID   string
			Name string
			Sub  struct {
				Random string
			}
		}{}

		bCount := b.N
		for i:=0; i<bCount; i++ {
			obj2, e := reflector.CopyAttributeByNames(obj1, obj2, "Name", "F64")
			if e!=nil {
				b.Fatalf("error copy: %s", e.Error())
			}
			if obj1.Name!=obj2.Name {
				b.Fatalf("error copy: %s", "name is not equal")
			}	
		}
}