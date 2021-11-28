package routes

import (
	"context"
	"main/job"
	"main/upstream"

	"github.com/johnjones4/Jabba/core"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func SaveJobUsecase(g job.AlertGenerator, u upstream.Upstream) usecase.IOInteractor {
	return usecase.NewIOI(new(core.Event), new(core.Event), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*core.Event)
			out = output.(*core.Event)
		)

		err := g.GenerateAlerts(in)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		err = u.LogEvent(in)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}

		*out = *in

		return nil
	})
}
