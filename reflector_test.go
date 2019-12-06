package reflector_test

import (
	"testing"
	"time"

	"github.com/ariefdarmawan/reflector"

	cv "github.com/smartystreets/goconvey/convey"
)

type childObj struct {
	Name   string
	Values []int
}

type obj struct {
	ID   string `name:"_id"`
	Name string
	Int  int
	Dec  float64 `name:"decimal"`
	Date time.Time

	Children []*childObj
}

func TestReflector(t *testing.T) {
	cv.Convey("reflector", t, func() {
		data := new(obj)
		err := reflector.From(data).
			Set("ID", "Obj1").
			Set("Name", "Obj1 Name").
			Set("Int", 10).
			Set("Dec", float64(20.30)).
			Set("Date", time.Now()).
			Flush()
		cv.So(err, cv.ShouldBeNil)
		cv.So(data.Dec, cv.ShouldEqual, 20.30)

		cv.Convey("update child", func() {
			children := []*childObj{}
			children = append(children, &childObj{"child1", []int{10, 20, 30}}, &childObj{"child2", []int{11, 21, 31}})
			err = reflector.From(data).Set("Children", children).Flush()
			cv.So(err, cv.ShouldBeNil)
			cv.So(data.Children[1].Values[1], cv.ShouldEqual, 21)

			cv.Convey("update child entity", func() {
				err = reflector.From(data.Children[0]).Set("Values", []int{1, 2, 3}).Flush()
				cv.So(err, cv.ShouldBeNil)
				cv.So(data.Children[0].Values[2], cv.ShouldEqual, 3)
			})
		})
	})
}

func TestNegative(t *testing.T) {
	cv.Convey("negative test", t, func() {
		data := obj{}
		err := reflector.From(data).Set("ID", "obj1").Flush()
		cv.So(err, cv.ShouldNotBeNil)
	})
}

func TestName(t *testing.T) {
	cv.Convey("get actual name", t, func() {
		data := obj{}
		names, err := reflector.From(&data).FieldNames("")
		cv.So(err, cv.ShouldBeNil)
		cv.So(names, cv.ShouldResemble, []string{"ID", "Name", "Int", "Dec", "Date", "Children"})

		cv.Convey("get masked names", func() {
			names, err := reflector.From(&data).FieldNames("name")
			cv.So(err, cv.ShouldBeNil)
			cv.So(names, cv.ShouldResemble, []string{"_id", "Name", "Int", "decimal", "Date", "Children"})
		})
	})
}
