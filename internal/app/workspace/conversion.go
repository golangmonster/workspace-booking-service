package workspace

import (
	"github.com/golangmonster/workspace-booking-service/internal/model/page"
	model "github.com/golangmonster/workspace-booking-service/internal/model/workspace"
	dto "github.com/golangmonster/workspace-booking-service/internal/service/workspace"
	pb "github.com/golangmonster/workspace-booking-service/pkg/api/workspace/v1"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	workspaceTypeToModel = map[pb.WorkspaceType]model.Type{
		pb.WorkspaceType_WORKSPACE_TYPE_UNSPECIFIED:  model.TypeUnspecified,
		pb.WorkspaceType_WORKSPACE_TYPE_DESK:         model.TypeDesk,
		pb.WorkspaceType_WORKSPACE_TYPE_MEETING_ROOM: model.TypeMeetingRoom,
		pb.WorkspaceType_WORKSPACE_TYPE_PHONE_BOOTH:  model.TypePhoneBooth,
	}

	workspaceTypeToProto = map[model.Type]pb.WorkspaceType{
		model.TypeUnspecified: pb.WorkspaceType_WORKSPACE_TYPE_UNSPECIFIED,
		model.TypeDesk:        pb.WorkspaceType_WORKSPACE_TYPE_DESK,
		model.TypeMeetingRoom: pb.WorkspaceType_WORKSPACE_TYPE_MEETING_ROOM,
		model.TypePhoneBooth:  pb.WorkspaceType_WORKSPACE_TYPE_PHONE_BOOTH,
	}

	workspaceStatusToModel = map[pb.WorkspaceStatus]model.Status{
		pb.WorkspaceStatus_WORKSPACE_STATUS_UNSPECIFIED: model.StatusUnspecified,
		pb.WorkspaceStatus_WORKSPACE_STATUS_AVAILABLE:   model.StatusAvailable,
		pb.WorkspaceStatus_WORKSPACE_STATUS_BOOKED:      model.StatusBooked,
		pb.WorkspaceStatus_WORKSPACE_STATUS_MAINTENANCE: model.StatusMaintenance,
	}

	workspaceStatusToProto = map[model.Status]pb.WorkspaceStatus{
		model.StatusUnspecified: pb.WorkspaceStatus_WORKSPACE_STATUS_UNSPECIFIED,
		model.StatusAvailable:   pb.WorkspaceStatus_WORKSPACE_STATUS_AVAILABLE,
		model.StatusBooked:      pb.WorkspaceStatus_WORKSPACE_STATUS_BOOKED,
		model.StatusMaintenance: pb.WorkspaceStatus_WORKSPACE_STATUS_MAINTENANCE,
	}
)

func toCreateWorkspaceRequest(req *pb.CreateWorkspaceRequest) *dto.CreateWorkspaceRequest {
	return &dto.CreateWorkspaceRequest{
		Name:        req.GetName(),
		Lat:         req.GetCoordinates().GetLat(),
		Lon:         req.GetCoordinates().GetLon(),
		FullAddress: req.GetFullAddress(),
		Type:        workspaceTypeToModel[req.GetType()],
		Capacity:    req.GetCapacity(),
	}
}

func toListWorkspacesRequest(req *pb.ListWorkspacesRequest) *dto.ListWorkspacesRequest {
	var (
		filter *dto.WorkspaceFilter
	)

	if req.Filter != nil {
		filter = &dto.WorkspaceFilter{
			Types: lo.Map(req.Filter.Types, func(t pb.WorkspaceType, _ int) model.Type {
				return workspaceTypeToModel[t]
			}),
			Statuses: lo.Map(req.Filter.Statuses, func(t pb.WorkspaceStatus, _ int) model.Status {
				return workspaceStatusToModel[t]
			}),
			CapacityGte: req.Filter.CapacityGte,
			FullAddress: req.Filter.FullAddress,
		}
	}

	return &dto.ListWorkspacesRequest{
		Page: &page.Page{
			Limit:  req.GetPage().GetLimit(),
			Offset: req.GetPage().GetOffset(),
		},
		Filter: filter,
	}
}

func toListWorkspacesResponse(resp *dto.ListWorkspacesResponse) *pb.ListWorkspacesResponse {
	return &pb.ListWorkspacesResponse{
		Workspaces: lo.Map(resp.Workspaces, func(w *model.Workspace, _ int) *pb.ListWorkspacesResponse_Workspace {
			return &pb.ListWorkspacesResponse_Workspace{
				Id:          w.ID,
				Name:        w.Name,
				Coordinates: toLatLon(w.Lat, w.Lon),
				FullAddress: w.FullAddress,
				Type:        workspaceTypeToProto[w.Type],
				Status:      workspaceStatusToProto[w.Status],
				Capacity:    w.Capacity,
				CreatedAt:   timestamppb.New(w.CreatedAt),
				UpdatedAt:   timestamppb.New(w.UpdatedAt),
			}
		}),
		TotalCount: resp.TotalCount,
	}
}

func toLatLon(lat, lon float64) *pb.LatLon {
	return &pb.LatLon{Lat: lat, Lon: lon}
}
