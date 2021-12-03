package routes

import (
	"context"
	"errors"
	"main/store"
	"time"

	"github.com/johnjones4/Jabba/core"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

const (
	StatusOk         = "ok"
	StatusRecovering = "recovering"
	StatusAbnormal   = "abnormal"
)

type getEventVendorTypeInfoInput struct {
	EventVendorType string `path:"eventVendorType"`
}

type getEventVendorTypeInfoOutput struct {
	EventVendorType string     `json:"eventVendorType"`
	EventVendorName string     `json:"eventVendorName"`
	Status          string     `json:"status"`
	LastEvent       core.Event `json:"lastEvent"`
}

func GetEventVendorTypeInfoUseCase(s store.Store, vendorInfo map[string]string) usecase.IOInteractor {
	return usecase.NewIOI(new(getEventVendorTypeInfoInput), new(getEventVendorTypeInfoOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*getEventVendorTypeInfoInput)
			out = output.(*getEventVendorTypeInfoOutput)
		)

		out.EventVendorType = in.EventVendorType
		if name, ok := vendorInfo[in.EventVendorType]; ok {
			out.EventVendorName = name
		} else {
			out.EventVendorName = in.EventVendorType
		}

		events, err := s.GetEventsForVendorType(in.EventVendorType, 2, 0)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		if len(events) == 0 {
			return status.Wrap(errors.New("not enough context for status"), status.Internal)
		}

		out.LastEvent = events[0]

		oneWeekAgo := time.Hour * 24 * 7 * -1

		if !out.LastEvent.IsNormal || out.LastEvent.Created.Before(time.Now().UTC().Add(oneWeekAgo)) {
			out.Status = StatusAbnormal
			return nil
		}

		if len(events) > 1 && !events[1].IsNormal {
			out.Status = StatusRecovering
			return nil
		}

		out.Status = StatusOk

		return nil
	})
}
