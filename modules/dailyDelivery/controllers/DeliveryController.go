package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
	"userv/commons/cache"
	"userv/commons/database"
	"userv/modules/dailyDelivery/dao"
	//"userv/modules/dailyDelivery/models"
)

/**
 *
 */
type deliveryController struct {
	mSession *database.MongoSession
}

/**
 *
 */
func DeliveryController(mongoSession *database.MongoSession) *deliveryController {
	return &deliveryController{mongoSession}
}

/**
 *
 */
func (uc *deliveryController) GetDelivery(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	deliveryDao := dao.NewDeliveryDao()

	var waitGroup sync.WaitGroup
	delivery, _ := deliveryDao.GetDelivery(&waitGroup, uc.mSession)
	waitGroup.Wait()
	delivery, _ = deliveryDao.IncrementField(&waitGroup, uc.mSession, "sussDlry", delivery)
	waitGroup.Wait()
	cache.Set("thiIsAKey", "keysValue", &waitGroup)
	fmt.Println(cache.Get("thiIsAKey", &waitGroup))
	waitGroup.Wait()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	dlryJson, _ := json.Marshal(delivery)
	fmt.Fprintf(w, "%s", dlryJson)
}