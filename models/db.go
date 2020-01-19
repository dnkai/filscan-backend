package models

import (
	"fmt"
	"log"
	"time"

	"github.com/globalsign/mgo"
)

var ps = fmt.Sprintf
var globalS *mgo.Session
var DB string
var TimeNow int64 = 0
var TimeNowStr string = "2006-01-02 15:04:05"

func TimenowInit() {
	go func() {
		for {
			tim_t := time.Now()
			TimeNow = tim_t.Unix()
			TimeNowStr = tim_t.Format("2006-01-02 15:04:05")
			time.Sleep(time.Second)
		}
	}()
	return
}

func init() { }

func GetGlobalSession(host, user, pass, mongoDB string) *mgo.Session {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{host},
		Username: user,
		Password: pass,
		Database: mongoDB,
	}
	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalln("create mongodb session error ", err)
	}
	globalS = s
	return s
}

func connect(collection string) (*mgo.Session, *mgo.Collection) {
	s := globalS.Copy()
	c := s.DB(DB).C(collection)
	return s, c
}

func Connect(collection string) (*mgo.Session, *mgo.Collection) {
	return connect(collection)
}

func Insert(collection string, docs ...interface{}) error {
	ms, c := connect(collection)
	defer ms.Close()
	return c.Insert(docs...)
}

func BulkUpsert(collection string, docs []interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Upsert(docs[:]...)
	return bulk.Run()
}

func BulkUpdate(collection string, docs []interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Update(docs[:]...)
	return bulk.Run()
}

func Upsert(collection string, selector interface{}, docs interface{}) (info *mgo.ChangeInfo, err error) {
	ms, c := connect(collection)
	defer ms.Close()
	return c.Upsert(selector, docs)
}

func IsExist(collection string, query interface{}) bool {
	ms, c := connect(collection)
	defer ms.Close()
	count, _ := c.Find(query).Count()
	return count > 0
}

func FindOne(collection string, query, selector, result interface{}) error {
	ms, c := connect(collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}

func FindAll(collection string, query, selector, result interface{}) error {
	ms, c := connect(collection)
	defer ms.Close()
	return c.Find(query).Select(selector).All(result)
}
func FindSortLimit(collection, sort string, query, selector, result interface{}, begindex, count int) error {
	ms, c := connect(collection)
	defer ms.Close()
	//fmt.Print(c.Find(query).Select(selector).Sort(sort).Limit(count).Skip(begindex).Count())
	return c.Find(query).Select(selector).Sort(sort).Limit(count).Skip(begindex).All(result)
}

func FindSortCollationLimit(collection, sort string, query, selector, result interface{}, begindex, count int, collation *mgo.Collation) error {
	ms, c := connect(collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Collation(collation).Sort(sort).Limit(count).Skip(begindex).All(result)
}

func FindAllLimit(collection string, query, selector, result interface{}, begindex, count int) error {
	ms, c := connect(collection)
	defer ms.Close()
	//fmt.Print(c.Find(query).Select(selector).Sort(sort).Limit(count).Skip(begindex).Count())
	return c.Find(query).Select(selector).Limit(count).Skip(begindex).All(result)
}

func FindCount(collection string, query, selector interface{}) (int, error) {
	ms, c := connect(collection)
	defer ms.Close()
	return c.Find(query).Select(selector).Count()
}
func Update(collection string, query, update interface{}) error {
	ms, c := connect(collection)
	defer ms.Close()
	return c.Update(query, update)
}

func Remove(collection string, query interface{}) error {
	ms, c := connect(collection)
	defer ms.Close()
	return c.Remove(query)
}

func Distinct(collection, distinctKey string, query, result interface{}) error {
	ms, c := connect(collection)
	defer ms.Close()
	return c.Find(query).Distinct(distinctKey, result)

}
func AggregateAll(collection string, query interface{}, result interface{}) error {
	ms, c := connect(collection)
	defer ms.Close()

	return c.Pipe(query).Iter().All(result)
}
func AggregateOne(collection string, query interface{}, result interface{}) error {
	ms, c := connect(collection)
	defer ms.Close()
	return c.Pipe(query).One(result)
}
