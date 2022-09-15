package dis

import (
	"net/http"

	"github.com/tkw1536/goprogram/stream"
)

func (dis Dis) info(io stream.IOStream) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		all, err := dis.Instances.All()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			io.EPrintln(err)
			return
		}

		for _, wk := range all {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(wk.Slug))
			w.Write([]byte("\n"))
		}

	}), nil
}
