package common

import (
	"encoding/json"
	"fmt"
	"os"
)

type sidebarMenu struct {
	Sidebar *[]sidebarTitle `json:"sidebar,omitempty"`
}

type sidebarTitle struct {
	Title string          `json:"title,omitempty"`
	Items *[]sidebarItems `json:"items,omitempty"`
}

type sidebarItems struct {
	Label string          `json:"label,omitempty"`
	Url   string          `json:"url,omitempty"`
	Tag   string          `json:"tag,omitempty"`
	Icon  string          `json:"icon,omitpay"`
	Items *[]sidebarItems `json:"items,omitempty"`
}

var sidebar *sidebarMenu

func loadMenu() {
	file, err := os.Open("menu.json")
	if err != nil {
		Danger("Cannot open menu file", err)
	}
	decoder := json.NewDecoder(file)
	sidebar = &sidebarMenu{}
	err = decoder.Decode(sidebar)
	if err != nil {
		Danger("Cannot get sidebar menu from file", err)
	}
}

func getSidebarItems(_items *[]sidebarItems, _tag string, _ulcls string) (tags string) {
	var cls string
	var has bool
	var act bool
	var cl string
	if len(_ulcls) > 0 {
		cls = " class=\"" + _ulcls + "\""
	}
	tags = "<ul" + cls + ">"
	for _, i := range *_items {
		cls, cl = "", ""
		has = i.Items != nil && len(*i.Items) > 0
		act = i.Tag == _tag
		if has {
			cl = " class=\"has-submenu\" @click.stop=\"toggleMenu($event)\""
		}
		tags += "<li" + cl + ">"
		if act {
			cls = " active"
		}
		if has {
			tags += fmt.Sprintf("<span class=\"sb-item%s\">", cls)
		} else {
			tags += fmt.Sprintf("<a href=\"%s\" class=\"sb-item%s\">", i.Url, cls)
		}
		tags += "<span class=\"icon icon-admin\">"
		tags += "<i class=\"" + i.Icon + "\"></i>"
		tags += "</span>"
		tags += "<span class=\"name\">" + i.Label + "</span>"
		if has {
			tags += "</span>"
		} else {
			tags += "</a>"
		}
		if has {
			tags += getSidebarItems(i.Items, _tag, "")
		}
		tags += "</li>"
	}
	tags += "</ul>"
	return
}

func getSidebarHtml(_items *[]sidebarTitle, _tag string) (tags string) {
	for _, v := range *_items {
		tags += "<div class=\"title-admin\">" + v.Title + "</div>"
		if v.Items != nil && len(*v.Items) > 0 {
			tags += getSidebarItems(v.Items, _tag, "zlev")
		}
	}
	return
}
