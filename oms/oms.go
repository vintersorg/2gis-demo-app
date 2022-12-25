package oms

import (
	"fmt"
	"sync"
	"time"

	"github.com/2gis-demo-app/log"
)

type OrderProvider struct {
	mux    sync.RWMutex
	orders map[string][]Order
	logger log.Logger
}

type Order struct {
	Room      string
	UserEmail string
	From      time.Time
	To        time.Time
}

type BusyInterval struct {
	From time.Time
	To   time.Time
}

func NewOMS(logger log.Logger) *OrderProvider {
	return &OrderProvider{
		mux:    sync.RWMutex{},
		orders: make(map[string][]Order),
		logger: logger,
	}
}
func (oms *OrderProvider) GetOrders(email string) ([]Order, error) {
	if email == "" {
		return nil, fmt.Errorf("email is empty")
	}
	oms.mux.RLock()
	defer oms.mux.RUnlock()
	res, isOK := oms.orders[email]
	if !isOK {
		//TODO: some logic?
		return []Order{}, nil
	}
	return res, nil
}
func (oms *OrderProvider) AddOrder(email, room string, from, to time.Time) error {
	if email == "" {
		return fmt.Errorf("email is empty")
	}
	if room == "" {
		return fmt.Errorf("room is empty")
	}
	currTime := time.Now()
	//oms.logger.LogErrorf("OMS AddOrder: currTime is %s", currTime)
	//oms.logger.LogErrorf("OMS AddOrder: from is %s", from)
	//oms.logger.LogErrorf("OMS AddOrder: to is %s", to)
	if from.Before(currTime) {
		return fmt.Errorf("invalid from time")
	}
	if to.Before(currTime) {
		return fmt.Errorf("invalid to time")
	}
	if to.Before(from) {
		return fmt.Errorf("invalid time range")
	}
	newOrder := Order{
		Room:      room,
		UserEmail: email,
		From:      from,
		To:        to,
	}
	oms.mux.Lock()
	defer oms.mux.Unlock()
	if oms.isAvailable(newOrder) {
		oms.orders[email] = append(oms.orders[email], newOrder)
	} else {
		return fmt.Errorf("room is busy")
	}
	return nil
}
func (oms *OrderProvider) isAvailable(order Order) bool {
	//TODO: checks
	//TODO: need optimization. Some Stock API
	for _, list := range oms.orders {
		for _, item := range list {
			//no need to check
			if item.Room != order.Room {
				continue
			}
			//new interval inside
			if order.From.After(item.From) && order.To.Before(item.To) {
				return false
			}
			//new interval cross right
			if order.From.After(item.From) && order.From.Before(item.To) {
				return false
			}
			//new interval cross left
			if order.To.After(item.From) && order.To.Before(item.To) {
				return false
			}
			//new interval bigger
			if order.From.Before(item.From) && order.To.After(item.To) {
				return false
			}
			//new interval dots equals
			if order.From.Equal(item.From) || order.To.Equal(item.To) || order.From.Equal(item.To) || order.To.Equal(item.From) {
				return false
			}
		}
	}
	return true
}
