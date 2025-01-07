import { createStore } from 'vuex'

// 定义 Todo 接口
interface Todo {
  id: number;
  task: string;
  description: string;
  category: string;
  link: string;
  deadline: string;
  reminder: string;
  important: boolean;
  daily: boolean;
  status: boolean;
}

// 定义 State 接口
interface State {
  todos: Todo[];
}

export default createStore<State>({
  state: {
    todos: [] // 现在这是一个 Todo[] 类型的数组
  },
  getters: {
  },
  mutations: {
    setTodos(state, todos: Todo[]) {
      state.todos = todos;
    },
    updateTodoStatus(state, todo: Todo) {
      const index = state.todos.findIndex(t => t.id === todo.id);
      if (index !== -1) {
        state.todos[index] = { ...todo };
      }
    }
  },
  actions: {
    fetchTodos({ commit }) {
      // 这里可以添加获取待办事项的逻辑
      // const response = await fetch('/api/todos');
      // const todos = await response.json();
      // commit('setTodos', todos);
    }
  },
  modules: {
  }
})
