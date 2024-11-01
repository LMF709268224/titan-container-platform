package order

import (
	"time"

	"titan-container-platform/core"
	"titan-container-platform/core/dao"
	"titan-container-platform/kubesphere"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("order")

const (
	timeInterval = 2 * time.Minute
)

// Init initializes the order manager.
func Init() {
	go startTimer()
}

func startTimer() {
	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()

	for {
		<-ticker.C

		checkOrderPaid()
		createSpaceFromOrders()
	}
}

func checkOrderPaid() {
	list, err := dao.LoadOrdersByStatus(core.OrderStatusCreated)
	if err != nil {
		log.Errorf("LoadOrdersByStatus err:%s", err.Error())
		return
	}

	for _, order := range list {
		// if time < xxx set status = OrderStatusTimeout
		// TODO check order from TOKEN
		// if true
		err := dao.UpdateOrderStatus(order.ID, core.OrderStatusPaid)
		if err != nil {
			log.Errorf("UpdateOrderStatus %s err:%s", order.ID, err.Error())
		}
	}
}

func createSpaceFromOrders() {
	list, err := dao.LoadOrdersByStatus(core.OrderStatusPaid)
	if err != nil {
		log.Errorf("LoadOrdersByStatus err:%s", err.Error())
		return
	}

	for _, order := range list {
		status := core.OrderStatusDone

		err = kubesphere.CreateSpaceAndResourceQuotas(order.ID, order.Account, order.CPUCores, order.RAMSize, order.StorageSize)
		if err != nil {
			log.Errorf("LoadOrdersByStatus err:%s", err.Error())

			status = core.OrderStatusExpired
		}

		err := dao.UpdateOrderStatus(order.ID, status)
		if err != nil {
			log.Errorf("UpdateOrderStatus %s err:%s", order.ID, err.Error())
		}
	}
}
