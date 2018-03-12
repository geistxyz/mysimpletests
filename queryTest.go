// ESTE es un mensaje introducido en la rama RELEASE1.2 
// VEREMOS QUE PASA! 
package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)



type msg struct {
	// Msg   string        `bson:"msg"`
	Count int `bson:"count"`
}

// Note: attribute name must be upper-case start. Or it will not save to DB. (could not identical with document (at least one upper-case))
type Book struct {
	ISBN  string
	TITLE string
	PRICE int
}

type Person struct {
	Name  string
	Phone string
}
type CardDetail struct {/*CMN10021F*/
	Id           bson.ObjectId  `json:"id" bson:"_id"`
	CollectorId   string         `json:"Collector_ID"  bson:"Collector ID"`
	ComplianceCondition       string           `json:"Compliance_Condition"  bson:"Compliance Condition"`
	Rows         int      `json:"Rows_For_Compliance"  bson:"Number of Rows For Compliance"`
	RowsLastRun      int      `json:"Number_of_Rows_In_Last_Run"  bson:"Number of Rows In Last Run"`
	RegulationDetails     string         `json:"Regulation_Details"  bson:"Regulation Details"`
	ReportCardId       string            `json:"Report_Card_ID"  bson:"Report Card ID"`
	ReportId      string            `json:"Report_ID"  bson:"Report ID"`
	DefaultValues      interface{}      `json:"Defaults_Values"  bson:"Defaults Values"`
	Exceptions      interface{}      `json:"Exceptions"  bson:"Exceptions"`
}
type CardDetails []CardDetail

type Card struct {/*CMN10020F*/
	Id           bson.ObjectId  `json:"id" bson:"_id"`
	CardName   string         `json:"Report_Card_Name"  bson:"Report Card Name"`
	CardId       string           `json:"Report_Card_ID"  bson:"Report Card ID"`
	Date         time.Time      `json:"date"  bson:"date"`
	Reports      interface{}      `json:"reports"  bson:"reports"`
	Details      CardDetails      `json:"details"  bson:"details"`
	Category     string         `json:"Category"  bson:"Category"`
	Reload       string            `json:"reload"  bson:"reload"`
	Ovrdfts      string            `json:"ovrdfts"  bson:"ovrdfts"`
	Regulation     string         `json:"Regulation"  bson:"Regulation ?"`
	LastRunStart     string         `json:"Last_Run_Start"  bson:"Last Run Start"`
	LastRunEnd     string         `json:"Last_Run_End"  bson:"Last Run End"`
	ComplianceThreshold     int         `json:"Compliance_Threshold"  bson:"Compliance Threshold"`
	TrinityGuardGroup     string         `json:"Trinity_Guard_Group"  bson:"Trinity Guard Group"`
}

type Report struct {
	Id           bson.ObjectId  `json:"id" bson:"_id"`
	Title   string         `json:"title"  bson:"title"`
	Server       string           `json:"server"  bson:"server"`
	Date         string      `json:"date"  bson:"date"`
	DateCreate         time.Time      `json:"date_create"  bson:"date_create"`
	Header         interface{}      `json:"header"  bson:"header"`
	Data         interface{}      `json:"data"  bson:"data"`
	ReportId     string         `json:"report_id"  bson:"report_id"`
	QueueId       string            `json:"queue_id"  bson:"queue_id"`
	UserId       int            `json:"user_id"  bson:"user_id"`
	Size       int            `json:"size"  bson:"size"`
	PageNumber       int            `json:"page_number"  bson:"page_number"`
	PageTotal       int            `json:"page_total"  bson:"page_total"`
	CardId     string         `json:"report_card_id"  bson:"report_card_id"`
	CollectorId       string            `json:"collector_id"  bson:"collector_id"`
	Error       string            `json:"error"  bson:"error"`
}

type Reports []Report


func main() {
	fmt.Println("Starting connect mongoDB....")
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	fmt.Println("Connection works....")
	// Optional. Switch the session to a monotonic behavior.
	session.SetSafe(&mgo.Safe{})
	c := session.DB("tg").C("report")
	var page  =1
	var id = "59fbea1dec006475c44645bd"
	iter := c.Find(bson.M{"report_card_id": bson.M{"$eq": id} , "page_number":page }).Iter()
	fmt.Println("Getting  Database count....")
	count, err2 := c.Find(bson.M{}).Count()
	if err2 != nil {
		panic(err)
	}
	fmt.Printf("total report count = %d\n", count)

	/*if count == 0 {
		err = c.Insert(&Book{"Ale1", "Book1", 35})
		err = c.Insert(&Book{"Ale2", "Book2", 20})
		err = c.Insert(&Book{"Ale3", "Book3", 40})
		err = c.Insert(&Book{"Ale4", "Book4", 15})
		err = c.Insert(&Book{"Ale5", "Book5", 55})
		err = c.Insert(&Book{"Ale6", "Book6", 45})
	}*/

	result := Report{}
	fmt.Println("Getting  data....")
	// Find book which prices is greater than(gt) 40
	// iter := c.Find(bson.M{"Size": bson.M{"$gte": 0}}).Iter()
	// iter := c.Find(nil).Iter()
	var index = 1
	reports := Reports{}
	for iter.Next(&result) {
		fmt.Printf("current result is [%d] result =%+v\n", index, result)
		index++
		reports = append(reports, result)
	}

	for i := 0; i < len(reports); i++ {
		fmt.Println("reporte ")
		fmt.Println(reports[i].Title)
	}
	//when search the DB it must all lower-case to avoid any error.
	if err2 := iter.Close(); err2 != nil {
		fmt.Printf("No data\n")
	} else {
		fmt.Printf("result =%+v\n", result)
	}
}
