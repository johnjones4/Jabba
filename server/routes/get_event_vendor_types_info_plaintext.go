package routes

import (
	"log"
	"main/shared"
	"main/store"
	"net/http"
	"strings"
)

func GetEventVendorTypesInfoPlaintext(s store.Store) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vendorTypes := strings.Split(req.URL.Query()["types"][0], ",")
		responses := make([]byte, len(vendorTypes))
		for i, t := range vendorTypes {
			_, infoStatus, err := shared.GetEventVendorTypeInfo(s, t)
			if err != nil {
				log.Println(err)
				w.WriteHeader(200)
				return
			}
			switch infoStatus {
			case shared.StatusOk:
				responses[i] = '0'
			case shared.StatusRecovering:
				responses[i] = '1'
			case shared.StatusAbnormal:
				responses[i] = '2'
			}
		}
		w.Write(responses)
	}
}
