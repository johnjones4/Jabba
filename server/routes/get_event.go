package routes

import (
	"context"
	"main/store"

	"github.com/johnjones4/Jabba/core"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type getEventInput struct {
	ID int `path:"id"`
}

func GetEventUseCase(s store.Store) usecase.IOInteractor {
	return usecase.NewIOI(new(getEventInput), new(core.Event), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*getEventInput)
			out = output.(*core.Event)
		)

		event, err := s.GetEvent(in.ID)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*out = event

		return nil
	})
}
