package main

import (
	"flag"
	"os"
	"time"

	"regexp"

	"github.com/MISingularity/deepshare2/deepstats/aggregate"
	"github.com/MISingularity/deepshare2/pkg/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//this program is only used for splitting collections by appID
func main() {
	fs := flag.NewFlagSet("deepstatsd", flag.ExitOnError)

	mongoAddr := fs.String("mongo-addr", "", "Specify the raw data mongo database URL")
	mongoDBDayAggreagate := fs.String("mongodb-day-aggregate", "deepstats", "Specify the Mongo database for aggregate")
	mongoCollDayAggregate := fs.String("mongocoll-day-aggregate", "day", "Specify the Mongo collection for aggregate")

	mongoDBTotalAggreagate := fs.String("mongodb-total-aggregate", "deepstats", "Specify the Mongo database for aggregate")
	mongoCollTotalAggregate := fs.String("mongocoll-total-aggregate", "total", "Specify the Mongo collection for aggregate")

	mongoDBHourAggreagate := fs.String("mongodb-hour-aggregate", "deepstats", "Specify the Mongo database for aggregate")
	mongoCollHourAggregate := fs.String("mongocoll-hour-aggregate", "hour", "Specify the Mongo collection for aggregate")

	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	if *mongoAddr == "" || *mongoDBHourAggreagate == "" || *mongoCollHourAggregate == "" ||
		*mongoDBDayAggreagate == "" || *mongoCollDayAggregate == "" ||
		*mongoDBTotalAggreagate == "" || *mongoCollTotalAggregate == "" {
		log.Fatal("Invalid args")
	}

	var session *mgo.Session
	if s, err := mgo.Dial(*mongoAddr); err != nil {
		log.Fatal(err)
	} else {
		session = s
	}

	start := time.Now()

	splitColl(session, *mongoDBHourAggreagate, *mongoCollHourAggregate, make(map[string]*mgo.Collection))
	log.Debug("hour splitred", time.Since(start))
	start = time.Now()
	splitColl(session, *mongoDBDayAggreagate, *mongoCollDayAggregate, make(map[string]*mgo.Collection))
	log.Debug("day splitred", time.Since(start))
	start = time.Now()
	splitColl(session, *mongoDBTotalAggreagate, *mongoCollTotalAggregate, make(map[string]*mgo.Collection))
	log.Debug("total splitred", time.Since(start))

	log.Debug("Finsihed! waiting for exit signal")
	select {}
}

func splitColl(session *mgo.Session, dbName, collName string, colls map[string]*mgo.Collection) {
	db := session.DB(dbName)
	log.Debug("Start split db:", db.Name, ",  collection:", collName)
	var e aggregate.CounterEvent
	c := db.C(collName)
	log.Debug(c.Find(bson.M{}).Count())
	iter := c.Find(bson.M{}).Iter()

	for iter.Next(&e) {
		reg, err := regexp.Compile("^[a-zA-Z0-9]+$")
		if err != nil {
			log.Fatal(err)
		}
		if reg.Match([]byte(e.AppID)) {
			coll := aggregate.GetColl(session, dbName, colls, collName, e.AppID)
			if err := coll.Insert(e); err != nil {
				log.Fatal(err, e.AppID)
			}
		}
	}
}
