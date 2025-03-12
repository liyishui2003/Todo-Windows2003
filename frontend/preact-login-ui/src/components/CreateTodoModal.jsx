import { useState, useEffect } from 'preact/hooks';
import axios from 'axios';
import { fetchTasks,triggerTasksUpdated } from './state';

export default function CreateTodoModal({onClose}) {
  // userID用来调用state.js里的fetchtask()

  const userID = localStorage.getItem('userID');
  // 获取 JWT 令牌
  const token = localStorage.getItem('jwt_token');

  // 配置 axios 请求头
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;

  const [newTodo, setNewTodo] = useState({
    title: '',
    description: '',
    due_date: '',
    importance: 0,
    urgency: 0,
  });

  // 处理输入字段变化
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setNewTodo({ ...newTodo, [name]: value });
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0'); // 补零
    const day = String(date.getDate()).padStart(2, '0'); // 补零
    return `${year}-${month}-${day}`;
  };

  // 提交新 Todo
  const handleSubmit = async () => {
    try {
      // 格式化日期
      const formattedDate = formatDate(newTodo.due_date);
      const todoData = {
        title: newTodo.title,
        description: newTodo.description,
        due_date: formattedDate,
        importance: parseInt(newTodo.importance, 10), // 转换为整数
        urgency: parseInt(newTodo.urgency, 10), // 转换为整数
      };

      // 提交数据
      const response = await axios.post('/api/tasks', todoData);
      console.log('Todo created:', response.data);
      await fetchTasks(userID); // 等待 fetchTasks 完成
      triggerTasksUpdated();
      onClose(); // 关闭弹窗
    } catch (error) {
      console.error('Failed to create todo:', error);
    }
  };

  return (
    <div
      style={{
        position: 'fixed',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        backgroundColor: '#c0c0c0',
        border: '2px solid black',
        padding: '16px',
        boxShadow: '3px 3px 0 rgba(0, 0, 0, 0.2)',
        zIndex: 1000,
        width: '350px', // 弹窗宽度
        borderRadius: '4px', // 圆角
      }}
    >
      {/* 弹窗标题栏 */}
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          backgroundColor: '#000080',
          color: 'white',
          padding: '4px',
          borderRadius: '4px 4px 0 0', // 圆角
        }}
      >
        <div>新建Todo</div>
        <button
          onClick={onClose}
          style={{
            backgroundColor: 'transparent',
            width:"20px",
            height:"20px",
            border: 'none',
            color: 'white',
            cursor: 'pointer',
            marginTop: '-10px',
            marginLeft: '100px'
          }}
        >
          ✕
        </button>
      </div>

      {/* 弹窗内容 */}
      <div style={{ marginTop: '16px' }}>
        <label>
          名称:
          <input
            type="text"
            name="title"
            value={newTodo.title}
            onChange={handleInputChange}
            style={{ marginLeft: '8px', width: '250px' }} // 缩小输入框长度
          />
        </label>
        <br />
        <label>
          描述:
          <input
            type="text"
            name="description"
            value={newTodo.description}
            onChange={handleInputChange}
            style={{ marginLeft: '8px', width: '250px' }} // 缩小输入框长度
          />
        </label>
        <br />
        <label>
          截止:
          <input
            type="date"
            name="due_date"
            value={newTodo.due_date}
            onChange={handleInputChange}
            style={{ marginLeft: '8px', width: '250px' }} // 缩小输入框长度
          />
        </label>
        <br />
        <label>
          Importance:
          <select
            name="importance"
            value={newTodo.importance}
            onChange={handleInputChange}
            style={{ marginLeft: '8px', width: '120px' }} // 缩小输入框长度
          >
            <option value={0}>Not Important</option>
            <option value={1}>Important</option>
          </select>
        </label>
        <br />
        <label>
          Urgency:
          <select
            name="urgency"
            value={newTodo.urgency}
            onChange={handleInputChange}
            style={{ marginLeft: '8px', width: '120px' }} // 缩小输入框长度
          >
            <option value={0}>Not Urgent</option>
            <option value={1}>Urgent</option>
          </select>
        </label>
        <br />
        <button onClick={handleSubmit} style={{ marginTop: '16px' }}>
          提交
        </button>
      </div>
    </div>
  );
}