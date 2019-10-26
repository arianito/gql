package gql

import (
	"database/sql"
	"errors"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"log"
	"time"
)

const GTChannel = "global-transactions"

type Transaction struct {
	tx        *sql.Tx
	expiresAt time.Time
}

type TransactionManager struct {
	srv          micro.Service
	db           *sql.DB
	transactions map[string]*Transaction
	broker       broker.Broker
	mutex        *Mutex
}

func NewTransactionManager(srv micro.Service, db *sql.DB) (*TransactionManager, error) {
	trs := make(map[string]*Transaction)
	b := srv.Options().Broker

	gt := &TransactionManager{
		transactions: trs, db: db, broker: b, srv: srv,
		mutex: &Mutex{},
	}

	go func() {
		for {
			time.Sleep(time.Second * 10)
			if unlock, err := WithLock(gt.mutex); err == nil {

				for key, value := range gt.transactions {
					if value.expiresAt.Before(time.Now()) {
						err := gt.rollbackLocal(key)
						if err != nil {
							log.Println("rollback failed -> ", key, err.Error())
						}
					}
				}
				unlock()
			} else {
				log.Println("failed to acquire lock")
			}
		}
	}()

	_, err := b.Subscribe(GTChannel, func(e broker.Event) error {
		head := e.Message().Header
		id := string(e.Message().Body)
		if head["action"] == "commit" {
			return gt.commitLocal(id)
		} else if head["action"] == "rollback" {
			return gt.rollbackLocal(id)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return gt, nil

}

func (t *TransactionManager) Transact(id string, txFunc func(*sql.Tx) error) (err error) {
	if unlock, err := WithLock(t.mutex); err == nil {
		defer unlock()
		tx, ok := t.transactions[id]
		if !ok {
			txx, err := t.db.Begin()
			if err != nil {
				return err
			}
			tx = &Transaction{
				tx:        txx,
				expiresAt: time.Now().Add(time.Minute * 5),
			}
			t.transactions[id] = tx
		}
		defer func() {
			if p := recover(); p != nil {
				t.broker.Publish(GTChannel, &broker.Message{
					Header: map[string]string{
						"action": "rollback",
					},
					Body: []byte(id),
				})
				panic(p)
			} else if err != nil {
				t.broker.Publish(GTChannel, &broker.Message{
					Header: map[string]string{
						"action": "rollback",
					},
					Body: []byte(id),
				})
			} else {
				err = t.broker.Publish(GTChannel, &broker.Message{
					Header: map[string]string{
						"action": "commit",
					},
					Body: []byte(id),
				})
			}
		}()
		err = txFunc(tx.tx)
		return err
	} else {
		return err
	}
}


func (t *TransactionManager) rollbackLocal(id string) error {
	if tx, ok := t.transactions[id]; ok {
		err := tx.tx.Rollback()
		if err != nil {
			return err
		}
		delete(t.transactions, id)
		return nil
	}
	return errors.New("transaction " + id + " not found")
}

func (t *TransactionManager) commitLocal(id string) error {
	if tx, ok := t.transactions[id]; ok {
		err := tx.tx.Commit()
		if err != nil {
			return err
		}
		delete(t.transactions, id)
		return nil
	}
	return errors.New("transaction " + id + " not found")
}
