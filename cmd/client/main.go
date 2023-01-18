package main

import (
	"context"
	"crud_grpc/pb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewCarClient(conn)

	runGetCars(client)
	// 	runGetCarById(client, "2")
	// 	runCreateCar(client, "TES123", "CARRO VELHO", 1900, 1900, "CLIENTE TESTE", "95478654123")
	// 	runDeleteCar(client, "2")
	// 	runGetCars(client)
	// 	runUpdateCar(client, "3", "TES123", "CARRO RELIQUIA", 1900, 1900, "CLIENTE TESTE", "95478654123")
	// 	runGetCars(client)
}

func runGetCars(client pb.CarClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Empty{}
	stream, err := client.GetCars(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetCars(_) = _, %v", client, err)
	}
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCars(_) = _, %v", client, err)
		}
		log.Printf("CarInfo: %v", row)
	}
}
func runGetCarById(client pb.CarClient, carid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Id: carid}
	carInfo, err := client.GetCarById(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetCarById(id) = _, %v", client, err)
	}
	log.Printf("CarInfo: %v", carInfo)
}
func runCreateCar(client pb.CarClient, placa string, modelo string, anofabricacao int32, anomodelo int32, nomeCliente string, cpf string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	carinfo := pb.CarInfo{Placa: placa,
		Modelo:        modelo,
		Anofabricacao: anofabricacao,
		Anomodelo:     anomodelo,
		Client: &pb.Client{
			Name: nomeCliente,
			Cpf:  cpf},
	}
	carid, err := client.CreatCar(ctx, &carinfo)
	if err != nil {
		log.Fatalf("%v.CreatCar(CarInfo) = _, %v", client, err)
	}
	log.Printf("CarId: %v", carid)
}
func runUpdateCar(client pb.CarClient, id string, placa string, modelo string, anofabricacao int32, anomodelo int32, nomeCliente string, cpf string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	carinfo := pb.CarInfo{Id: id,
		Placa:         placa,
		Modelo:        modelo,
		Anofabricacao: anofabricacao,
		Anomodelo:     anomodelo,
		Client: &pb.Client{
			Name: nomeCliente,
			Cpf:  cpf},
	}
	carid, err := client.UpdateCar(ctx, &carinfo)
	if err != nil {
		log.Fatalf("%v.CreatCar(CarInfo) = _, %v", client, err)
	}
	log.Printf("CarId: %v", carid)
}
func runDeleteCar(client pb.CarClient, carid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Id: carid}
	status, err := client.DeleteCar(ctx, req)
	if err != nil {
		log.Fatalf("%v.DeleteCar(id) = _, %v", client, err)
	}
	if int(status.GetStatus()) == 1 {
		log.Printf("DeleteCar Success")
	} else {
		log.Printf("DeleteCar Failed")
	}
}
