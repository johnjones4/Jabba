package routes

import (
	"context"
	statusEngine "main/status"

	"github.com/johnjones4/Jabba/core"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type getEventVendorTypeInfoInput struct {
	EventVendorType string `path:"eventVendorType"`
}

func GetEventVendorTypeInfoUseCase(se statusEngine.StatusEngine) usecase.IOInteractor {
	return usecase.NewIOI(new(getEventVendorTypeInfoInput), new(core.Status), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*getEventVendorTypeInfoInput)
			out = output.(*core.Status)
		)

		s, err := se.ProcessEventsForVendorType(in.EventVendorType)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*out = *s

		return nil
	})
}
