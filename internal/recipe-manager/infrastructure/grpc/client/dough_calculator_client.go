package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/cfioretti/recipe-manager/internal/infrastructure/logging"
	"github.com/cfioretti/recipe-manager/internal/recipe-manager/domain"
	pb "github.com/cfioretti/recipe-manager/internal/recipe-manager/infrastructure/grpc/proto/generated"
)

type CalculatorClient struct {
	client     pb.DoughCalculatorClient
	conn       *grpc.ClientConn
	serverAddr string
	timeout    time.Duration
}

func NewDoughCalculatorClient(serverAddr string, timeout time.Duration) (*CalculatorClient, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewDoughCalculatorClient(conn)

	return &CalculatorClient{
		client:     client,
		conn:       conn,
		serverAddr: serverAddr,
		timeout:    timeout,
	}, nil
}

func (c *CalculatorClient) Close() error {
	return c.conn.Close()
}

func (c *CalculatorClient) TotalDoughWeightByPans(ctx context.Context, pans domain.Pans) (*domain.Pans, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	correlationID := logging.GetCorrelationID(ctx)
	md := metadata.Pairs("x-correlation-id", correlationID)
	timeoutCtx = metadata.NewOutgoingContext(timeoutCtx, md)

	protoRequest := toProtoMessage(&pans)
	response, err := c.client.TotalDoughWeightByPans(timeoutCtx, &pb.PansRequest{
		Pans: protoRequest,
	})
	if err != nil {
		return nil, err
	}

	result := toDomainPans(response.Pans)
	return &result, nil
}

func toProtoMessage(domainPans *domain.Pans) *pb.PansProto {
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

func toDomainPans(protoMessage *pb.PansProto) domain.Pans {
	pans := make([]domain.Pan, 0, len(protoMessage.Pans))

	for _, p := range protoMessage.Pans {
		pan := domain.Pan{
			Shape: p.Shape,
			Measures: domain.Measures{
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

	return domain.Pans{
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
