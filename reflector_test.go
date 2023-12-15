package reflector_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/ariefdarmawan/reflector"
	"github.com/sebarcode/codekit"

	"github.com/smartystreets/goconvey/convey"
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

				cv.Convey("get", func() {
					name, err := reflector.From(data).Get("Name")
					cv.So(err, cv.ShouldBeNil)
					cv.So(name, cv.ShouldEqual, data.Name)
				})
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

type testObj struct {
	ID    string
	Name  string
	Value int
	Dt    time.Time
}

func TestAssignVar(t *testing.T) {
	cv.Convey("Assign Var", t, func() {
		source := testObj{"ID", "Name", 100, time.Now()}

		cv.Convey("Same Type Ptr", func() {
			dest := new(testObj)
			e := reflector.AssignValue(reflect.ValueOf(&source), reflect.ValueOf(dest))
			cv.So(e, cv.ShouldBeNil)
			cv.So(dest, cv.ShouldResemble, &source)
		})

		cv.Convey("Same Type Value", func() {
			dest := new(testObj)
			e := reflector.AssignValue(reflect.ValueOf(source), reflect.ValueOf(dest))
			cv.So(e, cv.ShouldBeNil)
			cv.So(dest, cv.ShouldResemble, &source)
		})

		cv.Convey("Same Type M", func() {
			dest := codekit.M{}
			rs := reflect.ValueOf(source)
			rd := reflect.ValueOf(&dest)
			e := reflector.AssignValue(rs, rd)
			cv.So(e, cv.ShouldBeNil)
			cv.So(codekit.ToInt(dest["Value"], codekit.RoundingAuto), cv.ShouldResemble, source.Value)
		})

		cv.Convey("Array", func() {
			sources := []testObj{
				{"ID0", "Name 0", 100, time.Now()},
				{"ID0", "Name 0", 100, time.Now()},
			}
			rs := reflect.ValueOf(sources)
			dest := []codekit.M{}
			rd := reflect.ValueOf(&dest)
			e := reflector.AssignValue(rs, rd)
			cv.So(e, cv.ShouldBeNil)
			cv.So(len(sources), cv.ShouldEqual, len(dest))
			cv.So(codekit.ToInt(dest[0]["Value"], codekit.RoundingAuto), cv.ShouldResemble, sources[0].Value)
		})

		cv.Convey("Diff type with some same field name", func() {
			dest := struct {
				ID   string
				Name string
			}{}
			e := reflector.AssignValue(reflect.ValueOf(source), reflect.ValueOf(&dest))
			cv.So(e, cv.ShouldBeNil)
			cv.So(dest.Name, cv.ShouldResemble, source.Name)
		})

		cv.Convey("Negative", func() {
			dest := struct {
				ID   int
				Name string
			}{}
			e := reflector.AssignValue(reflect.ValueOf(source), reflect.ValueOf(dest))
			cv.So(e, cv.ShouldNotBeNil)
		})
	})
}

func TestAssignSlice(t *testing.T) {
	cv.Convey("assing slice", t, func() {
		makeDest := func() []testObj {
			return []testObj{
				{"ID01", "Name01", 100, time.Now()},
				{"ID02", "Name02", 200, time.Now().Add(1 * time.Minute)},
				{"ID03", "Name03", 300, time.Now().Add(2 * time.Minute)},
			}
		}

		makeDestPtr := func() []*testObj {
			return []*testObj{
				{"ID01", "Name01", 100, time.Now()},
				{"ID02", "Name02", 200, time.Now().Add(1 * time.Minute)},
				{"ID03", "Name03", 300, time.Now().Add(2 * time.Minute)},
			}
		}

		cv.Convey("copy to same type with ptr source", func() {
			dest := makeDest()
			e := reflector.AssignSliceItem(reflect.ValueOf(&testObj{"ID04", "Name04", 400, time.Now()}), 3, reflect.ValueOf(&dest))
			cv.So(e, cv.ShouldBeNil)
			cv.So(len(dest), cv.ShouldEqual, 4)
			cv.So(dest[0].ID, cv.ShouldEqual, "ID01")
			cv.So(dest[3].ID, cv.ShouldEqual, "ID04")
		})

		cv.Convey("copy to same type with value source", func() {
			dest := makeDest()
			e := reflector.AssignSliceItem(reflect.ValueOf(testObj{"ID04", "Name04", 400, time.Now()}), 3, reflect.ValueOf(&dest))
			cv.So(e, cv.ShouldBeNil)
			cv.So(len(dest), cv.ShouldEqual, 4)
			cv.So(dest[0].ID, cv.ShouldEqual, "ID01")
			cv.So(dest[3].ID, cv.ShouldEqual, "ID04")
		})

		cv.Convey("copy to same type in ptr with ptr source", func() {
			dest := makeDestPtr()
			e := reflector.AssignSliceItem(reflect.ValueOf(&testObj{"ID04", "Name04", 400, time.Now()}), 3, reflect.ValueOf(&dest))
			cv.So(e, cv.ShouldBeNil)
			cv.So(len(dest), cv.ShouldEqual, 4)
			cv.So(dest[0].ID, cv.ShouldEqual, "ID01")
			cv.So(dest[3].ID, cv.ShouldEqual, "ID04")
		})

		cv.Convey("copy to same type in ptr with value source", func() {
			dest := makeDestPtr()
			e := reflector.AssignSliceItem(reflect.ValueOf(testObj{"ID04", "Name04", 400, time.Now()}), 3, reflect.ValueOf(&dest))
			cv.So(e, cv.ShouldBeNil)
			cv.So(len(dest), cv.ShouldEqual, 4)
			cv.So(dest[0].ID, cv.ShouldEqual, "ID01")
			cv.So(dest[3].ID, cv.ShouldEqual, "ID04")
		})

		cv.Convey("copy to same type in ptr with value source on existing index", func() {
			dest := makeDestPtr()
			e := reflector.AssignSliceItem(reflect.ValueOf(testObj{"ID04", "Name04", 400, time.Now()}), 2, reflect.ValueOf(&dest))
			cv.So(e, cv.ShouldBeNil)
			cv.So(len(dest), cv.ShouldEqual, 3)
			cv.So(dest[0].ID, cv.ShouldEqual, "ID01")
			cv.So(dest[2].ID, cv.ShouldEqual, "ID04")
		})

		cv.Convey("copy to []M with ptr source", func() {
			dest := []codekit.M{}
			rs := reflect.ValueOf(&testObj{"ID04", "Name04", 400, time.Now()})
			rd := reflect.ValueOf(&dest)
			e := reflector.AssignSliceItem(rs, 0, rd)
			cv.So(e, cv.ShouldBeNil)
			cv.So(len(dest), cv.ShouldEqual, 1)
			cv.So(dest[0].GetString("ID"), cv.ShouldEqual, "ID04")
		})

		cv.Convey("copy to []*M with ptr source", func() {
			dest := []*codekit.M{}
			e := reflector.AssignSliceItem(reflect.ValueOf(testObj{"ID04", "Name04", 400, time.Now()}), 0, reflect.ValueOf(&dest))
			cv.So(e, cv.ShouldBeNil)
			cv.So(len(dest), cv.ShouldEqual, 1)
			cv.So(dest[0].GetString("ID"), cv.ShouldEqual, "ID04")
		})

	})
}

func TestStructChild(t *testing.T) {
	type person struct {
		Name       string
		Salutation string
	}

	type e1 struct {
		Person person
		Role   string
	}

	type e2 struct {
		Person *person
		Role   string
	}

	cv.Convey("no ptr child", t, func() {
		p := &person{"Arief D", "Mr"}
		empl := new(e1)
		empl.Person = *p
		empl.Role = "Founder"

		refl := reflector.From(empl)
		refl.Set("Person.Salutation", "Tn.")
		refl.Flush()
		cv.So(empl.Person.Salutation, cv.ShouldEqual, "Tn.")

		get, _ := refl.Get("Person.Salutation")
		cv.So(get, cv.ShouldEqual, "Tn.")

		cv.Convey("ptr child", func() {
			p := &person{"Arief D", "Mr"}
			empl := new(e2)
			empl.Person = p
			empl.Role = "Founder"

			refl := reflector.From(empl)
			refl.Set("Person.Salutation", "Tn.")
			refl.Flush()
			cv.So(empl.Person.Salutation, cv.ShouldEqual, "Tn.")

			get, _ := refl.Get("Person.Salutation")
			cv.So(get, cv.ShouldEqual, "Tn.")

			cv.Convey("ptr child with null value", func() {
				empl := new(e2)
				empl.Role = "Founder"

				refl := reflector.From(empl)
				refl.Set("Person.Salutation", "Tn.")
				refl.Flush()
				cv.So(empl.Person.Salutation, cv.ShouldEqual, "Tn.")

				get, _ := refl.Get("Person.Salutation")
				cv.So(get, cv.ShouldEqual, "Tn.")
			})
		})
	})
}

func TestCreateFromPtr(t *testing.T) {
	convey.Convey("create from ptr", t, func() {
		objSource := new(testObj)
		objSource.ID = "create_from_ptr"
		objSource.Name = "random name"

		convey.Convey("copy", func() {
			objCopy, err := reflector.CreateFromPtr(objSource, true)
			convey.So(err, convey.ShouldBeNil)
			convey.So(objCopy.Name, convey.ShouldEqual, objSource.Name)
		})

		convey.Convey("not copy", func() {
			objNotCopy, err := reflector.CreateFromPtr(objSource, false)
			convey.So(err, convey.ShouldBeNil)
			convey.So(objNotCopy.Name, convey.ShouldEqual, "")
		})

		convey.Convey("not ptr", func() {
			objNotPtr := testObj{}
			_, err := reflector.CreateFromPtr(objNotPtr, false)
			convey.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestGetTo(t *testing.T) {
	convey.Convey("get to", t, func() {
		objSource := new(testObj)
		objSource.ID = "create_from_ptr"
		objSource.Name = "random name"
		objSource.Dt = time.Now()

		convey.Convey("positive ", func() {
			rf := reflector.From(objSource)
			name := ""
			dt := time.Time{}

			ename := rf.GetTo("Name", &name)
			edt := rf.GetTo("Dt", &dt)
			convey.So(ename, convey.ShouldBeNil)
			convey.So(name, convey.ShouldEqual, objSource.Name)
			convey.So(edt, convey.ShouldBeNil)
			convey.So(dt, convey.ShouldEqual, objSource.Dt)

			convey.Convey("negative", func() {
				ename = rf.GetTo("Name", name)
				convey.So(ename, convey.ShouldNotBeNil)
			})
		})
	})
}
