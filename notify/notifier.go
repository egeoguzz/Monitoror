package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"

	"errors"
	"net/http"
	"time"
)

var problem_count int
var success_msg bool
var Notif_Url string
var Notif_Emoji bool
var notif_channel bool
var Ment_Lst string
var err_code string
var tiles []Tile

type SlackRequestBody struct {
	Pretext string  `json:"pretext"`
	Color   string  `json:"color"`
	Fields  []Field `json:"fields"`
	Text    string  `json:"text"`
}
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Columns struct {
	Column_Count int                      `json:"columns"`
	Tiles        []map[string]interface{} `json:"tiles"`
	Mention_lst  string                   `json:"mention_list"`
	Label        string                   `json:"label"`
}
type Columns_Tile struct {
}
type Tile struct {
	Type                  string
	Label                 string
	Name                  string
	ID                    int
	current_problem_count int
	mnt_lst               []string
	error_cd              string
	Group                 bool
	GrpLabel              string
}

func isErrorOccured(ID int, isErr bool, ercode string) {

	if isErr {
		increaseErrCount(ID)
		if getErrCount(ID) == problem_count {
			attachErCode(ID, ercode)
			getMsg(ID, true)
		}
	} else {
		if getErrCount(ID) >= problem_count {
			if success_msg {
				getMsg(ID, false)
			}

		}
		resetErrCount(ID)
	}

}
func getMsg(ID int, flag bool) {

	var error_caution string
	var to_be_notified []string
	var special_notifications string
	var general_notifications string
	success_emoji := [5]string{":party_blob:", ":meow_fat:", ":cool-doge:", ":thumbsup_all:", ":fast_parrot:"}
	failure_emoji := [5]string{":typingcat:", ":mild-panic-intensifies:", ":confused_dog:", ":portalparrot:", ":ahhhhhhhhh:"}
	min := 0
	max := 5
	rnd := rand.Intn(max-min) + min
	emoji := ""
	if flag {
		error_caution = `"text": ":red_circle:  *Something bad is happening.*  :red_circle:"  ,`
		if Notif_Emoji {
			emoji = failure_emoji[rnd]
		}
	} else {
		error_caution = `"text": ":large_green_circle:  *Problem seems to be solved.*  :large_green_circle:" ,`
		if Notif_Emoji {
			emoji = success_emoji[rnd]
		}

	}
	for i := 0; i < len(tiles); i++ {
		if tiles[i].ID == ID {
			if tiles[i].mnt_lst != nil {
				for j := 0; j < len(tiles[i].mnt_lst); j++ {

					to_be_notified = append(to_be_notified, tiles[i].mnt_lst[j])
				}
			}

		}

	}

	special_notifications = ":warning: Special Caution: "
	if to_be_notified != nil {
		for i := 0; i < len(to_be_notified); i++ {

			special_notifications += ("<@" + to_be_notified[i] + ">")
		}
	} else {
		special_notifications += "Not Exist!  "
	}

	general_notifications = ":safety_vest: General Caution: "

	if Ment_Lst != "nil" {
		split := strings.Split(Ment_Lst, ",")
		for i := 0; i < len(split); i++ {
			temp_flag := false

			temp_str := strings.Replace(split[i], "[", "", -1)
			temp_str = strings.Replace(temp_str, "]", "", -1)
			temp_str = strings.Replace(temp_str, "\"", "", -1)
			temp_str = strings.Replace(temp_str, " ", "", -1)

			for j := 0; j < len(to_be_notified); j++ {

				if to_be_notified[j] == temp_str {
					temp_flag = true
					break
				}

			}
			if !temp_flag {

				general_notifications += ("<@" + temp_str + ">")
			}
		}

	} else {
		if notif_channel {
			general_notifications += ("<!channel>")
		} else {
			general_notifications += (" <Off> ")
		}

	}
	group_label := ""
	if findTile(ID).Group {
		group_label = (findTile(ID).GrpLabel)
	}
	str := ""

	if findTile(ID).Type != "HTTP-STATUS" && findTile(ID).Type != "HTTP-RAW" {

		group_label = "|" + group_label
		str += `"` + special_notifications + general_notifications + `"},{ "title":"Error Code" , "value": "` + findTile(ID).error_cd + `"}, {"title":"Related Server" ,  "value": "< ` + findTile(ID).Label + group_label + ` > ` + emoji + `"`
	} else {
		if findTile(ID).Group {
			group_label += " - "
		}
		str += `"` + special_notifications + general_notifications + `"},{ "title":"Error Code" , "value": "` + findTile(ID).error_cd + `"}, {"title":"Related Server" ,  "value": "<` + findTile(ID).Label + `|` + group_label + findTile(ID).Name + `> ` + emoji + `"`
	}

	str = `"fields": [ { "value":  ` + str + `}]}`

	msg := (`{ "pretext": "Monitoror Notification System" ,  `)
	if flag {
		msg += `"color":"DF342F" , `
	} else {
		msg += `"color":"31C55D" , `
	}

	msg += error_caution
	msg += str

	Notify(msg, ID)

}
func attachErCode(ID int, erCode string) {

	for i := 0; i < len(tiles); i++ {
		if tiles[i].ID == ID {
			tiles[i].error_cd = erCode
		}

	}
}
func decreaseErrCount(ID int) {
	for i := 0; i < len(tiles); i++ {
		if tiles[i].ID == ID {

			tiles[i].current_problem_count--

		}

	}
}
func findTile(ID int) Tile {
	var t Tile
	for i := 0; i < len(tiles); i++ {
		if tiles[i].ID == ID {
			t = tiles[i]
		}

	}
	return t
}

func increaseErrCount(ID int) {
	for i := 0; i < len(tiles); i++ {
		if tiles[i].ID == ID {

			tiles[i].current_problem_count++

		}

	}
}
func resetErrCount(ID int) {

	for i := 0; i < len(tiles); i++ {
		if tiles[i].ID == ID {
			tiles[i].current_problem_count = 0
		}

	}
}
func getErrCount(ID int) int {
	var temp int
	for i := 0; i < len(tiles); i++ {
		if tiles[i].ID == ID {
			temp = tiles[i].current_problem_count
		}

	}
	return temp
}

func parser(name string) {

	var err error
	attach_ID := 0
	file, err := ioutil.ReadFile(name)

	var clm Columns

	err = json.Unmarshal(file, &clm)
	if err != nil {
		fmt.Println("error:", err)
	}

	for a := 0; a < len(clm.Tiles); a++ {

		if clm.Tiles[a]["type"] == "GROUP" {

			marshaled, _ := json.Marshal(clm.Tiles[a])
			var tmp_clm Columns
			err = json.Unmarshal(marshaled, &tmp_clm)
			if err != nil {
				fmt.Println("error:", err)
			}
			for b := 0; b < len(tmp_clm.Tiles); b++ {
				attach_ID++

				parser_helper(tmp_clm.Tiles[b], attach_ID, true, tmp_clm.Mention_lst, tmp_clm.Label)
			}

		} else {
			attach_ID++
			parser_helper(clm.Tiles[a], a, false, "", "")
		}

	}

}
func parser_helper(t map[string]interface{}, atc_ID int, isGroup bool, men_List string, grp_lbl string) {
	var single_tile Tile
	var temp_label string
	var temp_type string
	var temp_name string
	var temp_ment_lst string
	var temp_label_2 string
	var temp_label_3 string

	single_tile.ID = atc_ID

	temp_label = fmt.Sprintf("%v", t["params"])

	temp_type = fmt.Sprintf("%v", t["type"])

	temp_name = fmt.Sprintf("%v", t["label"])

	if t["mention_list"] != nil {
		temp_ment_lst = fmt.Sprintf("%v", t["mention_list"])
		temp_ment_lst = strings.Replace(temp_ment_lst, "[", "", -1)
		temp_ment_lst = strings.Replace(temp_ment_lst, "]", "", -1)
		split := strings.Split(temp_ment_lst, ",")
		for i := 0; i < len(split); i++ {
			single_tile.mnt_lst = append(single_tile.mnt_lst, split[i])
		}
	}
	if isGroup {
		single_tile.Group = true
		single_tile.GrpLabel = grp_lbl
		if len(men_List) != 0 {
			split_2 := strings.Split(men_List, ",")
			for i := 0; i < len(split_2); i++ {
				flagger := false
				for j := 0; j < len(single_tile.mnt_lst); j++ {
					if single_tile.mnt_lst[j] == split_2[i] {
						flagger = true
					}
				}
				if !flagger {
					single_tile.mnt_lst = append(single_tile.mnt_lst, split_2[i])
				}

			}
		}
	} else {
		single_tile.Group = false
	}

	if strings.Contains(temp_label, "url:") {
		temp_label = strings.Split(temp_label, "url:")[1]
		temp_label = strings.Split(temp_label, "]")[0]
	} else if strings.Contains(temp_label, "port:") && strings.Contains(temp_label, "hostname:") {
		temp_label_2 = strings.Split(temp_label, " ")[0]
		temp_label_3 = strings.Split(temp_label, " ")[1]

		if strings.Contains(temp_label_2, "hostname:") {
			temp_label_2 = strings.Split(temp_label_2, "hostname:")[1]
			temp_label_3 = strings.Split(temp_label_3, "port:")[1]
			temp_label = temp_label_2 + ":" + temp_label_3
			temp_label = strings.Split(temp_label, "]")[0]
		}

	} else if strings.Contains(temp_label, "hostname:") {
		temp_label = strings.Split(temp_label, "hostname:")[1]
		temp_label = strings.Split(temp_label, "]")[0]
	}

	single_tile.Label = temp_label
	single_tile.Type = temp_type
	single_tile.Name = temp_name
	single_tile.error_cd = ""

	single_tile.current_problem_count = 0
	if single_tile.Type == "HTTP-STATUS" || single_tile.Type == "HTTP-RAW" || single_tile.Type == "PORT" || single_tile.Type == "PING" {
		tiles = append(tiles, single_tile)
	}

}
func Initialization_Notifier(nameConfig string, not_url string, men_list string, emoji bool, prb_cnt int, notf_chnl bool, scs_msg bool) {

	Notif_Url = not_url
	notif_channel = notf_chnl
	Ment_Lst = men_list
	Notif_Emoji = emoji
	problem_count = prb_cnt
	success_msg = scs_msg
	parser(nameConfig)

}

func Notify(msg string, ID int) {
	webhookUrl := Notif_Url

	err := SendSlackNotification(webhookUrl, msg)

	if err != nil {
		decreaseErrCount(ID)
	}

}

func SendSlackNotification(webhookUrl string, msg string) error {

	var srb SlackRequestBody
	err := json.Unmarshal([]byte(msg), &srb)
	if err != nil {
		return err
	}

	out, _ := json.Marshal(&srb)
	req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(out))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 25 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}

func ParaMatcher(label string, isErr bool, ercd string) {
	err_code = ercd
	for i := 0; i < len(tiles); i++ {

		if strings.TrimRight(tiles[i].Label, "\n") == label {

			isErrorOccured(tiles[i].ID, isErr, ercd)
		}

	}
}
