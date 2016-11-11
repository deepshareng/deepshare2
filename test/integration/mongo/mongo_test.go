package mongo

import (
	"reflect"
	"testing"

	"github.com/MISingularity/deepshare2/pkg/testutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type person struct {
	Name  string
	Phone string
}

// TestMongoSetup assumes that mongod is running in background
// and tests some basic functions of mgo library.
func TestMongoSetup(t *testing.T) {
	var err error
	testDBName := "test"

	session := testutil.MustNewLocalSession()
	defer session.Close()
	// "people" collection
	c := session.DB(testDBName).C("people")

	// Prepare phase:
	// - Insert some entries into collection
	// We can later reuse partial information of inserted entries to test query.
	persons := []*person{
		&person{"Ale", "+55 53 8116 9639"},
		&person{"Cla", "+55 53 8402 8510"},
	}
	// mgo Insert method accepts `interface{}` variadic arugments.
	personsForInsert := make([]interface{}, len(persons))
	for i, p := range persons {
		personsForInsert[i] = p
	}
	err = c.Insert(personsForInsert...)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		p   *person
		err error
	}{
		{persons[0], nil},
		{persons[1], nil},
		{&person{"No Such Person", ""}, mgo.ErrNotFound},
	}

	for i, tt := range tests {
		result := &person{}
		err = c.Find(bson.M{"name": tt.p.Name}).One(result)
		if err != nil {
			if err != tt.err {
				t.Errorf("#%d: err=%v, want=%s", i, err, tt.err)
			}
			continue
		}
		if !reflect.DeepEqual(result, tt.p) {
			t.Errorf("#%d: person=%#v, want=%#v", i, result, tt.p)
		}
	}

	err = session.DB(testDBName).DropDatabase()
	if err != nil {
		t.Fatal(err)
	}
}
