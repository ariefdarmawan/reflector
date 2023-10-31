package reflector_test

import (
	"testing"

	"github.com/ariefdarmawan/reflector"
	"github.com/smartystreets/goconvey/convey"
)

func TestCopyField(t *testing.T) {
	convey.Convey("copy", t, func() {
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

		convey.Convey("validate", func() {
			obj2, e := reflector.CopyAttributes(obj1, obj2, "ID")
			convey.So(e, convey.ShouldBeNil)
			convey.So(obj2.ID, convey.ShouldBeBlank)
			convey.So(obj2.Sub.Random, convey.ShouldEqual, obj1.Sub.Random)
		})
	})
}

func TestCopyFieldByNames(t *testing.T) {
	convey.Convey("copy", t, func() {
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

		obj2, e := reflector.CopyAttributeByNames(obj1, obj2, "Name", "F64")
		convey.Convey("validate", func() {
			convey.So(e, convey.ShouldBeNil)
			convey.So(obj2.ID, convey.ShouldBeBlank)
			convey.So(obj2.Sub.Random, convey.ShouldBeBlank)
			convey.So(obj2.Name, convey.ShouldEqual, obj1.Name)
		})
	})
}

func TestCopyFieldNegative(t *testing.T) {
	convey.Convey("copy", t, func() {
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
			Sub  float64
		}{}

		_, e := reflector.CopyAttributes(obj1, obj2, "ID")
		convey.Convey("validate", func() {
			convey.So(e, convey.ShouldNotBeNil)
		})
	})
}
