package main

import (
	"context"
	"crud_grpc/pb"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var cars []*pb.CarInfo

type Server struct {
	pb.UnimplementedCarServer
}

func main() {
	println("Running gRPC server")

	initCar()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(err)
	}

	s := grpc.NewServer()

	pb.RegisterCarServer(s, &Server{})

	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initCar() {
	car1 := &pb.CarInfo{
		Id:            "1",
		Placa:         "ABC1235",
		Modelo:        "BMW",
		Anofabricacao: 2021,
		Anomodelo:     2022,
		Client: &pb.Client{
			Name: "Lucas Euzebio",
			Cpf:  "01234567890",
		},
	}
	car2 := &pb.CarInfo{
		Id:            "2",
		Placa:         "DEV1235",
		Modelo:        "FUSCA",
		Anofabricacao: 1885,
		Anomodelo:     1885,
		Client: &pb.Client{
			Name: "Harley Mari",
			Cpf:  "45678901234",
		},
	}

	cars = append(cars, car1)
	cars = append(cars, car2)
}
func (s *Server) GetCars(in *pb.Empty,
	stream pb.Car_GetCarsServer) error {
	log.Printf("receivedGet: %v", in)
	for _, car := range cars {
		if err := stream.Send(car); err != nil {
			return err
		}
	}
	return nil
}
func (s *Server) GetCarById(ctx context.Context, in *pb.Id) (*pb.CarInfo, error) {
	log.Printf("receivedBy: %v", in)

	res := &pb.CarInfo{}

	for _, car := range cars {
		if car.GetId() == in.GetId() {
			res = car
			break
		}
	}
	return res, nil
}
func (s *Server) CreatCar(ctx context.Context, in *pb.CarInfo) (*pb.Id, error) {
	log.Printf("receivedCreate: %v", in)

	res := &pb.Id{}
	in.Id = strconv.Itoa(len(cars) + 1)
	cars = append(cars, in)
	res.Id = in.GetId()
	return res, nil
}
func (s *Server) UpdateCar(ctx context.Context, in *pb.CarInfo) (*pb.Status, error) {
	log.Printf("receivedDel: %v", in)

	res := &pb.Status{}

	for _, car := range cars {
		if car.GetId() == in.GetId() {
			car.Anofabricacao = in.Anofabricacao
			car.Anomodelo = in.Anomodelo
			car.Modelo = in.GetModelo()
			car.Placa = in.GetPlaca()
			car.Client.Cpf = in.Client.GetCpf()
			car.Client.Name = in.Client.GetName()
			res.Status = 1
			break
		}
	}
	return res, nil
}
func (s *Server) DeleteCar(ctx context.Context, in *pb.Id) (*pb.Status, error) {
	log.Printf("receivedDel: %v", in)

	res := &pb.Status{}

	for index, car := range cars {
		if car.GetId() == in.GetId() {
			cars = append(cars[:index], cars[index+1:]...)
			res.Status = 1
			break
		}
	}

	return res, nil
}
