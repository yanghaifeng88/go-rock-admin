package menu

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/daodao97/egin/db"
	"github.com/daodao97/egin/utils/logger"

	"oms/model"
	"oms/service/user"
)

type TreeList struct {
	Id         int         `json:"id"`
	Pid        int         `json:"pid" comment:"父id"`
	ModuleId   int         `json:"module_id" comment:"所属模块"`
	Name       string      `json:"name" comment:"名称 唯一"`
	Type       int         `json:"type" comment:"类型 0原样渲染,1脚手架"`
	PageType   int         `json:"page_type" comment:"类型 0原样渲染,1脚手架"`
	Path       string      `json:"path" comment:"页面路由 唯一"`
	IsShow     bool        `json:"is_show"`
	PageSchema interface{} `json:"page_schema" comment:"页面配置"`
	View       string      `json:"view" comment:"自定义组件路径"`
	Icon       string      `json:"icon" comment:"图片"`
	Children   []TreeList  `json:"children"`
}

type Route struct {
	model.FrontRouteEntity
	IsShow bool `json:"is_show"`
}

func GetMenu(menuList []Route, pid int) []TreeList {
	var treeList []TreeList
	for _, v := range menuList {
		if v.Pid == pid {
			child := GetMenu(menuList, v.Id)
			var pageSchema interface{}
			_ = json.Unmarshal([]byte(v.PageSchema), &pageSchema)
			node := TreeList{
				Id:         v.Id,
				Name:       v.Name,
				Pid:        v.Pid,
				ModuleId:   v.ModuleId,
				Type:       v.Type,
				PageType:   v.PageType,
				PageSchema: pageSchema,
				IsShow:     v.IsShow,
				Path:       v.Path,
				Icon:       v.Icon,
			}
			node.Children = child
			treeList = append(treeList, node)
		}
	}
	return treeList
}

func GetRoutes(uid int) []Route {
	filter := db.Filter{"status": 1}
	if uid > 0 {
		var roleIdsStr []string
		roleIds, _ := user.New().Role(uid)
		for _, v := range roleIds {
			roleIdsStr = append(roleIdsStr, strconv.Itoa(v))
		}
		if !user.New().IsSupperMan(uid) {
			filter["find_in_set"] = strings.Join(roleIdsStr, ",")
		}

	}
	var routes []Route
	err := model.FrontRoute.Select(filter, db.Attr{Select: []string{"id", "pid", "module_id", "name", "type", "page_type", "path", "page_schema", "view", "icon"}}, &routes)
	if err != nil {
		logger.NewLogger("user").Error(err)
	}

	return routes
}

func GetPages(uid int) []TreeList {
	routes := GetRoutes(uid)
	for i, _ := range routes {
		routes[i].IsShow = routes[i].Type != 3
	}
	list := GetMenu(routes, 0)
	return list
}

type TreeSelect struct {
	Value    int          `json:"value"`
	Label    string       `json:"label"`
	Children []TreeSelect `json:"children"`
}

func GetMenuSelect(id int) []TreeSelect {
	userRoute := GetRoutes(id)
	return getMenuSelect(userRoute, 0)
}

func getMenuSelect(menuList []Route, pid int) []TreeSelect {
	var treeList []TreeSelect
	for _, v := range menuList {
		if v.Pid == pid {
			child := getMenuSelect(menuList, v.Id)
			node := TreeSelect{
				Value: v.Id,
				Label: v.Name,
			}
			node.Children = child
			treeList = append(treeList, node)
		}
	}
	return treeList
}

type TreeOptions struct {
	Value    string        `json:"value"`
	Label    string        `json:"label"`
	Children []TreeOptions `json:"children"`
}

func GetMenuOptions(id int) []TreeOptions {
	userRoute := GetRoutes(id)
	return getMenuOptions(userRoute, 0)
}

func getMenuOptions(menuList []Route, pid int) []TreeOptions {
	var treeList []TreeOptions
	for _, v := range menuList {
		if v.Pid == pid {
			child := getMenuOptions(menuList, v.Id)
			node := TreeOptions{
				Value: strconv.Itoa(v.Id),
				Label: v.Name,
			}
			node.Children = child
			if v.PageType == model.PageTypeTable || v.PageType == model.PageTypeForm {
				node.Children = pageSchemaResourceList(v.PageSchema)
			}
			treeList = append(treeList, node)
		}
	}
	return treeList
}

func pageSchemaResourceList(schemaString string) []TreeOptions {
	var pageSchema schema
	_ = json.Unmarshal([]byte(schemaString), &pageSchema)
	var options []TreeOptions

	var formItemOptions []TreeOptions
	if len(pageSchema.FormItems) > 0 {
		for _, v := range pageSchema.FormItems {
			formItemOptions = append(formItemOptions, TreeOptions{Value: v.Field, Label: v.Label})
		}
	}
	if len(formItemOptions) > 0 {
		options = append(options, TreeOptions{
			Value:    "formItems",
			Label:    "表单项",
			Children: formItemOptions,
		})
	}

	if pageSchema.SaveApi != "" {
		options = append(options, TreeOptions{
			Value: "saveApi",
			Label: "保存表单",
		})
	}

	var tableFilterOptions []TreeOptions
	if len(pageSchema.Filter) > 0 {
		for _, v := range pageSchema.Filter {
			tableFilterOptions = append(tableFilterOptions, TreeOptions{Value: v.Field, Label: v.Label})
		}
	}
	if len(tableFilterOptions) > 0 {
		options = append(options, TreeOptions{
			Value:    "filter",
			Label:    "筛选条件",
			Children: tableFilterOptions,
		})
	}
	var tableHeaderOptions []TreeOptions
	if len(pageSchema.Headers) > 0 {
		for _, v := range pageSchema.Headers {
			tableHeaderOptions = append(tableHeaderOptions, TreeOptions{Value: v.Field, Label: v.Label})
		}
	}
	if len(tableHeaderOptions) > 0 {
		options = append(options, TreeOptions{
			Value:    "headers",
			Label:    "列表项",
			Children: tableHeaderOptions,
		})
	}

	var tableNormalButton []TreeOptions
	if len(pageSchema.NormalButton) > 0 {
		for _, v := range pageSchema.NormalButton {
			tableNormalButton = append(tableNormalButton, TreeOptions{Value: v.Target, Label: v.Text})
		}
	}
	if len(tableNormalButton) > 0 {
		options = append(options, TreeOptions{
			Value:    "normalButton",
			Label:    "列表按钮",
			Children: tableNormalButton,
		})
	}

	var tableBatchButton []TreeOptions
	if len(pageSchema.BatchButton) > 0 {
		for _, v := range pageSchema.BatchButton {
			tableBatchButton = append(tableBatchButton, TreeOptions{Value: v.Target, Label: v.Text})
		}
	}
	if len(tableBatchButton) > 0 {
		options = append(options, TreeOptions{
			Value:    "batchButton",
			Label:    "批量按钮",
			Children: tableBatchButton,
		})
	}

	var tableRowButton []TreeOptions
	if len(pageSchema.RowButton) > 0 {
		for _, v := range pageSchema.RowButton {
			tableRowButton = append(tableRowButton, TreeOptions{Value: v.Target, Label: v.Text})
		}
	}
	if len(tableRowButton) > 0 {
		options = append(options, TreeOptions{
			Value:    "rowButton",
			Label:    "行操作按钮",
			Children: tableRowButton,
		})
	}

	return options
}

type filter struct {
	Label string `json:"label"`
	Field string `json:"field"`
}
type header struct {
	Label string `json:"label"`
	Field string `json:"field"`
}
type normalButton struct {
	Text   string `json:"text"`
	Target string `json:"target"`
}
type rowButton struct {
	Text   string `json:"text"`
	Target string `json:"target"`
}
type batchButton struct {
	Text   string `json:"text"`
	Target string `json:"target"`
}
type formItems struct {
	Label string `json:"label"`
	Field string `json:"field"`
}
type schema struct {
	Filter       []filter       `json:"filter"`
	Headers      []header       `json:"headers"`
	NormalButton []normalButton `json:"normalButton"`
	RowButton    []rowButton    `json:"rowButton"`
	BatchButton  []batchButton  `json:"batchButton"`
	FormItems    []formItems    `json:"formItems"`
	SaveApi      string         `json:"saveApi"`
	GetApi       string         `json:"getApi"`
	ListApi      string         `json:"listApi"`
}
