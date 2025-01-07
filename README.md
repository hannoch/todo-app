# Todo List 待办事项管理系统

一个基于 Vue 3 + TypeScript + Element Plus 的现代化待办事项管理系统。

## 功能特点

- 📝 待办事项的增删改查
- 📋 分类管理
- 🔍 多条件筛选
- 📅 日期范围查询
- 📊 分页展示
- ⭐ 重要性标记
- 🔄 状态切换
- 📎 文档链接
- ⏰ 提醒功能
- 🎯 每日任务

## 技术栈

- Vue 3 - 渐进式 JavaScript 框架
- TypeScript - JavaScript 的超集
- Element Plus - 基于 Vue 3 的组件库
- Axios - HTTP 客户端
- Vue Router - 路由管理
- date-fns - 日期处理库

## 项目结构

src/
├── components/ # 组件目录
│ ├── TodoList.vue # 待办事项列表组件
│ └── TodoSidebar.vue # 侧边栏组件
├── router/ # 路由配置
├── types/ # TypeScript 类型定义
├── utils/ # 工具函数
└── App.vue # 根组件

backend/
├── app/ # 应用逻辑
│ ├── api/ # API 接口
│ └── model/ # 数据模型
├── main.go # 主入口文件
└── manifest/ # 配置文件
    ├── config.toml # 配置文件
    └── ... # 其他配置文件

## API 接口

### 待办事项

- `GET /api/todos` - 获取待办事项列表
- `POST /api/todos` - 创建新待办事项
- `PUT /api/todos/:id` - 更新待办事项
- `DELETE /api/todos/:id` - 删除待办事项
- `DELETE /api/todos/batch` - 批量删除待办事项

### 分类

- `GET /api/categories` - 获取分类列表
- `POST /api/categories` - 创建新分类
- `PUT /api/categories/:id` - 更新分类
- `DELETE /api/categories/:id` - 删除分类

## 配置说明

### 环境变量
VITE_API_BASE_URL=http://localhost:8080/api

### 分页配置

- 默认每页显示：10条
- 可选每页显示：10, 20, 50, 100条

## 开发指南

1. 代码规范
   - 使用 ESLint 进行代码检查
   - 使用 Prettier 进行代码格式化
   - 遵循 Vue 3 组合式 API 风格指南

2. 组件开发
   - 使用 `<script setup>` 语法
   - 使用 TypeScript 类型注解
   - 遵循单一职责原则

3. 提交规范
   - feat: 新功能
   - fix: 修复
   - docs: 文档更新
   - style: 代码格式
   - refactor: 重构
   - test: 测试
   - chore: 构建过程或辅助工具的变动

## 注意事项

1. 确保后端 API 服务正常运行
2. 正确配置跨域访问
3. 注意处理大数据量的性能优化
4. 做好错误处理和用户提示

## 贡献指南

1. Fork 本仓库
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

[MIT License](LICENSE)
