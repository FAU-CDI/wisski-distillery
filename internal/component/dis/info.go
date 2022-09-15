package dis

import (
	"encoding/json"
	"net/http"

	"github.com/FAU-CDI/wisski-distillery/internal/component/instances"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/sync/errgroup"
)

func (dis *Dis) info(io stream.IOStream) (http.Handler, error) {
	return http.HandlerFunc(dis.handleDis), nil
}

const disLimit = 2

func (dis *Dis) handleDis(w http.ResponseWriter, r *http.Request) {
	// make sure the user is authorized
	if !dis.authDis(r) {
		w.Header().Add("WWW-Authenticate", `Basic realm="WissKI Distillery Admin"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	// create a new error group
	var errgroup errgroup.Group
	errgroup.SetLimit(disLimit)

	// list all the instances
	all, err := dis.Instances.All()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	// get all of their info!
	infos := make([]instances.Info, len(all))
	for i, instance := range all {
		{
			i := i
			instance := instance

			errgroup.Go(func() (err error) {
				infos[i], err = instance.Info()
				return err
			})
		}
	}

	// if some info call failed
	if err := errgroup.Wait(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		w.Write([]byte("\n"))
		return
	}

	// and return the json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(infos)
}

func (dis *Dis) authDis(r *http.Request) bool {
	user, pass, ok := r.BasicAuth()
	return ok && user == dis.Config.DisAdminUser && pass == dis.Config.DisAdminPassword
}
