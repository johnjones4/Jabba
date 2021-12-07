package routes

import (
	"context"
	"main/store"

	"github.com/johnjones4/Jabba/core"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type getEventsInput struct {
	Offset          int    `query:"offset"`
	Limit           int    `query:"limit"`
	EventVendorType string `query:"eventVendorType"`
}

type getEventsOutput struct {
	Items []core.Event `json:"items"`
}

func GetEventsUseCase(s store.Store) usecase.IOInteractor {
	return usecase.NewIOI(new(getEventsInput), new(getEventsOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*getEventsInput)
			out = output.(*getEventsOutput)
		)

		var events []core.Event
		var err error
		if in.EventVendorType != "" {
			events, err = s.GetEventsForVendorType(in.EventVendorType, in.Limit, in.Offset)
		} else {
			events, err = s.GetEvents(in.Limit, in.Offset)
		}
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		out.Items = events

		return nil
	})
}
