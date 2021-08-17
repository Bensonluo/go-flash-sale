package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-flash-sale/irisMVC/common"
	"go-flash-sale/irisMVC/datamodels"
	"go-flash-sale/irisMVC/encrypt"
	"go-flash-sale/irisMVC/rabbitmq"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var hostArray = []string{"127.0.0.1", "127.0.0.1"}
var localHost = ""
var port = "8083"

//数量控制接口服务内网IP，或者getOne的SLB内网IP
var GetOneIp = "127.0.0.1"
var GetOnePort = "8084"
var hashConsistent *common.Consistent
var rabbitmqValidate *rabbitmq.RabbitMQ

// AccessControl to store control info
type AccessControl struct {
	sourcesArray map[int]time.Time
	sync.RWMutex
}

//server interval
var interval = 20

var accessControl = &AccessControl{sourcesArray: make(map[int]time.Time)}

func (m *AccessControl) GetNewRecord(uid int) time.Time {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	return m.sourcesArray[uid]
}

func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.sourcesArray[uid] = time.Now()
	m.RWMutex.Unlock()
}

//blacklist
type BlackList struct {
	listArray map[int]bool
	sync.RWMutex
}

var blackList = &BlackList{listArray: make(map[int]bool)}

func (m *BlackList) GetBlackListByID(uid int) bool {
	m.RLock()
	defer m.RUnlock()
	return m.listArray[uid]
}

func (m *BlackList) SetBlackListByID(uid int) bool {
	m.Lock()
	defer m.Unlock()
	m.listArray[uid] = true
	return true
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	//get user UID
	uid, err := req.Cookie("uid")
	if err != nil {
		return false
	}

	//by consistent hash，get the server by user ID
	hostRequest, err := hashConsistent.Get(uid.Value)
	if err != nil {
		return false
	}

	//check if localhost
	if hostRequest == localHost {
		//local process
		return m.GetDataFromMap(uid.Value)
	} else {
		//as proxy
		return GetDataFromOtherMap(hostRequest, req)
	}

}

//get result from local machine
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}
	//add to blacklist
	if blackList.GetBlackListByID(uidInt) {
		return false
	}
	//get record
	dataRecord := m.GetNewRecord(uidInt)
	if !dataRecord.IsZero() {
		if dataRecord.Add(time.Duration(interval) * time.Second).After(time.Now()) {
			return false
		}
	}
	m.SetNewRecord(uidInt)
	return true
}

//get result from other machine
func GetDataFromOtherMap(host string, request *http.Request) bool {
	hostUrl := "http://" + host + ":" + port + "/check"
	response, body, err := GetCurl(hostUrl, request)
	if err != nil {
		return false
	}
	//check status
	if response.StatusCode == 200 {
		if string(body) == "true" {
			return true
		} else {
			return false
		}
	}
	return false
}

func Auth(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("processing！")
	//check cookie
	err := CheckUserInfo(r)
	if err != nil {
		return err
	}
	return nil
}

func GetCurl(hostUrl string, request *http.Request) (response *http.Response, body []byte, err error) {
	uidPre, err := request.Cookie("uid")
	if err != nil {
		return
	}
	//get sign
	uidSign, err := request.Cookie("sign")
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+hostUrl+":"+port+"/check", nil)
	if err != nil {
		return
	}

	//check if unused cookie
	cookieUid := &http.Cookie{Name: "uid", Value: uidPre.Value, Path: "/"}
	cookieSign := &http.Cookie{Name: "sign", Value: uidSign.Value, Path: "/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	response, err = client.Do(req)
	defer response.Body.Close()

	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}

func Check(w http.ResponseWriter, r *http.Request) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil && len(queryForm["productID"]) <= 0 && len(queryForm["productID"][0]) <= 0 {
		w.Write([]byte("false"))
		return
	}
	productString := queryForm["productID"][0]
	fmt.Println(productString)
	//get user cookie
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("false"))
		return
	}
	//分布式鉴权
	right := accessControl.GetDistributedRight(r)
	if right == false {
		w.Write([]byte("false"))
	}
	//获取数量权限，防止超卖
	hostUrl := "http ://" + GetOneIp + ":" + GetOnePort + "/getOne"
	responseValidate, validateBody, err := GetCurl(hostUrl, r)
	if err != nil {
		w.Write([]byte("false"))
	}
	//
	if responseValidate.StatusCode == 200 {
		if string(validateBody) == "true" {
			//place order
			productID, err := strconv.ParseInt(productString, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
			}
			userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
			if err != nil {
				w.Write([]byte("false"))
				return
			}

			message := datamodels.Message{userID, productID}
			byteMessage, err := json.Marshal(message)
			if err != nil {
				w.Write([]byte("false"))
				return
			}
			err = rabbitmqValidate.PublishSimple(string(byteMessage))
			w.Write([]byte("true"))
			return
		}
	}
	w.Write([]byte("false"))
	return
}

func CheckRight(w http.ResponseWriter, r *http.Request) {
	right := accessControl.GetDistributedRight(r)
	if !right {
		w.Write([]byte("false"))
		w.Write([]byte("false"))
		return

	}
	w.Write([]byte("true"))
	return
}

func CheckUserInfo(r *http.Request) error {
	//get Uid，cookie
	uidCookie, err := r.Cookie("uid")
	if err != nil {
		return errors.New("Fail to get UID Cookie！")
	}

	signCookie, err := r.Cookie("sign")
	if err != nil {
		return errors.New("Fail to get Cookie sign！")
	}

	signByte, err := encrypt.DePwdCode(signCookie.Value)
	if err != nil {
		return errors.New("Sign error")
	}

	fmt.Println("compare")
	fmt.Println("user ID：" + uidCookie.Value)
	fmt.Println("decrypted ID：" + string(signByte))
	if checkInfo(uidCookie.Value, string(signByte)) {
		return nil
	}
	//return errors.New("Invalid！")
	return nil
}

func checkInfo(checkStr string, signStr string) bool {
	if checkStr == signStr {
		return true
	}
	return false
}

func main() {
	hashConsistent = common.NewConsistent()
	//add to hash circle
	for _, v := range hostArray {
		hashConsistent.Add(v)
	}
	localIp, err := common.GetIp()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIp
	fmt.Println(localHost)

	rabbitmqValidate = rabbitmq.NewRabbitMQSimple("goFlashSale")
	defer rabbitmqValidate.Destroy()

	//static files
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./irisMVC/frontend/web/htmlProductShow"))))
	//resource
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./irisMVC/frontend/web/public"))))

	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth)
	filter.RegisterFilterUri("/checkRight", Auth)

	http.HandleFunc("/check", filter.Handle(Check))
	http.HandleFunc("/checkRight", filter.Handle(CheckRight))

	http.ListenAndServe(":8083", nil)
}
