package routes

import (
	"context"
	"main/shared"
	"main/store"

	"github.com/johnjones4/Jabba/core"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
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

		lastEvent, infoStatus, err := shared.GetEventVendorTypeInfo(s, in.EventVendorType)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		out.LastEvent = lastEvent
		out.Status = infoStatus

		return nil
	})
}
