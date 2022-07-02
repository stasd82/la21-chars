package mid

import (
	"context"
	"net/http"

	"github.com/stasd82/la21-chars/domain/sys/validate"
	v1 "github.com/stasd82/la21-chars/domain/web/v1"
	"github.com/stasd82/tux"
	"go.uber.org/zap"
)

func Errors(sl *zap.SugaredLogger) (c tux.Circuit) {

	c = func(in tux.Route) (out tux.Route) {

		out = func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			v, err := tux.GetValues(ctx)
			if err != nil {
				return tux.NewShutdownError("web values missing from context")
			}

			// Run the next handler and catch any propagated error.
			if err = in(ctx, w, r); err != nil {

				// Log the error.
				sl.Errorw("ERROR", "traceid", v.TraceID, "ERROR", err)

				// Build out the error response.
				var er v1.ErrorResponse
				var status int
				switch {
				case validate.IsFieldErrors(err):
					status = http.StatusBadRequest
					er = v1.ErrorResponse{
						Error:  "data validation error",
						Fields: validate.GetFieldErrors(err).Fields(),
					}
				case v1.IsRequestErr(err):
					re := v1.GetRequestErr(err)
					status = re.Status
					er = v1.ErrorResponse{
						Error: re.Error(),
					}
				default:
					status = http.StatusInternalServerError
					er = v1.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
				}

				// Respond with the error back to the client
				if err := tux.Respond(ctx, w, er, status); err != nil {
					return err
				}

				if tux.IsShutdown(err) {
					return err
				}
			}
			// The error has been handled so we can stop propagating it.
			return nil
		}
		return
	}
	return
}
