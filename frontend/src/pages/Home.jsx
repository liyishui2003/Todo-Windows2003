import Sidebar from '../components/Sidebar';
import TaskPool from '../components/TaskPool';
import axios from 'axios';

import { useState, useEffect } from 'preact/hooks';

export default function Home() {
  const [error, setError] = useState(null); // 定义 error 和 setError
  const username = localStorage.getItem('username');
  const token = localStorage.getItem('jwt_token');

  const GetUserID = async () => {
    try {
      const userIDResponse = await axios.get(`api/users/${username}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      const userID = userIDResponse.data.userID;
      localStorage.setItem('userID', userID);
      console.log('User ID:', userID);
    } catch (err) {
      setError('请求失败，未能正确得到ID'); // 使用 setError 设置错误信息
      console.error('Failed to fetch userID:', err);
    }
  };

  useEffect(() => {
    GetUserID();
  }, []);

  return (
    <div className="home-container">
      {error && <div className="error-message">{error}</div>} {/* 显示错误信息 */}
      <Sidebar />
      <TaskPool />
    </div>
  );
}