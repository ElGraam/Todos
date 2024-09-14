/* eslint-disable @typescript-eslint/naming-convention */
import React, { useState, useEffect, memo } from 'react';
import axios from 'axios';
import './App.css';

const apiUrl = 'http://localhost:5000/api';

interface Todo {
  id: number;
  body: string;
  completed: boolean;
  created_at: string | null;
}

interface TodoItemProps {
  todo: Todo;
  onToggle: (id: number) => void;
  onDelete: (id: number) => void;
}

const TodoItem: React.FC<TodoItemProps> = ({ todo, onToggle, onDelete }) => (
  <li className={`todo-item ${todo.completed ? 'completed' : ''}`}>
    <span className="todo-text">{todo.body}</span>
    <div className="todo-actions">
      <button onClick={() => onToggle(todo.id)} className="toggle-btn">
        {todo.completed ? '未完了にする' : '完了にする'}
      </button>
      <button onClick={() => onDelete(todo.id)} className="delete-btn">
        削除
      </button>
    </div>
  </li>
);

TodoItem.displayName = 'TodoItem';

const MemoizedTodoItem = memo(TodoItem);

const TodoList: React.FC = () => {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTodo, setNewTodo] = useState<string>('');

  useEffect(() => {
    fetchTodos();
  }, []);

  const fetchTodos = async (): Promise<void> => {
    try {
      const response = await axios.get<Todo[]>(`${apiUrl}/todos`);
      setTodos(response.data);
    } catch (error) {
      console.error('Error fetching todos:', error);
    }
  };

  const addTodo = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    if (!newTodo.trim()) return;
    try {
      const response = await axios.post<Todo>(`${apiUrl}/todos`, { body: newTodo });
      setTodos(prevTodos => [...prevTodos, response.data]);
      setNewTodo('');
    } catch (error) {
      console.error('Error adding todo:', error);
    }
  };

  const toggleTodo = async (id: number): Promise<void> => {
    try {
      const response = await axios.patch<Todo>(`${apiUrl}/todos/${id}`);
      setTodos(prevTodos => prevTodos.map(todo =>
        todo.id === id ? response.data : todo
      ));
    } catch (error) {
      console.error('Error toggling todo:', error);
    }
  };

  const deleteTodo = async (id: number): Promise<void> => {
    try {
      await axios.delete(`${apiUrl}/todos/${id}`);
      setTodos(prevTodos => prevTodos.filter(todo => todo.id !== id));
    } catch (error) {
      console.error('Error deleting todo:', error);
    }
  };

  return (
    <div className="todo-container">
      <h1>Todo リスト</h1>
      <form onSubmit={addTodo} className="todo-form">
        <input
          type="text"
          value={newTodo}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setNewTodo(e.target.value)}
          placeholder="新しいタスクを入力..."
          className="todo-input"
        />
        <button type="submit" className="add-btn">追加</button>
      </form>
      <ul className="todo-list">
        {todos.map((todo) => (
          <MemoizedTodoItem
            key={todo.id}
            todo={todo}
            onToggle={toggleTodo}
            onDelete={deleteTodo}
          />
        ))}
      </ul>
    </div>
  );
};

TodoList.displayName = 'TodoList';

const App: React.FC = () => {
  return (
    <div className="App">
      <TodoList />
    </div>
  );
};

App.displayName = 'App';

export default App;