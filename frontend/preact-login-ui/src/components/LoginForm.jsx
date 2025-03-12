import { useState } from 'preact/hooks';

export default function LoginForm({ onSubmit, error }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    onSubmit(username, password);
  };

  return (
    <div className="login-box">
      <h1>登录</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label for="username">用户名</label>
          <input
            type="text"
            id="username"
            value={username}
            onInput={(e) => setUsername(e.target.value)}
            placeholder="请输入用户名"
            required
          />
        </div>
        <div className="form-group">
          <label for="password">密码</label>
          <input
            type="password"
            id="password"
            value={password}
            onInput={(e) => setPassword(e.target.value)}
            placeholder="请输入密码"
            required
          />
        </div>
        {error && <p className="error-message">{error}</p>}
        <button type="submit">登录</button>
      </form>
    </div>
  );
}