package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	pb2 "learn-protobuf/0/pb"
	"log"
)

func main() {
	employee := &pb2.Employee{
		Id:          1,
		Name:        "Suzuki",
		Email:       "test@example.com",
		Occupation:  pb2.Occupation_ENGINEER,
		PhoneNumber: []string{"080-1234-5678", "090-1234-5678"},
		Project:     map[string]*pb2.Company_Project{"ProjectX": &pb2.Company_Project{}},
		Profile:     &pb2.Employee_Text{Text: "My name is Suzuki."},
		Birthday:    &pb2.Date{Year: 2000, Month: 1, Day: 1},
	}

	binData, err := proto.Marshal(employee)
	if err != nil {
		log.Fatalln("can't serialize", err)
	}

	if err := ioutil.WriteFile("test.bin", binData, 0666); err != nil {
		log.Fatalln("can't write", err)
	}

	in, err := ioutil.ReadFile("test.bin")
	if err != err {
		log.Fatalln("can't read file", err)
	}
	readEmployee := &pb2.Employee{}
	if err = proto.Unmarshal(in, readEmployee); err != nil {
		log.Fatalln("can't deserialize", err)
	}
	fmt.Println(readEmployee)

}
