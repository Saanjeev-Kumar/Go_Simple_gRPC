package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"simpleGrpc/protoPackage"
)

type AddressBook struct {
	People map[int32]*protoPackage.Student
}

func (a *AddressBook) AddPerson(ctx context.Context, req *protoPackage.AddPersonReq) (*protoPackage.Empty, error) {
	//if req.Student.Id == 0 {
	//	return &protoPackage.Empty{}, fmt.Errorf("student cannot be empty")
	//}
	//
	//if _, ok := a.People[req.Student.Id]; ok {
	//	return &protoPackage.Empty{}, fmt.Errorf("student already exist %d", req.Student.Id)
	//}
	a.People[req.Student.Id] = req.Student
	return &protoPackage.Empty{}, nil
}

func (a *AddressBook) GetStudentDetails(ctx context.Context, req *protoPackage.GetDetailsReq) (*protoPackage.Student, error) {
	student, _ := a.People[req.Id]
	//student, ok := a.People[req.Id]
	//if !ok {
	//	return nil, fmt.Errorf("Student with that doesn't exist %d", req.Id)
	//}
	return student, nil

}

func StartGrpcServer() {
	s := grpc.NewServer()
	//&AddressBook{People: make(map[int32]*protoPackage.Student, 0)}
	protoPackage.RegisterAddressBookServer(s, &AddressBook{People: make(map[int32]*protoPackage.Student, 0)})
	//	protoPackage.RegisterAddressBookServer(s, &AddressBook{People: make(map[int32]*protoPackage.Student, 0)})
	//&AddressBook{People: make(map[int32]*protoPackage.Student, 0)}

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("failed to listen to grpc port", err)
		return
	}

	if err = s.Serve(lis); err != nil {
		fmt.Println("failed tp serve GRPC", err)
		return
	}

}

func main() {
	go StartGrpcServer()

	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		fmt.Println("failed to dial grpc server", err)
	}
	client := protoPackage.NewAddressBookClient(conn)
	_, err = client.AddPerson(context.Background(), &protoPackage.AddPersonReq{
		Student: &protoPackage.Student{
			Id:    11,
			Name:  "saanjeevkuamr",
			Email: "saanjeevkumar@mail.com",
		},
	})
	_, err = client.AddPerson(context.Background(), &protoPackage.AddPersonReq{
		Student: &protoPackage.Student{
			Id:    12,
			Name:  "saanjeev",
			Email: "saanjeev@mail.com",
		},
	})
	//_, err = client.AddPerson(context.Background(), &protoPackage.Student{
	//		Id:    12,
	//		Name:  "saanjeev",
	//		Email: "saanjeev@mail.com",
	//	},
	//)
	if err != nil {
		fmt.Println("Error in creating a student details")
		return
	}
	person, err := client.GetStudentDetails(context.Background(), &protoPackage.GetDetailsReq{Id: 11})
	person1, err := client.GetStudentDetails(context.Background(), &protoPackage.GetDetailsReq{Id: 12})
	if err != nil {
		fmt.Println("error in retrival of student deatils", err)
		return
	}
	fmt.Printf("Retrived Student Id : %d, name:%s, email :%s \n", person.Id, person.Name, person.Email)
	fmt.Printf("Retrived Student Id : %d, name:%s, email :%s\n", person1.Id, person1.Name, person1.Email)
}
