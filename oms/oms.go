package oms

import (
	"applicationDesignTest/app"
	"fmt"
	"sync"
	"time"
)

type OrderProvider struct {
	mux    sync.RWMutex
	orders map[string][]Order
}

type Order struct {
	Room      string
	UserEmail string
	From      time.Time
	To        time.Time
}

func NewOMS(app *app.Application) *OrderProvider {
	return &OrderProvider{
		mux:    sync.RWMutex{},
		orders: make(map[string][]Order),
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
		return nil, fmt.Errorf("email %s not found", email)
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
	if currTime.Before(from) {
		return fmt.Errorf("invalid from time")
	}
	if currTime.Before(from) {
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
	if oms.IsAvailable(newOrder) {
		oms.orders[email] = append(oms.orders[email], newOrder)
	} else {
		return fmt.Errorf("room is busy")
	}
	return nil
}
func (oms *OrderProvider) IsAvailable(order Order) bool {
	//TODO: checks
	return true
}
