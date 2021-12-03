package routes

import (
	"context"
	"main/store"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type getEventVendorTypesInput struct {
}

type getEventVendorTypesOutput struct {
	EventVendorTypes []string `json:"eventVendorTypes"`
}

func GetEventVendorTypesUseCase(s store.Store) usecase.IOInteractor {
	return usecase.NewIOI(new(getEventVendorTypesInput), new(getEventVendorTypesOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*getEventVendorTypesOutput)
		)

		eventVendorTypes, err := s.GetEventVendorTypes()
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		out.EventVendorTypes = eventVendorTypes

		return nil
	})
}
