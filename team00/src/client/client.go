package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math"
	"sync"

	"team00/db"
	pb "team00/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

func main() {
	coef := flag.Float64("k", 0.0, "Установка среднеквадратичного разброса коэфициэнта аномалии")
	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial("localhost:3333", opts...)
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	client := pb.NewTransmitterClient(conn)
	stream, err := client.ListRequests(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу %v", err)
	}

	psql, err := db.Connect()
	if err != nil {
		log.Fatalf("Ошибка при подключении к бд %v", err)
	}

	detectAnomalies(stream, psql, *coef)
}

func calcMean(elem float64, mean float64, count int) float64 {
	return (mean*float64(count-1) + elem) / float64(count)
}

func calcSd(mean float64, pool *sync.Pool) (float64, *sync.Pool) {
	sum := 0.0
	count := 0

	newPool := sync.Pool{}
	for elem := pool.Get(); elem != nil; elem = pool.Get() {
		sum += math.Pow(elem.(float64)-mean, 2.0)
		newPool.Put(elem)
		count++
	}

	return math.Sqrt(sum / float64(count)), &newPool
}

func detectAnomalies(stream pb.Transmitter_ListRequestsClient, psql *gorm.DB, coef float64) {
	mean := 0.0
	count := 0
	sd := 0.0
	anomalyCount := 0

	pool := new(sync.Pool)
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("Ошибка при подключении к серверу %v", err)
		}

		count++

		if count < 150 {
			pool.Put(response.Frequency)
			mean = calcMean(response.Frequency, mean, count)
			sd, pool = calcSd(mean, pool)

			log.Printf("После %v ответов среднее отклонение %v, стандартное отклонение %v\n", count, mean, sd)
		}

		if count >= 150 {
			low := mean - (coef * sd)
			high := mean + (coef * sd)

			if !(response.Frequency >= low && response.Frequency <= high) {
				log.Printf("Аномалия обнаружена в ответе %v; частота %v; минимум %v, максимум %v\n", count, response.Frequency, low, high)

				anomalyCount++

				log.Println(anomalyCount, float64(anomalyCount)/float64(count)*100.0)

				psql.Create(&db.Record{
					Uuid:      response.Uuid,
					Frequency: response.Frequency,
					Timestamp: response.Timestamp,
				})
			}
		}
	}
}
