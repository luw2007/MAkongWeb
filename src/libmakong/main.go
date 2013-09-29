// makongs project main.go
package libmakong

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type ResponseNiuxba struct {
	Error int
	Data  interface{}
}

type Friend struct {
	CardId string `json:"cardId"`
	Id     string `json:"id"`
	Name   string `json:"name"`
	Level  string `json"level"`
	Role   string `json"role"`
	Max    string `json"max"`
	Now    string `json"now"`
	Active string `json"active"`
}

type Friends struct {
	Users []Friend `json:"friends"`
}

type User struct {
	UserId   string `json:"userId"`
	GroupId  int    `json:"groupId"`
	ServerNo string `json:"serverNo"`
	Phone    string `json:"phone"`
}

type Regist struct {
	Phone    string `json:"phone"`
	NickName string `json:"nickName"`
	GroupId  int    `json:"groupId"`
}

type POST_DATA struct {
	Data   User `json:"data"`
	Result int  `json:"result"`
}

func get_niuxba(_url string) []byte {

	fmt.Println("get_niuxba:")
	resp, err := http.Get(_url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//
	}
	return body
}

func post_niuxba(_url string, data string) []byte {
	resp, err := http.PostForm(_url, url.Values{"data": {data}})
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//
	}
	return body
}

func Get_ranking(post_data User) []byte {
	_url := "http: //www.niuxba.com/ma/backend/cgi-bin/getZsjRankInfo.php?phone=" + post_data.Phone + "&groupId=" + string(post_data.GroupId) + "&serverNo=" + post_data.ServerNo + "&userId=" + post_data.UserId

	return get_niuxba(_url)
}
func Get_football(post_data User, keyword string) []byte {
	_url := "http: //www.niuxba.com/ma/backend/cgi-bin/getFootball2.php?phone=" + post_data.Phone + "&groupId=" + string(post_data.GroupId) + "&serverNo=" + post_data.ServerNo + "&userId=" + post_data.UserId
	resp, err := http.PostForm(_url, url.Values{"keyword": {keyword}})
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//
	}
	return body

}

func Get_friends_info(data string) []byte {
	body := post_niuxba("http://www.niuxba.com/ma/backend/cgi-bin/getFriendInfo.php", data)
	return body
}

func Get_user_info(data string) []byte {
	body := post_niuxba("http://www.niuxba.com/ma/backend/cgi-bin/getUserInfo.php", data)
	return body
}

func regeist(registInfo string) (POST_DATA, error) {
	body := post_niuxba("http://www.niuxba.com/ma/backend/cgi-bin/bindUser3.php", registInfo)
	var data POST_DATA
	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(body, data)
	return data, err
}

func Get_post_data(userFile string, phone string, nickName string, groupId int) (User, string, error) {
	r := Regist{Phone: phone, NickName: nickName, GroupId: groupId}
	registInfo, err := json.Marshal(r)
	if err != nil {
		fmt.Println("json err:", err)
	}
	result, err := regeist(string(registInfo))

	basePostData, err := json.Marshal(result.Data)
	if err != nil {
		fmt.Println("json err:", err)
	}
	Write_config(userFile, basePostData)
	fmt.Println("Write_config", string(basePostData))
	return result.Data, string(basePostData), err
}

func Write_config(userFile string, data []byte) {
	fout, err := os.Create(userFile)
	defer fout.Close()
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	fout.Write(data)
}

var (
	phone    string
	nickName string
	serverNo string
	groupId  int
	userId   string
)

func Read_config(userFile string, post_data *User) (string, error) {
	_, err := os.Stat(userFile)
	if err != nil {
		//no such file
		fmt.Println("read file error", err)
		return "", err
	}
	config, err := ioutil.ReadFile(userFile)
	if err != nil {
		fmt.Println("read file error", err)
		return "", err
	}

	fmt.Println("Write_config", string(config))
	return string(config), json.Unmarshal(config, post_data)
}

func command() {

	userFile := phone + "_user.json"
	var post_data User
	basePostData, err := Read_config(userFile, &post_data)
	if err != nil {
		post_data, basePostData, err = Get_post_data(userFile, phone, nickName, groupId)
	}

	if post_data.UserId == "" {
		panic("error, check input")
		//fmt.Println("error, check input")
		//flag.Usage()
		//os.Exit(2)
	}

	fmt.Printf("%d区[%s] %s serverNo:%s userId:%s\n",
		post_data.GroupId, nickName,
		post_data.Phone, post_data.ServerNo, post_data.UserId)
	fmt.Printf("排位赛：http: //www.niuxba.com/ma/backend/cgi-bin/getZsjRankInfo.php?phone=%s&groupId=%d&serverNo=%s&userId=%s\n",
		post_data.Phone, post_data.GroupId, post_data.ServerNo, post_data.UserId)

	user_info := Get_user_info(basePostData)
	js, err := simplejson.NewJson(user_info)
	if err != nil {
		fmt.Println("json err:", err)
	}
	info, err := js.Get("data").Map()
	if err != nil {
		fmt.Println("json err:", err)
	}
	for k, v := range info {
		switch vv := v.(type) {
		case int:
			fmt.Printf("%s: %d\t", k, vv)
		case string:
			fmt.Printf("%s: %s\t", k, vv)
		case float64:
			fmt.Printf("%s: %d\t", k, int(vv))
		default:
			fmt.Println(k, "is of a type I don't know how to handle", vv)
		}
	}
	fmt.Println()

	data := Get_friends_info(basePostData)
	js, err = simplejson.NewJson(data)
	if err != nil {
		fmt.Println("json err:", err)
	}
	friends, err := js.Get("data").Get("friends").Array()
	if err != nil {
		fmt.Println("json err:", err)
	}
	for _, friend := range friends {
		m := friend.(map[string]interface{})
		v := m["active"]
		str_str, _ := strconv.Atoi(v.(string))
		if int(str_str) < 150 {
			for k, v := range friend.(map[string]interface{}) {
				fmt.Printf("%s: %s\t", k, v.(string))
			}
			fmt.Printf("\n")
		}

	}

	var waint string
	fmt.Println("任意键推出")
	fmt.Scanln(&waint)
}
