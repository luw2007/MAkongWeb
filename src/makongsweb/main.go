// Copyright (C) 2013 Andras Belicza. All rights reserved.
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// libmakongweb project main.go
package main

import (
	"code.google.com/p/gowut/gwu"
	"github.com/bitly/go-simplejson"

	"fmt"
	"libmakong"
	"os"
	"strconv"
)

var post_data libmakong.User
var basePostData string
var userFile string
var username string

// plural returns an empty string if i is equal to 1,
// "s" otherwise.
func plural(i int) string {
	if i == 1 {
		return ""
	}
	return "s"
}

func Version() string {
	version := "0.4"
	return version
}

type demo struct {
	link      gwu.Label
	buildFunc func(gwu.Event) gwu.Comp
	comp      gwu.Comp // Lazily initialized demo comp
}
type pdemo *demo

func buildLoginWin(s gwu.Session) {
	win := gwu.NewWindow("login", "登录")
	win.Style().SetFullSize()
	win.SetAlign(gwu.HA_CENTER, gwu.VA_MIDDLE)

	p := gwu.NewPanel()
	p.SetHAlign(gwu.HA_CENTER)
	p.SetCellPadding(2)

	l := gwu.NewLabel("登  录")
	l.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD).SetFontSize("150%")
	p.Add(l)

	errL := gwu.NewLabel("")
	errL.Style().SetColor(gwu.CLR_RED)
	p.Add(errL)

	table := gwu.NewTable()
	table.SetCellPadding(2)
	table.EnsureSize(2, 2)
	table.Add(gwu.NewLabel("手机号:"), 0, 0)
	tb := gwu.NewTextBox("")
	tb.Style().SetWidthPx(160)
	table.Add(tb, 0, 1)

	table.Add(gwu.NewLabel("昵称:"), 1, 0)
	pb := gwu.NewTextBox("白人")
	pb.Style().SetWidthPx(160)
	table.Add(pb, 1, 1)

	table.Add(gwu.NewLabel("第几区:"), 2, 0)
	lb := gwu.NewListBox([]string{"一区", "二区"})
	lb.SetSelected(1, true)
	lb.AddEHandlerFunc(func(e gwu.Event) {
		table.Add(gwu.NewLabel(lb.SelectedValue()+":"), 2, 0)
		e.ReloadWin("login")
	}, gwu.ETYPE_CHANGE)
	table.Add(lb, 2, 1)

	p.Add(table)
	b := gwu.NewButton("OK")
	b.AddEHandlerFunc(func(e gwu.Event) {
		phone := tb.Text()
		userFile = tb.Text() + "_user.json"
		username = pb.Text()
		groupId := lb.SelectedIdx() + 1

		_basePostData, err := libmakong.Read_config(userFile, &post_data)
		basePostData = _basePostData

		fmt.Println(basePostData, "1", username, userFile, post_data)
		if err != nil {
			post_data, basePostData, err = libmakong.Get_post_data(userFile, phone, username, groupId)
		}

		fmt.Println(basePostData, "2", username, userFile, post_data)
		fmt.Println("login:", tb.Text(), pb.Text(), groupId, post_data, basePostData)
		if err != nil || post_data.UserId != "" {
			e.Session().RemoveWin(win) // Login win is removed, password will not be retrievable from the browser

			buildPrivateWins(e.Session())

			e.ReloadWin("main")
		} else {
			e.SetFocusedComp(tb)
			errL.SetText("登录失败")
			e.MarkDirty(errL)
		}
	}, gwu.ETYPE_CLICK)
	p.Add(b)
	l = gwu.NewLabel("")
	p.Add(l)
	p.CellFmt(l).Style().SetHeightPx(200)

	win.Add(p)
	s.AddWin(win)
}

func buildPrivateWins(sess gwu.Session) {
	win := gwu.NewWindow("main", "main - MAKong's")
	win.Style().SetFullSize()
	win.AddEHandlerFunc(func(e gwu.Event) {
		switch e.Type() {
		case gwu.ETYPE_WIN_LOAD:
			fmt.Println("LOADING window:", e.Src().Id())
		case gwu.ETYPE_WIN_UNLOAD:
			fmt.Println("UNLOADING window:", e.Src().Id())
		}
	}, gwu.ETYPE_WIN_LOAD, gwu.ETYPE_WIN_UNLOAD)

	hiddenPan := gwu.NewNaturalPanel()
	sess.SetAttr("hiddenPan", hiddenPan)

	header := gwu.NewHorizontalPanel()
	header.Style().SetFullWidth().SetBorderBottom2(2, gwu.BRD_STYLE_SOLID, "#777777")
	l := gwu.NewLabel("MA控 - 简单版 MAKong's " + Version())
	l.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD).SetFontSize("120%")
	header.Add(l)
	header.AddHConsumer()
	header.Add(gwu.NewLabel("Theme:"))
	themes := gwu.NewListBox([]string{"default", "debug"})
	themes.AddEHandlerFunc(func(e gwu.Event) {
		win.SetTheme(themes.SelectedValue())
		e.ReloadWin("main")
	}, gwu.ETYPE_CHANGE)
	header.Add(themes)
	header.AddHSpace(10)
	header.Add(gwu.NewLabel(username))
	reset := gwu.NewLink("退出", "#")
	reset.SetTarget("")
	reset.AddEHandlerFunc(func(e gwu.Event) {
		e.RemoveSess()
		e.Session().RemoveWin(win)
		buildLoginWin(e.Session())
		e.ReloadWin("login")
	}, gwu.ETYPE_CLICK)
	header.Add(reset)
	setNoWrap(header)
	win.Add(header)

	content := gwu.NewHorizontalPanel()
	content.SetCellPadding(1)
	content.SetVAlign(gwu.VA_TOP)
	content.Style().SetFullSize()

	demoWrapper := gwu.NewPanel()
	demoWrapper.Style().SetPaddingLeftPx(5)
	demoWrapper.AddVSpace(10)
	demoTitle := gwu.NewLabel("")
	demoTitle.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD).SetFontSize("110%")
	demoWrapper.Add(demoTitle)
	demoWrapper.AddVSpace(10)

	links := gwu.NewPanel()
	links.SetCellPadding(1)
	links.Style().SetPaddingRightPx(5)

	demos := make(map[string]pdemo)
	var selDemo pdemo

	selectDemo := func(d pdemo, e gwu.Event) {
		if selDemo != nil {
			selDemo.link.Style().SetBackground("")
			if e != nil {
				e.MarkDirty(selDemo.link)
			}
			demoWrapper.Remove(selDemo.comp)
		}
		selDemo = d
		d.link.Style().SetBackground("#88ff88")
		demoTitle.SetText(d.link.Text())
		if d.comp == nil {
			d.comp = d.buildFunc(e)
		}
		demoWrapper.Add(d.comp)
		if e != nil {
			e.MarkDirty(d.link, demoWrapper)
		}
	}

	createDemo := func(name string, buildFunc func(gwu.Event) gwu.Comp) pdemo {
		link := gwu.NewLabel(name)
		link.Style().SetFullWidth().SetCursor(gwu.CURSOR_POINTER).SetDisplay(gwu.DISPLAY_BLOCK).SetColor(gwu.CLR_BLUE)
		demo := &demo{link: link, buildFunc: buildFunc}
		link.AddEHandlerFunc(func(e gwu.Event) {
			selectDemo(demo, e)
		}, gwu.ETYPE_CLICK)
		links.Add(link)
		demos[name] = demo
		return demo
	}

	links.Style().SetFullHeight().SetBorderRight2(2, gwu.BRD_STYLE_SOLID, "#777777")
	links.AddVSpace(5)
	aboutPage := createDemo("介绍", buildAboutPage)
	selectDemo(aboutPage, nil)
	links.AddVSpace(5)
	l = gwu.NewLabel("功能")
	l.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD).SetFontSize("110%")
	links.Add(l)
	links.AddVSpace(5)
	createDemo("用户信息", buildUserInfoPage)
	createDemo("好友信息", buildFriendsPage)
	createDemo("收集排行榜", buildRankingPage)

	links.AddVSpace(5)
	l = gwu.NewLabel("TODO")
	l.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD).SetDisplay(gwu.DISPLAY_BLOCK)
	links.Add(l)
	createDemo("模拟器", buildComingsoonDemo)
	createDemo("卡片浏览", buildComingsoonDemo)
	createDemo("模拟器", buildComingsoonDemo)
	links.AddVConsumer()
	setNoWrap(links)
	content.Add(links)
	content.Add(demoWrapper)
	content.CellFmt(demoWrapper).Style().SetFullWidth()

	win.Add(content)
	win.CellFmt(content).Style().SetFullSize()

	footer := gwu.NewHorizontalPanel()
	footer.Style().SetFullWidth().SetBorderTop2(2, gwu.BRD_STYLE_SOLID, "#777777")
	footer.Add(hiddenPan)
	footer.AddHConsumer()
	l = gwu.NewLabel("Copyright © 2013 luw2007. All rights reserved.")
	l.Style().SetFontStyle(gwu.FONT_STYLE_ITALIC).SetFontSize("95%")
	footer.Add(l)
	footer.AddHSpace(10)
	link := gwu.NewLink("Thanks Gowut", "https://sites.google.com/site/gowebuitoolkit/")
	link.Style().SetFontStyle(gwu.FONT_STYLE_ITALIC).SetFontSize("95%")
	footer.Add(link)
	setNoWrap(footer)
	win.Add(footer)

	sess.AddWin(win)
}

func buildAboutPage(event gwu.Event) gwu.Comp {
	p := gwu.NewPanel()
	html := "<p>使用ma控未公开的接口， 查询数据。</p><h3>目前完成：</h3><ol><li>用户信息显示</li><li>手机排名</li><li>不活跃的小伙伴</li></ol><h3>TODO:</h3><ul><li>合成模拟器</li><li>闪卡提醒</li><li>卡牌攻击力计算</li></ul><p>需要增加的功能，请在 <a href=\"http://42qu.cc/luw2007\">http://42qu.cc/luw2007</a>  留言</p>"
	h := gwu.NewHtml(html)
	p.Add(h)
	return p
}

func buildUserInfoPage(event gwu.Event) gwu.Comp {
	p := gwu.NewPanel()
	fmt.Println(basePostData, "c", username, userFile, post_data)
	if len(basePostData) > 0 {
		user_info := libmakong.Get_user_info(basePostData)
		js, err := simplejson.NewJson(user_info)
		if err != nil {
			fmt.Println("json err:", err)
		}
		info, err := js.Get("data").Map()
		if err != nil {
			fmt.Println("json err:", err)
		}
		html := "<table border=0.1>"

		for k, v := range info {
			html += "<tr><td>" + k + "</td>: <td>"

			switch vv := v.(type) {
			case int:
				html += strconv.Itoa(vv)
			case string:
				html += vv
			case float64:
				html += strconv.Itoa(int(vv))
			default:
				fmt.Println(k, "is of a type I don't know how to handle", vv)
				html += "</td>"
			}
			html += "</td></tr>"

		}
		html += "</tr></table)"
		h := gwu.NewHtml(html)
		p.Add(h)
	} else {

		p.Add(gwu.NewHtml("error"))
	}
	return p
}

func buildFriendsPage(event gwu.Event) gwu.Comp {
	p := gwu.NewPanel()
	if len(basePostData) == 0 {
		p.Add(gwu.NewHtml("error"))
	}

	data := libmakong.Get_friends_info(basePostData)
	js, err := simplejson.NewJson(data)
	if err != nil {
		fmt.Println("json err:", err)
	}
	friends, err := js.Get("data").Get("friends").Array()
	if err != nil {
		fmt.Println("json err:", err)
	}
	html := "<table><tr><td>昵称</td><td><级别></td><td><活跃度></td><td><好友数></td></tr>"

	//name	level active now/max

	for _, friend := range friends {
		m := friend.(map[string]interface{})
		html += "<tr><td>" +
			m["name"].(string) + "</td><td>" +
			m["level"].(string) + "</td><td>" +
			m["active"].(string) + "</td><td>" +
			m["now"].(string) + "//" + m["max"].(string) + "</td></tr>"

	}
	html += "</table>"
	h := gwu.NewHtml(html)
	p.Add(h)
	return p

}

func buildRankingPage(event gwu.Event) gwu.Comp {
	p := gwu.NewPanel()
	query_string := post_data.Phone + "&groupId=" + strconv.Itoa(post_data.GroupId) + "&serverNo=" + post_data.ServerNo + "&userId=" + post_data.UserId
	l := gwu.NewLabel("少女之心")
	l.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD).SetDisplay(gwu.DISPLAY_BLOCK)
	p.Add(l)

	_url := "http://www.niuxba.com/ma/backend/cgi-bin/getZsjRankInfo.php?phone=" + query_string
	p.Add(gwu.NewLink("我的排名", _url))
	_url = "http://www.niuxba.com/ma/backend/cgi-bin/getZsjFriendRank.php?phone=" + query_string
	p.Add(gwu.NewLink("好友排名", _url))

	l = gwu.NewLabel("比基尼")
	l.Style().SetFontWeight(gwu.FONT_WEIGHT_BOLD).SetDisplay(gwu.DISPLAY_BLOCK)
	p.Add(l)
	p.Add(l)
	_url = "http://www.niuxba.com/ma/backend/cgi-bin/getRankInfo.php?phone=" + query_string
	p.Add(gwu.NewLink("我的排名", _url))
	_url = "http://www.niuxba.com/ma/backend/cgi-bin/getFriendRank.php?phone=" + query_string
	p.Add(gwu.NewLink("好友排名", _url))

	return p
}

func buildComingsoonDemo(event gwu.Event) gwu.Comp {
	p := gwu.NewPanel()
	p.Add(gwu.NewLabel("coming soon!"))

	return p

}

// setNoWrap sets WHITE_SPACE_NOWRAP to all children of the specified panel.
func setNoWrap(panel gwu.Panel) {
	count := panel.CompsCount()
	for i := count - 1; i >= 0; i-- {
		panel.CompAt(i).Style().SetWhiteSpace(gwu.WHITE_SPACE_NOWRAP)
	}
}

// SessHandler is our session handler to build the showcases window.
type SessHandler struct{}

func (h SessHandler) Created(s gwu.Session) {
	buildLoginWin(s)
	//buildShowcaseWin(s)
}

func (h SessHandler) Removed(s gwu.Session) {}

func main() {
	// Allow app control from command line (in co-operation with the starter script):
	fmt.Println("Type 'r' to restart, 'e' to exit.")
	go func() {
		var cmd string
		for {
			fmt.Scanf("%s", &cmd)
			switch cmd {
			case "r": // restart
				os.Exit(1)
			case "e": // exit
				os.Exit(0)
			}
		}
	}()

	// Create GUI server
	server := gwu.NewServer("", "")
	server.SetText("MAKong's - simple makong on pc")

	server.AddSessCreatorName("login", "登录")

	server.AddSHandler(SessHandler{})

	// Start GUI server
	if err := server.Start("login"); err != nil {
		fmt.Println("Error: Cound not start GUI server:", err)
		return
	}
}
