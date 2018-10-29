package phone

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/jonnywang/go-kits/redis"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type TPhoneDetail struct {
	pType     string
	pRegion   string
	pCity     string
	pDistrict string
	pZipCode  string
}

type phonesRedisHandler struct {
	redis.RedisHandler
	phones  map[string]*TPhoneDetail
	num     int
	version string
	sync.Mutex
}

func (obj *phonesRedisHandler) Init() error {
	obj.Initiation(func() {
		go obj.loadPhones()
	})

	return nil
}

func (obj *phonesRedisHandler) loadPhones() {
	obj.Lock()
	defer obj.Unlock()

	obj.phones = make(map[string]*TPhoneDetail, 0)
	obj.num = 0

	content, err := ioutil.ReadFile(GConfig.PhoneDict)
	if err != nil {
		Logger.Printf("read phone.dat failed %v", err)
		return
	}

	Logger.Printf("start load phone.dat")

	totalLen := len(content)
	obj.version = string(content[:4])

	firstIndexPosition := binary.LittleEndian.Uint32(content[4:8])

	var phoneItem []byte
	i := int(firstIndexPosition)
	for ; i < totalLen; i += 9 {
		phoneItem = content[i : i+9]

		phoneNumber := fmt.Sprintf("%d", binary.LittleEndian.Uint32(phoneItem[0:4]))
		obj.phones[phoneNumber] = &TPhoneDetail{
			pType:     fmt.Sprintf("%d", phoneItem[8]), //1移动 2联通 3电信 4电信虚拟运营商 5联通虚拟运营商 6移动虚拟运营商
			pRegion:   "",
			pCity:     "",
			pDistrict: "",
			pZipCode:  "",
		}

		phoneDetailIndex := int(binary.LittleEndian.Uint32(phoneItem[4:8]))
		stepLen := 500

		for {
			phoneDetail := content[phoneDetailIndex : phoneDetailIndex+stepLen]
			phoneDeatilEndIndex := bytes.Index(phoneDetail, []byte("\000"))
			if -1 == phoneDeatilEndIndex {
				stepLen += 100
				continue
			}

			phoneDetailList := bytes.Split(phoneDetail[:phoneDeatilEndIndex], []byte("|"))
			obj.phones[phoneNumber].pRegion = string(phoneDetailList[0])
			obj.phones[phoneNumber].pCity = string(phoneDetailList[1])
			obj.phones[phoneNumber].pZipCode = string(phoneDetailList[2])
			obj.phones[phoneNumber].pDistrict = string(phoneDetailList[3])
			break
		}

		obj.num++
	}

	Logger.Printf("load phone.dat finished, found total %d phones", obj.num)
}

func (obj *phonesRedisHandler) Shutdown() {
	Logger.Print("mailer server will shutdown")
}

func (obj *phonesRedisHandler) Version() (string, error) {
	return VERSION, nil
}

func (obj *phonesRedisHandler) Ping(message string) (string, error) {
	if len(message) > 0 {
		return message, nil
	}

	return "PONG", nil
}

func (obj *phonesRedisHandler) Total() (int, error) {
	obj.Lock()
	defer obj.Unlock()

	return obj.num, nil
}

func (obj *phonesRedisHandler) Reload() (error) {
	go obj.loadPhones()

	return nil
}

func (obj *phonesRedisHandler) Hgetall(phoneNumber string) (map[string]interface{}, error) {
	if len(phoneNumber) == 0 || len(phoneNumber) < 7 {
		return nil, errors.New("error params")
	}

	obj.Lock()
	defer obj.Unlock()

	phoneDetail, ok := obj.phones[phoneNumber[:7]]
	if !ok {
		return nil, nil
	}

	Logger.Printf("found %s detail %v", phoneNumber, phoneDetail)

	out := make(map[string]interface{}, 0)
	out["type"] = phoneDetail.pType
	out["region"] = phoneDetail.pRegion
	out["city"] = phoneDetail.pCity
	out["district"] = phoneDetail.pDistrict
	out["zip"] = phoneDetail.pZipCode

	return out, nil
}

func Run() {
	phoneHandler := &phonesRedisHandler{}

	err := phoneHandler.Init()
	if err != nil {
		Logger.Print(err)
		return
	}

	phoneRedisServer, err := redis.NewServer(GConfig.ListenServer, phoneHandler)
	if err != nil {
		Logger.Print(err)
		return
	}

	serverStop := make(chan bool)
	stopSignal := make(chan os.Signal)
	signal.Notify(stopSignal, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-stopSignal
		Logger.Print("catch exit signal")
		phoneRedisServer.Stop(10)
		serverStop <- true
	}()

	err = phoneRedisServer.Start()
	if err != nil {
		Logger.Print(err)
		stopSignal <- syscall.SIGTERM
	}

	<-serverStop

	close(serverStop)
	close(stopSignal)

	Logger.Print("all server shutdown")
}
