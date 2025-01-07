<template>
  <div class="todo-sidebar">
    <div class="header">
      <span class="title">待办类目</span>
      <el-button type="primary" @click="openNewCategoryDialog" size="small">新建类目</el-button>
    </div>
    <el-menu
      v-loading="loading"
      default-active="1"
      class="el-menu-vertical-demo"
      @select="handleSelect"
    >
      <template v-for="(category, index) in categories" :key="index">
        <el-sub-menu v-if="category.children && category.children.length" :index="String(index + 1)">
          <template #title>
            <span>{{ category.name }}</span>
            <el-dropdown trigger="hover" @command="(command) => handleCommand(command, index)">
              <el-icon class="el-dropdown-link">
                <More />
              </el-icon>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">编辑</el-dropdown-item>
                  <el-dropdown-item command="delete">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <el-menu-item 
            v-for="(child, childIndex) in category.children" 
            :key="`${index}-${childIndex}`"
            :index="`${index + 1}-${childIndex + 1}`"
          >
            <span>{{ child.name }}</span>
            <el-dropdown trigger="hover" @command="(command) => handleChildCommand(command, index, childIndex)">
              <el-icon class="el-dropdown-link">
                <More />
              </el-icon>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">编辑</el-dropdown-item>
                  <el-dropdown-item command="delete">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </el-menu-item>
        </el-sub-menu>
        <el-menu-item v-else :index="String(index + 1)">
          <span>{{ category.name }}</span>
          <el-dropdown trigger="hover" @command="(command) => handleCommand(command, index)">
            <el-icon class="el-dropdown-link">
              <More />
            </el-icon>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="edit">编辑</el-dropdown-item>
                <el-dropdown-item command="delete">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-menu-item>
      </template>
    </el-menu>

    <el-dialog :title="isEditing ? '编辑类目' : '新建类目'" v-model="dialogVisible">
      <el-form :model="newCategory" :rules="categoryRules" ref="categoryForm">
        <el-form-item label="分类名称" :label-width="formLabelWidth" prop="name">
          <el-input v-model="newCategory.name" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="描述" :label-width="formLabelWidth">
          <el-input type="textarea" v-model="newCategory.description"></el-input>
        </el-form-item>
        <el-form-item label="父级分类" :label-width="formLabelWidth">
          <el-select
            v-model="newCategory.parent_id"
            placeholder="请选择父级分类"
            clearable
          >
            <el-option
              v-for="cat in categories"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCategoryForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from 'vue';
import { More } from '@element-plus/icons-vue';
import axios from 'axios';
import { ElMessage } from 'element-plus';

interface Category {
  id: number;
  name: string;
  description: string;
  parent_id?: number;
  children?: Category[];
}

interface NewCategory {
  name: string;
  description: string;
  parent_id: number | undefined;
}

export default defineComponent({
  name: 'TodoSidebar',
  components: {
    More,
  },
  setup() {
    const dialogVisible = ref(false);
    const isEditing = ref(false);
    const loading = ref(true);
    const categories = ref<Category[]>([]);
    const categoryForm = ref();
    const currentCategory = ref<Category | null>(null);
    
    const categoryRules = {
      name: [
        { required: true, message: '请输入分类名称', trigger: 'blur' },
      ]
    };

    const newCategory = ref<NewCategory>({
      name: '',
      description: '',
      parent_id: undefined,
    });

    // 处理子分类的操作
    const handleChildCommand = async (command: string, parentIndex: number, childIndex: number) => {
      const parentCategory = categories.value[parentIndex];
      if (!parentCategory.children) return;
      
      const childCategory = parentCategory.children[childIndex];
      
      if (command === 'edit') {
        isEditing.value = true;
        currentCategory.value = childCategory;
        newCategory.value = {
          name: childCategory.name,
          description: childCategory.description,
          parent_id: childCategory.parent_id,
        };
        dialogVisible.value = true;
      } else if (command === 'delete') {
        try {
          await axios.delete(`/api/categories/${childCategory.id}`);
          ElMessage.success('删除成功');
          await fetchCategories();
        } catch (error) {
          console.error('删除失败:', error);
          ElMessage.error('删除失败');
        }
      }
    };

    // 处理父分类的操作
    const handleCommand = async (command: string, index: number) => {
      const category = categories.value[index];
      
      if (command === 'edit') {
        isEditing.value = true;
        currentCategory.value = category;
        newCategory.value = {
          name: category.name,
          description: category.description,
          parent_id: category.parent_id,
        };
        dialogVisible.value = true;
      } else if (command === 'delete') {
        try {
          await axios.delete(`/api/categories/${category.id}`);
          ElMessage.success('删除成功');
          await fetchCategories();
        } catch (error) {
          console.error('删除失败:', error);
          ElMessage.error('删除失败');
        }
      }
    };

    // 添加 handleSelect 方法
    const handleSelect = (index: string) => {
      console.log('选中的菜单项:', index);
      // 这里可以添加选中菜单项后的处理逻辑
    };

    // 获取分类列表
    const fetchCategories = async () => {
      try {
        loading.value = true;
        const response = await axios.get('/api/categories');
        const rawCategories = response.data.data;
        categories.value = buildCategoryTree(rawCategories);
      } catch (error) {
        console.error('获取分类失败:', error);
        ElMessage.error('获取分类失败');
      } finally {
        loading.value = false;
      }
    };

    // 构建分类树结构
    const buildCategoryTree = (items: Category[]): Category[] => {
      const result: Category[] = [];
      const map = new Map<number, Category>();

      items.forEach(item => {
        map.set(item.id, { ...item, children: [] });
      });

      items.forEach(item => {
        const node = map.get(item.id);
        if (node && item.parent_id && map.has(item.parent_id)) {
          const parent = map.get(item.parent_id);
          if (parent) {
            parent.children = parent.children || [];
            parent.children.push(node);
          }
        } else if (node) {
          result.push(node);
        }
      });

      return result;
    };

    // 打开新建分类对话框
    const openNewCategoryDialog = () => {
      isEditing.value = false;
      newCategory.value = { name: '', description: '', parent_id: undefined };
      dialogVisible.value = true;
    };

    // 重置表单
    const resetForm = () => {
      newCategory.value = {
        name: '',
        description: '',
        parent_id: undefined,
      };
      currentCategory.value = null;
      isEditing.value = false;
    };

    // 提交分类表单
    const submitCategoryForm = async () => {
      if (!categoryForm.value) {
        console.log('表单引用不存在');
        return;
      }
      
      try {
        console.log('开始验证表单');
        const valid = await categoryForm.value.validate();
        console.log('表单验证结果:', valid);
        
        if (!valid) {
          return;
        }
        
        console.log('编辑状态:', isEditing.value);
        console.log('当前分类:', currentCategory.value);
        
        if (isEditing.value && currentCategory.value) {
          const updateData = {
            id: currentCategory.value.id,
            name: newCategory.value.name,
            description: newCategory.value.description,
            parent_id: newCategory.value.parent_id
          };
          console.log('准备更新数据:', updateData);
          
          await axios.put('/api/categories', updateData);
          ElMessage.success('更新成功');
        } else {
          // 创建分类
          await axios.post('/api/categories', {
            name: newCategory.value.name,
            description: newCategory.value.description,
            parent_id: newCategory.value.parent_id
          });
          ElMessage.success('创建成功');
        }

        // 重新获取分类列表
        await fetchCategories();
        dialogVisible.value = false;
        resetForm();
      } catch (error: unknown) {
        console.error('提交表单时出错:', error);
        if (error instanceof Error) {
          ElMessage.error(error.message);
        } else {
          ElMessage.error('操作失败');
        }
      }
    };

    // 在组件挂载时获取分类列表
    onMounted(() => {
      fetchCategories();
    });

    return {
      dialogVisible,
      isEditing,
      newCategory,
      categories,
      loading,
      categoryRules,
      categoryForm,
      currentCategory,
      openNewCategoryDialog,
      submitCategoryForm,
      handleCommand,
      handleChildCommand,
      handleSelect,
      formLabelWidth: '120px',
    };
  },
});
</script>

<style scoped>
.todo-sidebar {
  height: 100%;
  padding: 20px;
  background-color: #f9f9f9;
  box-shadow: none;
  border-radius: 0;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}
.title {
  font-weight: bold;
  font-size: 16px;
}
.el-menu-vertical-demo {
  width: 100%;
  border-right: none;
}
.el-dropdown-link {
  cursor: pointer;
  color: #409EFF;
  margin-left: 8px;
}

.el-sub-menu :deep(.el-sub-menu__title) {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.el-menu-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style> 