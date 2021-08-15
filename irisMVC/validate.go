package main

import (
	"errors"
	"fmt"
	"go-flash-sale/irisMVC/common"
	"go-flash-sale/irisMVC/encrypt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

var hostArray = []string{"127.0.0.1", "127.0.0.1"}

var localHost = "127.0.0.1"

var port = "8081"

var hashConsistent *common.Consistent

// AccessControl to store control info
type AccessControl struct {
	//store meta data
	sourcesArray map[int]interface{}
	sync.RWMutex
}

var accessControl = &AccessControl{sourcesArray:make(map[int]interface{})}

// GetNewRecord get data
func (m *AccessControl) GetNewRecord(uid int) interface{} {
	m.RWMutex.RLock()
	defer m.RWMutex.RUnlock()
	data:=m.sourcesArray[uid]
	return data
}

// SetNewRecord set data
func (m *AccessControl) SetNewRecord(uid int) {
	m.RWMutex.Lock()
	m.sourcesArray[uid]="hello world"
	m.RWMutex.Unlock()
}

func (m *AccessControl) GetDistributedRight(req *http.Request) bool {
	//get user UID
	uid ,err := req.Cookie("uid")
	if err !=nil {
		return false
	}

	//by consistent hash，get the server by user ID
	hostRequest,err:=hashConsistent.Get(uid.Value)
	if err !=nil {
		return false
	}

	//check if localhost
	if hostRequest == localHost {
		//local process
		return m.GetDataFromMap(uid.Value)
	} else {
		//as proxy
		return GetDataFromOtherMap(hostRequest,req)
	}

}

//get result from local machine
func (m *AccessControl) GetDataFromMap(uid string) (isOk bool) {
	uidInt,err := strconv.Atoi(uid)
	if err !=nil {
		return false
	}
	data:=m.GetNewRecord(uidInt)

	if data !=nil {
		return true
	}
	return
}

//get result from other machine
func GetDataFromOtherMap(host string,request *http.Request) bool  {
	uidPre,err := request.Cookie("uid")
	if err !=nil {
		return false
	}
	//get sign
	uidSign,err:=request.Cookie("sign")
	if err !=nil {
		return  false
	}

	client :=&http.Client{}
	req,err:= http.NewRequest("GET","http://"+host+":"+port+"/check",nil)
	if err !=nil {
		return false
	}

	//check if unused cookie
	cookieUid :=& http.Cookie{Name:"uid",Value:uidPre.Value,Path:"/"}
	cookieSign :=& http.Cookie{Name:"sign",Value:uidSign.Value,Path:"/"}
	req.AddCookie(cookieUid)
	req.AddCookie(cookieSign)

	response,err :=client.Do(req)
	if err !=nil {
		return false
	}
	body,err:=ioutil.ReadAll(response.Body)
	if err !=nil {
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

func Check(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Validating！")
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

	filter := common.NewFilter()
	filter.RegisterFilterUri("/check", Auth)
	http.HandleFunc("/check", filter.Handle(Check))
	http.ListenAndServe(":8083", nil)
}
