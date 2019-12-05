## Reflector
Reflector is a Go object manipulator

Sometime we will need to update object dynamically or programatically. This is what reflector is used for.

```go
type childObj struct {
	Name   string
	Values []int
}

type obj struct {
	ID   string
	Name string
	Int  int
	Dec  float64
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
```