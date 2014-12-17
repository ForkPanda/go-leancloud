package leancloud

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

var cloud = &Client{}

func init() {
	log.SetFlags(log.Lshortfile)

	cfg := Config{}
	cfg.AppId = ""
	cfg.AppKey = ""
	cfg.MasterKey = ""
	cfg.UsingMaster = true
	cloud.Cfg = cfg
	cloud.BeforeRequest = func(r *http.Request) *http.Request {
		//data, _ := httputil.DumpRequestOut(r, true)
		//log.Println(string(data))
		return r
	}

	rand.Seed(time.Now().UnixNano())
}

func randString() string {
	return fmt.Sprintf("abc%d", rand.Int())
}

func TestObject(t *testing.T) {
	className := "NewClass"
	o1 := NewObject()
	o1.Set("key", "value")
	err := o1.Save(cloud, className, true)
	if err != nil {
		t.Fatal(err)
	}
	if o1.ObjectId() == "" {
		t.Fatal("null objectId")
	}
	o1.Set("updatekey", "updatevalue")
	err = o1.Update(cloud, className)
	if err != nil {
		t.Fatal(err)
	}
	o2, err := FetchObject(cloud, className, o1.ObjectId(), "")
	if err != nil {
		t.Fatal(err)
	}
	err = o2.Delete(cloud, className)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDate(t *testing.T) {
	className := "Class2"
	o1 := NewObject()
	d := FormatDate(time.Now())
	o1.Set("key", d)
	err := o1.Save(cloud, className, true)
	if err != nil {
		t.Fatal(err)
	}
	err = o1.Delete(cloud, className)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUser(t *testing.T) {
	u1 := NewUser()
	email := fmt.Sprintf("%s@email.com", randString())
	phone := fmt.Sprintf("1386818%0d%0d", rand.Intn(99), rand.Intn(99))
	username := randString()
	password := "password"
	r1, err := u1.Register(cloud, username, password, email, phone)
	if err != nil {
		t.Fatal(r1, err)
	}
	r2, err := u1.Login(cloud, username, password)
	if err != nil {
		t.Fatal(r2, err)
	}
}

func TestCQL(t *testing.T) {
	r, err := CQL(cloud, "select * from _User where username like 'abc%'")
	if err != nil {
		t.Fatal(err, r)
	}
}

func TestCloudFunction(t *testing.T) {
	r, err := CallFunction(cloud, "syncDate", "")
	if err != nil {
		t.Fatal(err, r)
	}
}

func Test1(t *testing.T) {
	//t.Fatal(time.Now().UTC().Format(time.RFC3339))
}
