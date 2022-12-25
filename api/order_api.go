package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/2gis-demo-app/log"
	"github.com/2gis-demo-app/oms"
)

var AvailableRooms = map[string]struct{}{"econom": {}, "standart": {}, "lux": {}} //Some Stock API

type OrderApi struct {
	Mux           *http.ServeMux
	OrderProvider *oms.OrderProvider
	logger        log.Logger
}

func NewApi(logger log.Logger, provider *oms.OrderProvider) *OrderApi {
	api := &OrderApi{
		Mux:           http.NewServeMux(),
		OrderProvider: provider,
		logger:        logger,
	}
	api.Mux.HandleFunc("/order", api.MakeOrder)
	api.Mux.HandleFunc("/orders", api.GetOrders)
	return api
}

func (s *OrderApi) GetOrders(w http.ResponseWriter, r *http.Request) {
	//s.logger.LogErrorf("API GetOrders: %s", "start")
	userEmail := r.URL.Query().Get("email")
	if userEmail == "" {
		s.logger.LogErrorf("error in getOrders method: %s", "empty email")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//s.logger.LogErrorf("API GetOrders: %s", "email is OK")
	res, err := s.OrderProvider.GetOrders(userEmail)
	if err != nil {
		s.logger.LogInfo("error in getOrders method: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(res)
	if err != nil {
		s.logger.LogErrorf("error in getOrders method: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		s.logger.LogErrorf("error in getOrders method: %s", err.Error())
		return
	}

	s.logger.LogInfo("Method getOrders was successfully done")
}

func (s *OrderApi) MakeOrder(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("email")
	if userEmail == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	room := r.URL.Query().Get("room")
	if room == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if _, isOK := AvailableRooms[room]; !isOK {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	from := r.URL.Query().Get("from")
	if from == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fromTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	to := r.URL.Query().Get("to")
	if to == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	toTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = s.OrderProvider.AddOrder(userEmail, room, fromTime, toTime)
	if err != nil {
		s.logger.LogErrorf("error in makeOrders method: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	s.logger.LogInfo("Method makeOrder was successfully done")
}
