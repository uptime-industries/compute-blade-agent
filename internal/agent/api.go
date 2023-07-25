package agent

import (
	"context"

	bladeapiv1alpha1 "github.com/xvzf/computeblade-agent/api/bladeapi/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ComputeBladeAgent implementing the BladeAgentServiceServer
type agentGrpcService struct {
	bladeapiv1alpha1.UnimplementedBladeAgentServiceServer

	Agent ComputeBladeAgent
}

// NewGrpcServiceFor creates a new gRPC service for a given agent
func NewGrpcServiceFor(agent ComputeBladeAgent) *agentGrpcService {
	return &agentGrpcService{
		Agent: agent,
	}
}

// EmitEvent emits an event to the agent runtime
func (service *agentGrpcService) EmitEvent(
	ctx context.Context,
	req *bladeapiv1alpha1.EmitEventRequest,
) (*emptypb.Empty, error) {
	switch req.GetEvent() {
	case bladeapiv1alpha1.Event_IDENTIFY:
		return &emptypb.Empty{}, service.Agent.EmitEvent(ctx, IdentifyEvent)
	case bladeapiv1alpha1.Event_IDENTIFY_CONFIRM:
		return &emptypb.Empty{}, service.Agent.EmitEvent(ctx, IdentifyConfirmEvent)
	case bladeapiv1alpha1.Event_CRITICAL:
		return &emptypb.Empty{}, service.Agent.EmitEvent(ctx, CriticalEvent)
	case bladeapiv1alpha1.Event_CRITICAL_RESET:
		return &emptypb.Empty{}, service.Agent.EmitEvent(ctx, CriticalResetEvent)
	default:
		return &emptypb.Empty{}, status.Errorf(codes.InvalidArgument, "invalid event type")
	}
}

func (service *agentGrpcService) WaitForIdentifyConfirm(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, service.Agent.WaitForIdentifyConfirm(ctx)
}

// SetFanSpeed sets the fan speed of the blade
func (service *agentGrpcService) SetFanSpeed(
	ctx context.Context,
	req *bladeapiv1alpha1.SetFanSpeedRequest,
) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, service.Agent.SetFanSpeed(ctx, uint8(req.GetPercent()))
}

// SetStealthMode enables/disables stealth mode on the blade
func (service *agentGrpcService) SetStealthMode(ctx context.Context, req *bladeapiv1alpha1.StealthModeRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, service.Agent.SetStealthMode(ctx, req.GetEnable())
}

// GetStatus aggregates the status of the blade
func (service *agentGrpcService) GetStatus(context.Context, *emptypb.Empty) (*bladeapiv1alpha1.StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
