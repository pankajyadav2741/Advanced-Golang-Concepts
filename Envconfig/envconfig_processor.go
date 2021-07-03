/*
	Alarm which are included need to be forwarded to some other service.
	Purpose of this program is to show how to read environment variables and how to decode them into map
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type AlarmData struct {
	Alarm map[string]bool `envconfig:"ALARM_DATA"`
}

var (
	alarm1 = `{"alarmID":"1001","severity":"minor"}`    //EXISTENT Alarm
	alarm2 = `{"alarmID":"1002","severity":"major"}`    //EXISTENT Alarm
	alarm3 = `{"alarmID":"1003","severity":"critical"}` //EXISTENT Alarm
	alarm4 = `{"alarmID":"1004","severity":"minor"}`    //EXISTENT Alarm
	alarm5 = `{"alarmID":"1005","severity":"major"}`    //EXISTENT Alarm
	alarm6 = `{"alarmID":"1006","severity":"minor"}`    //NON-EXISTENT Alarm: Considered as excluded
	alarm7 = `{"alarmID":"1007","severity":"critical"}` //NON-EXISTENT Alarm: Considered as excluded
)

func prepareAlarm(msg []byte) map[string]interface{} {
	alarm := map[string]interface{}{}
	err := json.Unmarshal([]byte(msg), &alarm)
	if err != nil {
		panic("Error in unmarshalling alarm")
		return nil
	}
	return alarm
}

func processAlarm(alarm []byte, alarmData map[string]bool) {
	al := prepareAlarm([]byte(alarm))
	if !isAlarmExcluded(al, alarmData) {
		fmt.Printf("INFO: Alarm %s excluded\n", al["alarmID"])
	} else {
		fmt.Printf("INFO: Alarm %s included\n", al["alarmID"])
	}
}

func main() {
	//Alarm IDs with value true are excluded
	os.Setenv("ALARM_DATA", "1001:false,1002:true,1003:false,1004:true,1005:true") //There must not be any spaces else envconfig.Process will not be able to decode to map.
	al := AlarmData{}
	err := envconfig.Process("", &al)
	if err != nil {
		panic("Error in processing environment variable")
	}
	processAlarm([]byte(alarm1), al.Alarm)
	processAlarm([]byte(alarm2), al.Alarm)
	processAlarm([]byte(alarm3), al.Alarm)
	processAlarm([]byte(alarm4), al.Alarm)
	processAlarm([]byte(alarm5), al.Alarm)
	processAlarm([]byte(alarm6), al.Alarm)
	processAlarm([]byte(alarm7), al.Alarm)
}

func isAlarmExcluded(alarmAttributes map[string]interface{}, alarmData map[string]bool) bool {
	alarmId := fmt.Sprintf("%v", alarmAttributes["alarmID"])
	val, _ := alarmData[alarmId]
	return val
}
