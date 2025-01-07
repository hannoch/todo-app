package api

import (
	"todo-go/app/model"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var Todo = todoApi{}

type todoApi struct{}

// List 获取待办事项列表
func (a *todoApi) List(r *ghttp.Request) {
	// 获取分页参数
	page := r.GetInt("page", 1)
	pageSize := r.GetInt("pageSize", 10)
	query := r.GetString("query")

	// 构建查询条件
	model := g.DB().Model("todo")

	if query != "" {
		model = model.Where("title LIKE ? OR description LIKE ?",
			"%"+query+"%",
			"%"+query+"%")
	}

	// 先获取总数
	total, err := model.Count()
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	var todos []struct {
		Id           int    `json:"id"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		CategoryId   int    `json:"category_id"`
		DocumentLink string `json:"document_link"`
		DueDate      string `json:"due_date"`
		ReminderTime string `json:"reminder_time"`
		IsImportant  bool   `json:"is_important"`
		IsDaily      bool   `json:"is_daily"`
		IsCompleted  bool   `json:"is_completed"`
		CreatedAt    string `json:"created_at"`
		UpdatedAt    string `json:"updated_at"`
	}

	// 添加分页查询
	err = model.Limit(pageSize).Offset(offset).Scan(&todos)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    1,
			"message": err.Error(),
		})
		return
	}

	r.Response.WriteJsonExit(g.Map{
		"code":    0,
		"message": "success",
		"data": g.Map{
			"list":     todos,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// Create 创建待办事项
func (a *todoApi) Create(r *ghttp.Request) {
	var data *model.Todo
	if err := r.Parse(&data); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
	}

	// 不需要传入 ID，使用数据库自增
	result, err := g.DB().Model("todo").Data(g.Map{
		"title":         data.Title,
		"description":   data.Description,
		"category_id":   data.CategoryId,
		"document_link": data.DocumentLink,
		"due_date":      data.DueDate,
		"reminder_time": data.ReminderTime,
		"is_important":  data.IsImportant,
		"is_daily":      data.IsDaily,
		"is_completed":  data.IsCompleted,
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
	var todo model.Todo
	if err := g.DB().Model("todo").Where("id", id).Scan(&todo); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
		return
	}

	r.Response.WriteJsonExit(g.Map{"data": todo})
}

// Update 更新待办事项
func (a *todoApi) Update(r *ghttp.Request) {
	var data *model.Todo
	if err := r.Parse(&data); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
	}

	if _, err := g.DB().Model("todo").Data(g.Map{
		"title":         data.Title,
		"description":   data.Description,
		"category_id":   data.CategoryId,
		"document_link": data.DocumentLink,
		"due_date":      data.DueDate,
		"reminder_time": data.ReminderTime,
		"is_important":  data.IsImportant,
		"is_daily":      data.IsDaily,
		"is_completed":  data.IsCompleted,
	}).Where("id", data.Id).Update(); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
		return
	}

	// 查询更新后的记录
	var todo model.Todo
	if err := g.DB().Model("todo").Where("id", data.Id).Scan(&todo); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
		return
	}

	r.Response.WriteJsonExit(g.Map{"data": todo})
}

// Delete 删除待办事项
func (a *todoApi) Delete(r *ghttp.Request) {
	// 获取待办事项 ID /api/delete/{id}
	id := r.GetInt("id")
	if _, err := g.DB().Model("todo").Where("id", id).Delete(); err != nil {
		r.Response.WriteJsonExit(g.Map{"error": err.Error()})
	}
	r.Response.WriteJsonExit(g.Map{"message": "删除成功"})
}
