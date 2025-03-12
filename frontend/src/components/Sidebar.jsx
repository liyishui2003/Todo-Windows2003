import { useState, useEffect } from 'preact/hooks'; // 导入 useEffect
import axios from 'axios'; // 导入 axios
import EditProfileModal from './EditProfileModal';
import CreateTodoModal from './CreateTodoModal'; // 导入新建Todo的弹窗组件

export default function Sidebar( ) {
  const [isProfileModalOpen, setIsProfileModalOpen] = useState(false); // 控制编辑个人信息弹窗
  const [isTodoModalOpen, setIsTodoModalOpen] = useState(false); // 控制新建Todo弹窗

  // 获取 JWT 令牌
  const token = localStorage.getItem('jwt_token');

  // 配置 axios 请求头
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;

  // 打开编辑个人信息弹窗
  const openProfileModal = () => setIsProfileModalOpen(true);
  const closeProfileModal = () => setIsProfileModalOpen(false);

  // 打开新建Todo弹窗
  const openTodoModal = () => setIsTodoModalOpen(true);
  const closeTodoModal = () => setIsTodoModalOpen(false);

  const userName = localStorage.getItem('username');
  return (
    <div className="sidebar">
      <h1>{userName}</h1>
      {/* 编辑个人信息按钮 */}
      <button onClick={openProfileModal}>编辑个人信息</button>
      {/* 新建Todo按钮 */}
      <button onClick={openTodoModal}>新建Todo</button>

      {/* 编辑个人信息弹窗 */}
      {isProfileModalOpen && <EditProfileModal onClose={closeProfileModal} />}

      {/* 新建Todo弹窗 */}
      {isTodoModalOpen && (
        <CreateTodoModal
          onClose={closeTodoModal}
        />
      )}
    </div>
  );
}