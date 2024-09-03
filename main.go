package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/TEST-SDLJWNKELK/notice/pkg/dingtalk"
	"github.com/TEST-SDLJWNKELK/notice/pkg/holiday"
)

var (
	AccessToken = os.Getenv("DING_ACCESS_TOKEN")
	Secret      = os.Getenv("DING_SECRET")
	JuHeKey     = os.Getenv("JU_HE_KEY")
)

var template = "明日社区值班工程师 %s 运营伙伴 %s 辛苦大家关注社群消息，及时回复哦～"

func main() {
	//设定北京时间
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("无法加载时区:", err)
		return
	}

	t := time.Now().In(location).AddDate(0, 0, 1)

	flag := holiday.IsHoliday(JuHeKey, t)

	if !flag {
		send(t)
	}
}

func send(t time.Time) {
	d := dingtalk.Webhook{
		AccessToken: AccessToken,
		Secret:      Secret,
	}
	t.Weekday()
	g := duty[t.Weekday()]

	engineer := "@" + strings.Join(g.engineer, " @")
	devrels := "@" + strings.Join(g.devrels, " @")

	all := append(g.engineer, g.devrels...)

	if err := d.SendMessageText(fmt.Sprintf(template, engineer, devrels), all...); err != nil {
		fmt.Println(err)
	}

}

var duty = map[time.Weekday]group{
	time.Monday: {
		engineer: []string{"18434391952", "13265162439"},
		devrels:  []string{"13161354391"},
	},
	time.Tuesday: {
		engineer: []string{"18088642209", "15511340988"},
		devrels:  []string{"15354874060"},
	},
	time.Wednesday: {
		engineer: []string{"17600686802", "15294561913"},
		devrels:  []string{"17624047637"},
	},
	time.Thursday: {
		engineer: []string{"15901359231", "18346072982"},
		devrels:  []string{"13161354391"},
	},
	time.Friday: {
		engineer: []string{"18501154050", "18332559886"},
		devrels:  []string{"15354874060"},
	},
}

type group struct {
	engineer []string
	devrels  []string
}
