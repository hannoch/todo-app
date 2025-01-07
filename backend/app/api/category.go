package api

import (
	"todo-go/app/model"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var Category = categoryApi{}

type categoryApi struct{}

// List 获取类别列表
func (a *categoryApi) List(r *ghttp.Request) {
	var categories []model.Category
	if err := g.DB().Model("category").Scan(&categories); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
	}
	r.Response.WriteJsonExit(g.Map{"data": categories})
}

// Create 创建类别
func (a *categoryApi) Create(r *ghttp.Request) {
	var data *model.Category
	if err := r.Parse(&data); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
	}

	// 不需要传入 ID，使用数据库自增
	result, err := g.DB().Model("category").Data(g.Map{
		"name":        data.Name,
		"description": data.Description,
		"parent_id":   data.ParentId,
	}).Insert()

	if err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
		return
	}

	// 获取插入的 ID
	id, err := result.LastInsertId()
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
		return
	}

	// 查询新创建的记录
	var category model.Category
	if err := g.DB().Model("category").Where("id", id).Scan(&category); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
		return
	}

	r.Response.WriteJsonExit(g.Map{"data": category})
}

// Update 更新类别
func (a *categoryApi) Update(r *ghttp.Request) {
	var data *model.Category
	if err := r.Parse(&data); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
	}

	if _, err := g.DB().Model("category").Data(g.Map{
		"name":        data.Name,
		"description": data.Description,
		"parent_id":   data.ParentId,
	}).Where("id", data.Id).Update(); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
		return
	}

	// 查询更新后的记录
	var category model.Category
	if err := g.DB().Model("category").Where("id", data.Id).Scan(&category); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
		return
	}

	r.Response.WriteJsonExit(g.Map{"data": category})
}

// Delete 删除类别
func (a *categoryApi) Delete(r *ghttp.Request) {
	id := r.GetInt("id")

	// 检查是否有待办事项使用此类别
	count, err := g.DB().Model("todo").Where("category_id", id).Count()
	if err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
		return
	}
	if count > 0 {
		r.Response.WriteJsonExit(g.Map{"error": "该类别下还有待办事项，无法删除"})
		return
	}

	if _, err := g.DB().Model("category").Where("id", id).Delete(); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
		return
	}
	r.Response.WriteJsonExit(g.Map{"message": "删除成功"})
}
