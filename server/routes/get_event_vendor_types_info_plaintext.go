package routes

import (
	"log"
	statusEngine "main/status"
	"net/http"
	"strings"
)

func GetEventVendorTypesInfoPlaintext(se statusEngine.StatusEngine) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vendorTypes := strings.Split(req.URL.Query()["types"][0], ",")
		responses := make([]byte, len(vendorTypes))
		for i, t := range vendorTypes {
			s, err := se.ProcessEventsForVendorType(t)
			if err != nil {
				log.Println(err)
				w.WriteHeader(200)
				return
			}
			switch s.Status {
			case statusEngine.StatusOk:
				responses[i] = '0'
			case statusEngine.StatusRecovering:
				responses[i] = '1'
			case statusEngine.StatusAbnormal:
				responses[i] = '2'
			}
		}
		w.Write(responses)
	}
}
