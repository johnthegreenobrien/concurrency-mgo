package controllers

import (
	//"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	//"sync"
	"time"
	"userv/commons/cache"
	"userv/commons/database"
	"userv/commons/log"
	"userv/modules/dailyDelivery/dao"
	"userv/modules/dailyDelivery/models"
)

/**
 *
 */
type deliveryController struct {
	mSession  *database.MongoSession
	rClient   *cache.RedisClient
	dailyColl string
}

/**
 *
 */
func DeliveryController() *deliveryController {
	return &deliveryController{mSession: database.ConnMongo(), rClient: cache.ConnRedis()}
}

/**
 * Compare length of the daily delivery and addresses collection
 */
func (uc *deliveryController) compareDeliveryAddress() bool {
	resultDaily, errd := uc.mSession.CountColl(uc.dailyColl)
	if errd != nil {

	}
	resultAddress, erra := uc.mSession.CountColl("address")
	if erra != nil {

	}
	fmt.Println(resultDaily)
	fmt.Println(resultAddress)
	if resultAddress > resultDaily {
		fmt.Println("Address Collection is bigger than daily delivery Collection")
		return false
	}
	return true
}

/**
 * Checks if the daily delivery collection exists
 */
func (uc *deliveryController) checkDeliveryCollection() bool {
	resultDaily, _ := uc.mSession.CountColl(uc.dailyColl)
	if resultDaily > 0 {
		//fmt.Println("Bigger than 0")
		return true
	}
	//fmt.Println("Collection", collName, "is empty")
	return false
}

/**
 * Creates a daily delivery collection
 */
func (uc *deliveryController) createDailyDeliveryCollection() {
	fmt.Println("We're going to ensure an index and the existence of the", uc.dailyColl, "Collection")
	deliveryDao := dao.NewDeliveryDao(uc.dailyColl)
	err := deliveryDao.CreateDailyCollection(uc.mSession, uc.dailyColl)
	if err != nil {
		fmt.Println("An error has occurred when trying to ensure", uc.dailyColl, "Collection")
	} else {
		fmt.Println("Collection", uc.dailyColl, " was successfully ensured")
	}
}

/**
 *
 */
func (uc *deliveryController) importDeliveryAddress() *models.Address {
	addressDao := dao.NewAddressDao()
	skip, _ := uc.mSession.CountColl(uc.dailyColl)
	fmt.Println("We're going to import an address, skipping", skip, "docs")
	address, _ := addressDao.GetAddress(uc.mSession, skip)
	fmt.Println("We're going to import", address, "to the daily delivery Collection.")
	return address
}

/**
 * Add the new address and
 */
func (uc *deliveryController) updateDeliveryAddresses() {
	fmt.Println("We're going to update the daily delivery Collection docs")
	address := uc.importDeliveryAddress()

	deliveryDao := dao.NewDeliveryDao(uc.dailyColl)
	delivery := models.Delivery{Address: address.Address, ZipCode: address.ZipCode, SussDlry: 0}
	deliveryDao.InsertDelivery(uc.mSession, &delivery)
}

/**
 *
 */
func (uc *deliveryController) getDailyDelivery() *models.Delivery {
	deliveryDao := dao.NewDeliveryDao(uc.dailyColl)
	delivery, _ := deliveryDao.GetDelivery(uc.mSession)
	delivery, _ = deliveryDao.IncrementField(uc.mSession, "sussDlry", delivery)
	return delivery
}

/**
 * @method GetDelivery
 */
func (uc *deliveryController) GetDelivery(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uc.dailyColl = time.Now().Format("2006_01_02")

	logging, _ := log.NewSysLog()
	logging.Debug()
	fmt.Println(logging.GetTraceMsg())

	if uc.compareDeliveryAddress() {
		fmt.Println(uc.getDailyDelivery())
	} else if uc.checkDeliveryCollection() {
		uc.updateDeliveryAddresses()
		uc.getDailyDelivery()
	} else {
		uc.createDailyDeliveryCollection()
		uc.updateDeliveryAddresses()
		uc.getDailyDelivery()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	//dlryJson, _ := json.Marshal(delivery)
	//fmt.Fprintf(w, "%s", dlryJson)
}
