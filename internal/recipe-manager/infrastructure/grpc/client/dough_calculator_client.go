package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	bdomain "github.com/cfioretti/recipe-manager/internal/ingredients-balancer/domain"
	pb "github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/grpc/proto/generated"
)

type DoughCalculatorClient struct {
	client     pb.DoughCalculatorClient
	conn       *grpc.ClientConn
	serverAddr string
	timeout    time.Duration
}

func NewDoughCalculatorClient(serverAddr string, timeout time.Duration) (*DoughCalculatorClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewDoughCalculatorClient(conn)

	return &DoughCalculatorClient{
		client:     client,
		conn:       conn,
		serverAddr: serverAddr,
		timeout:    timeout,
	}, nil
}

func (c *DoughCalculatorClient) Close() error {
	return c.conn.Close()
}

func (c *DoughCalculatorClient) TotalDoughWeightByPans(pans bdomain.Pans) (*bdomain.Pans, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	protoRequest := toProtoMessage(&pans)
	response, err := c.client.TotalDoughWeightByPans(ctx, &pb.PansRequest{
		Pans: protoRequest,
	})
	if err != nil {
		return nil, err
	}

	result := toDomainPans(response.Pans)
	return &result, nil
}

func toProtoMessage(domainPans *bdomain.Pans) *pb.PansProto {
	panProtos := make([]*pb.PanProto, 0, len(domainPans.Pans))

	for _, p := range domainPans.Pans {
		panProto := &pb.PanProto{
			Shape: p.Shape,
			Measures: &pb.MeasuresProto{
				Diameter: fromPointer(p.Measures.Diameter),
				Edge:     fromPointer(p.Measures.Edge),
				Width:    fromPointer(p.Measures.Width),
				Length:   fromPointer(p.Measures.Length),
			},
			Name: p.Name,
			Area: p.Area,
		}
		panProtos = append(panProtos, panProto)
	}

	return &pb.PansProto{
		Pans:      panProtos,
		TotalArea: domainPans.TotalArea,
	}
}

func toDomainPans(protoMessage *pb.PansProto) bdomain.Pans {
	pans := make([]bdomain.Pan, 0, len(protoMessage.Pans))

	for _, p := range protoMessage.Pans {
		pan := bdomain.Pan{
			Shape: p.Shape,
			Measures: bdomain.Measures{
				Diameter: toPointer(p.Measures.Diameter),
				Edge:     toPointer(p.Measures.Edge),
				Width:    toPointer(p.Measures.Width),
				Length:   toPointer(p.Measures.Length),
			},
			Name: p.Name,
			Area: p.Area,
		}
		pans = append(pans, pan)
	}

	return bdomain.Pans{
		Pans:      pans,
		TotalArea: protoMessage.TotalArea,
	}
}

func toPointer(value *int32) *int {
	if value == nil {
		return nil
	}
	val := int(*value)
	return &val
}

func fromPointer(value *int) *int32 {
	if value == nil {
		return nil
	}
	val := int32(*value)
	return &val
}
