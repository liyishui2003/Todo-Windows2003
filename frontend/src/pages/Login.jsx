import { useState } from 'preact/hooks';
import axios from 'axios';
import LoginForm from '../components/LoginForm';

export default function Login() {
  const [error, setError] = useState('');

  const handleLogin = async (username, password) => {
    try {
      const response = await axios.post('api/login', {
        username,
        password,
      });
      const token = response.data.token;
      localStorage.setItem('jwt_token', token); // 存储 JWT 令牌
      localStorage.setItem('username', username); // 存储用户名
      
      window.location.href = '/home'; // 跳转到首页
    } catch (err) {
      setError('登录失败，请检查用户名和密码');
      console.error(err);
    }
  };

  return (
    <div className="login-container">
      <LoginForm onSubmit={handleLogin} error={error} />
    </div>
  );
}