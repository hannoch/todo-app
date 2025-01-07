<template>
  <div class="todo-list">
    <el-row>
      <el-col :span="4">
        <TodoSidebar />
      </el-col>
      <el-col :span="20">
        <div class="filter-bar">
          <el-input
            placeholder="待办名称或者描述"
            v-model="searchQuery"
            @input="handleSearch"
            clearable
            style="width: 200px;"
          />
          <el-select
            v-model="filterStatus"
            placeholder="待办状态"
            @change="filterTodos"
            clearable
            style="width: 150px;"
          >
            <el-option label="全部" value=""></el-option>
            <el-option label="已完成" value="completed"></el-option>
            <el-option label="未完成" value="uncompleted"></el-option>
            <el-option label="重要" value="important"></el-option>
            <el-option label="今日待办" value="today"></el-option>
          </el-select>
          <el-date-picker
            v-model="startDate"
            type="date"
            placeholder="开始日期"
            style="width: 130px;"
            @change="filterTodos"
          />
          <el-date-picker
            v-model="endDate"
            type="date"
            placeholder="结束日期"
            style="width: 130px;"
            @change="filterTodos"
          />
          <el-button type="primary" @click="filterTodos">搜索</el-button>
          <el-button @click="clearFilters">清空</el-button>
        </div>
        <div class="action-bar">
          <el-button type="primary" @click="openTodoDialog">新建待办</el-button>
          <el-button type="danger" @click="bulkDelete" :disabled="!selectedTodos.length">批量删除</el-button>
        </div>
        <el-table 
          :data="todos" 
          style="width: 100%" 
          @selection-change="handleSelectionChange"
          :height="tableHeight"
          v-loading="loading"
        >
          <el-table-column type="selection" width="55"></el-table-column>
          <el-table-column prop="is_completed" label="状态" width="100">
            <template #default="scope">
              <el-checkbox 
                :model-value="scope.row.is_completed" 
                @change="handleStatusChange(scope.row)"
              ></el-checkbox>
            </template>
          </el-table-column>
          <el-table-column prop="is_important" label="重要" width="100">
            <template #default="scope">
              <svg 
                :style="{ fontSize: '20px' }" 
                width="1em" 
                height="1em" 
                viewBox="0 0 24 24"
              >
                <path
                  :fill="scope.row.is_important ? '#b7b703' : 'white'"
                  stroke="black"
                  stroke-width="1"
                  d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
                />
              </svg>
            </template>
          </el-table-column>
          <el-table-column prop="title" label="待办事项">
            <template #default="scope">
              <div style="display: flex; align-items: center; gap: 8px;">
                {{ scope.row.title }}
                <svg 
                  v-if="scope.row.is_daily"
                  width="20px" 
                  height="20px" 
                  viewBox="0 0 24 24"
                  style="color: #409EFF;"
                >
                  <path
                    fill="currentColor"
                    d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 14H4V8l8 5 8-5v10zm-8-7L4 6h16l-8 5z"
                  />
                </svg>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述">
            <template #default="scope">
              <el-tooltip
                v-if="scope.row.description && scope.row.description.length > 30"
                :content="scope.row.description"
                placement="top"
                effect="light"
              >
                <span class="truncate-text">{{ scope.row.description }}</span>
              </el-tooltip>
              <span v-else>{{ scope.row.description }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="due_date" label="截止日期">
            <template #default="scope">
              <span class="date-display">
                {{ formatDate(scope.row.due_date) || '-' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="分类" prop="category_id" min-width="120">
            <template #default="scope">
              {{ getCategoryName(scope.row.category_id) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="scope">
              <el-dropdown trigger="click">
                <el-icon class="el-dropdown-link">
                  <More />
                </el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="viewDetails(scope.row)">查看详情</el-dropdown-item>
                    <el-dropdown-item @click="editTask(scope.row)">编辑</el-dropdown-item>
                    <el-dropdown-item @click="deleteTask(scope.row)">删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>
        </el-table>
        
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="total"
            :current-page="currentPage"
            layout="total, sizes, prev, pager, next"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </el-col>
    </el-row>

    <el-dialog title="待办详情" v-model="detailsDialogVisible" width="500px">
      <div class="todo-details">
        <p><strong>任务:</strong> {{ currentTodo.title }}</p>
        <p><strong>描述:</strong> {{ currentTodo.description }}</p>
        <p><strong>分类:</strong> {{ getCategoryName(currentTodo.category_id) }}</p>
        <p><strong>截止日期:</strong> {{ formatDate(currentTodo.due_date) }}</p>
        <p><strong>提醒时间:</strong> {{ formatDate(currentTodo.reminder_time) }}</p>
        <p><strong>重要:</strong> {{ currentTodo.is_important ? '是' : '否' }}</p>
        <p><strong>每日待办:</strong> {{ currentTodo.is_daily ? '是' : '否' }}</p>
      </div>
      <template #footer>
        <el-button @click="detailsDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <el-dialog :title="isEditing ? '编辑待办' : '新建待办'" v-model="todoDialogVisible">
      <el-form :model="newTodo" :rules="todoRules" ref="todoForm">
        <el-form-item label="待办事项" :label-width="formLabelWidth" prop="title">
          <el-input v-model="newTodo.title" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="描述" :label-width="formLabelWidth">
          <el-input type="textarea" v-model="newTodo.description"></el-input>
        </el-form-item>
        <el-form-item label="所属分类" :label-width="formLabelWidth" prop="category_id">
          <el-select v-model="newTodo.category_id" placeholder="请选择分类">
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="文档链接" :label-width="formLabelWidth">
          <el-input type="textarea" v-model="newTodo.document_link"></el-input>
        </el-form-item>
        <el-form-item label="截止日期" :label-width="formLabelWidth">
          <el-date-picker
            v-model="newTodo.due_date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="提醒时间" :label-width="formLabelWidth">
          <el-date-picker
            v-model="newTodo.reminder_time"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="标记重要" :label-width="formLabelWidth">
          <el-switch v-model="newTodo.is_important"></el-switch>
        </el-form-item>
        <el-form-item label="每日待办" :label-width="formLabelWidth">
          <el-switch v-model="newTodo.is_daily"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="todoDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitTodoForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import TodoSidebar from './TodoSidebar.vue'
import { More } from '@element-plus/icons-vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

// 添加 NodeJS 类型
import type { Timeout } from 'node:timers'

// 定义接口
interface Todo {
  id: number;
  title: string;
  description: string;
  category_id: number;
  document_link: string;
  due_date: string | Date | null;
  reminder_time: string | Date | null;
  is_important: boolean;
  is_daily: boolean;
  is_completed: boolean;
  created_at?: string;
  updated_at?: string;
}

// 定义组件属性
const formLabelWidth = '120px'

// 定义响应式变量
const todoDialogVisible = ref(false)
const detailsDialogVisible = ref(false)
const isEditing = ref(false)
const searchQuery = ref('')
const filterCategory = ref<number | null>(null)
const filterStatus = ref('')
const startDate = ref<Date | null>(null)
const endDate = ref<Date | null>(null)
const selectedTodos = ref<Todo[]>([])
const todos = ref<Todo[]>([])
const loading = ref(true)
const categories = ref<Array<{id: number, name: string}>>([])
const categoryMap = ref<{[key: string]: number}>({})
const reverseCategoryMap = ref<{[key: number]: string}>({})
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const tableHeight = ref('calc(100vh - 280px)')

// 修改搜索防抖定时器类型
let searchTimer: Timeout | null = null;

// 添加搜索处理函数
const handleSearch = () => {
  if (searchTimer) {
    clearTimeout(searchTimer);
  }
  searchTimer = setTimeout(() => {
    fetchTodos();
  }, 300);
};

// 添加筛选处理函数
const filterTodos = () => {
  currentPage.value = 1; // 重置页码
  fetchTodos();
};

// 添加清空筛选条件函数
const clearFilters = () => {
  searchQuery.value = '';
  filterCategory.value = null;
  filterStatus.value = '';
  startDate.value = null;
  endDate.value = null;
  currentPage.value = 1;
  fetchTodos();
};

// 添加打开新建待办对话框函数
const openTodoDialog = () => {
  isEditing.value = false;
  todoDialogVisible.value = true;
};

// 添加批量删除函数
const bulkDelete = async () => {
  try {
    await ElMessageBox.confirm('确定要删除选中的待办事项吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    const ids = selectedTodos.value.map(todo => todo.id);
    await axios.delete('/api/todos/batch', { data: { ids } });
    ElMessage.success('删除成功');
    fetchTodos();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error);
      ElMessage.error('删除失败');
    }
  }
};

// 添加选择变化处理函数
const handleSelectionChange = (selection: Todo[]) => {
  selectedTodos.value = selection;
};

// 添加状态变更处理函数
const handleStatusChange = async (todo: Todo) => {
  try {
    await axios.put(`/api/todos/${todo.id}`, {
      ...todo,
      is_completed: !todo.is_completed
    });
    await fetchTodos();
  } catch (error) {
    console.error('更新状态失败:', error);
    ElMessage.error('更新状态失败');
  }
};

// 添加日期格式化函数
const formatDate = (date: string | Date | null) => {
  if (!date || date === '0001-01-01') return '-';
  if (typeof date === 'string') return date;
  return date.toISOString().split('T')[0];
};

// 添加表单引用
const todoForm = ref();

// 添加表单规则
const todoRules = {
  title: [{ required: true, message: '请输入待办事项', trigger: 'blur' }],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change' }]
};

// 添加获取待办事项列表的方法
const fetchTodos = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/api/todos', {
      params: {
        page: currentPage.value,
        pageSize: pageSize.value,
        query: searchQuery.value,
        status: filterStatus.value,
        category: filterCategory.value,
        startDate: startDate.value ? formatDate(startDate.value) : null,
        endDate: endDate.value ? formatDate(endDate.value) : null
      }
    });
    
    // 确保正确设置数据和总数
    todos.value = response.data.data.list || [];
    total.value = response.data.data.total || 0;
  } catch (error) {
    console.error('获取待办事项失败:', error);
    ElMessage.error('获取待办事项失败');
  } finally {
    loading.value = false;
  }
};

// 添加获取分类列表的方法
const fetchCategories = async () => {
  try {
    const response = await axios.get('/api/categories');
    categories.value = response.data.data;
    // 构建分类映射
    categories.value.forEach(cat => {
      categoryMap.value[cat.name] = cat.id;
      reverseCategoryMap.value[cat.id] = cat.name;
    });
  } catch (error) {
    console.error('获取分类列表失败:', error);
    ElMessage.error('获取分类列表失败');
  }
};

// 添加新待办的默认值
const newTodo = ref({
  title: '',
  description: '',
  category_id: null,
  document_link: '',
  due_date: null,
  reminder_time: null,
  is_important: false,
  is_daily: false,
  is_completed: false
});

// 当前编辑的待办
const currentTodo = ref({
  id: 0,
  title: '',
  description: '',
  category_id: null,
  document_link: '',
  due_date: null,
  reminder_time: null,
  is_important: false,
  is_daily: false,
  is_completed: false
});

// 添加提交表单的方法
const submitTodoForm = async () => {
  if (!todoForm.value) return;
  
  try {
    const valid = await todoForm.value.validate();
    if (!valid) return;

    const todoData = {
      ...newTodo.value,
      due_date: formatDate(newTodo.value.due_date),
      reminder_time: formatDate(newTodo.value.reminder_time)
    };
    console.log(todoData);
    if (isEditing.value) {
      await axios.put(`/api/todos/${currentTodo.value.id}`, todoData);
      ElMessage.success('更新成功');
    } else {
      await axios.post('/api/todos', todoData);
      ElMessage.success('创建成功');
    }

    todoDialogVisible.value = false;
    await fetchTodos();
  } catch (error) {
    console.error('提交表单时出错:', error);
    ElMessage.error('操作失败');
  }
};

// 添加编辑任务的方法
const editTask = (todo: Todo) => {
  currentTodo.value = { ...todo };
  newTodo.value = { ...todo };
  isEditing.value = true;
  todoDialogVisible.value = true;
};

// 添加查看详情的方法
const viewDetails = (todo: Todo) => {
  currentTodo.value = { ...todo };
  detailsDialogVisible.value = true;
};

// 添加删除任务的方法
const deleteTask = async (todo: Todo) => {
  try {
    await ElMessageBox.confirm('确定要删除这个待办事项吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    
    await axios.delete(`/api/todos/${todo.id}`);
    ElMessage.success('删除成功');
    await fetchTodos();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error);
      ElMessage.error('删除失败');
    }
  }
};

// 修改分页处理函数
const handleSizeChange = async (val: number) => {
  pageSize.value = val;
  currentPage.value = 1; // 重置到第一页
  await fetchTodos(); // 重新获取数据
};

const handleCurrentChange = async (val: number) => {
  currentPage.value = val;
  await fetchTodos(); // 重新获取数据
};

// 添加获取分类名称的方法
const getCategoryName = (categoryId: number) => {
  return reverseCategoryMap.value[categoryId] || '未分类';
};

// 组件挂载时获取数据
onMounted(async () => {
  await Promise.all([fetchCategories(), fetchTodos()]);
});
</script>

<style scoped>
.todo-list {
  padding: 0;
  background-color: #f5f5f5;
  min-height: 100vh;
  box-sizing: border-box;
}

.el-row {
  height: 100vh;
}

.el-col {
  height: 100%;
}

.el-col:last-child {
  overflow-y: auto;
  padding: 20px;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}
.action-bar {
  margin-bottom: 20px;
}
.el-dropdown-link {
  cursor: pointer;
  color: #409EFF;
}
.todo-details {
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
.todo-details p {
  margin: 5px 0;
}
.el-table .completed {
  text-decoration: line-through;
  color: #999;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  padding: 16px 0;
}

.truncate-text {
  display: inline-block;
  max-width: 300px; /* 可以根据需要调整宽度 */
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 添加表格容器样式 */
.el-table {
  margin-bottom: 16px;
}

.date-display {
  text-align: center; /* 居中显示 */
  display: inline-block; /* 使其成为块级元素以支持居中 */
  width: 100%; /* 使其占满整个单元格宽度 */
}
</style> 