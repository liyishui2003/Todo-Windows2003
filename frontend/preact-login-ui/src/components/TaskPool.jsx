import { useEffect, useState,useCallback } from 'preact/hooks';
import axios from 'axios';
import { tasks, fetchTasks, triggerTasksUpdated } from './state';

export default function TaskPool({ }) {
  
  const [currentTime, setCurrentTime] = useState(new Date().toLocaleTimeString());
  const [location, setLocation] = useState("");
  const [ipAddress, setIpAddress] = useState("");
  const [quadrants, setQuadrants] = useState({
    quadrant1: [],
    quadrant2: [],
    quadrant3: [],
    quadrant4: [],
  });
  const userID = localStorage.getItem("userID");
  // 获取 JWT 令牌
  const token = localStorage.getItem('jwt_token');
  // 配置 axios 请求头
  axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;

  // 更新当前时间
  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(new Date().toLocaleTimeString());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  // 获取当前所在地和 IP 地址
  useEffect(() => {
    const fetchLocationAndIp = async () => {
      try {
        // 使用 ipwhois.io 获取 IP 地址和地理位置信息
        const response = await axios.get('https://ipwhois.app/json/');
        if (response.data.success) {
          setIpAddress(response.data.ip); // 设置 IP 地址
          setLocation(`${response.data.city}, ${response.data.country}`); // 设置地理位置
        } else {
          console.error('Failed to fetch location:', response.data);
        }
      } catch (error) {
        console.error('Failed to fetch location or IP:', error);
      }
    };
  
    fetchLocationAndIp();
  }, []);

  const handleDelete = useCallback(async (taskID) => {
    console.log('Deleting task:', taskID); // 调试日志
    try {
      await axios.delete(`/api/tasks/${taskID}`);
      await fetchTasks(userID); // 等待 fetchTasks 完成
      triggerTasksUpdated();
    } catch (error) {
      console.error('Failed to delete task:', error);
    }
  }, []); 

  const categorizeTasks = (tasks) => {
    if (!Array.isArray(tasks)) {
      console.error('Tasks is not an array:', tasks);
      return {
        quadrant1: [],
        quadrant2: [],
        quadrant3: [],
        quadrant4: [],
      };
    }
  
    const quadrants = {
      quadrant1: [], // Important & Not Urgent
      quadrant2: [], // Important & Urgent
      quadrant3: [], // Not Important & Urgent
      quadrant4: [], // Not Important & Not Urgent
    };
  
    tasks.forEach(task => {
      if (task.importance === 1 && task.urgency === 0) {
        quadrants.quadrant1.push(task);
      } else if (task.importance === 1 && task.urgency === 1) {
        quadrants.quadrant2.push(task);
      } else if (task.importance === 0 && task.urgency === 1) {
        quadrants.quadrant3.push(task);
      } else {
        quadrants.quadrant4.push(task);
      }
    });
  
    return quadrants;
  };

  // 初始化四个象限
  useEffect(() => {
    const fetchData = async () => {
      await fetchTasks(userID); // 等待 fetchTasks 完成
      setQuadrants(categorizeTasks(tasks)); // 更新任务象限
      console.log("Quadrants:", quadrants); // 打印 quadrants
    };
    fetchData();
  }, [userID]);

  // 通过自定义事件实时监测
  useEffect(() => {
    const handleTasksUpdated = () => {
      setQuadrants(categorizeTasks(tasks)); // 重新计算象限
    };

    window.addEventListener('tasksUpdated', handleTasksUpdated);
    return () => window.removeEventListener('tasksUpdated', handleTasksUpdated);
  }, []);

  return (
    <div className="window">
      {/* 顶边栏 */}
      <div className="title-bar">
        <div className="title" style={{ fontSize: '24px', fontWeight: 'bold', fontFamily: 'Arial, sans-serif' }}>Task Pool</div>
        <div className="clock-info">
          <div className="clock">{currentTime}</div>
          <div className="ip-address">IP: {ipAddress}</div>
          <div className="location">{location}</div>
        </div>
      </div>

      {/* 内容区域 */}
      <div className="content">
        {/* 任务池 */}
        <div className="quadrant-container">
          {/* 第一象限：重要不紧急 */}
          <div className="quadrant quadrant1">
          <QuadrantTitle importance={1} urgency={0} />
            <div className="task-list">
              {quadrants.quadrant1.map(task => (
                <TaskCard key={task.id} task={task} handleDelete={handleDelete} />
              ))}
            </div>
          </div>

          {/* 第二象限：重要且紧急 */}
          <div className="quadrant quadrant2">
          <QuadrantTitle importance={1} urgency={1} />
            <div className="task-list">
              {quadrants.quadrant2.map(task => (
                <TaskCard key={task.id} task={task} handleDelete={handleDelete} />
              ))}
            </div>
          </div>

          {/* 第三象限：紧急不重要 */}
          <div className="quadrant quadrant3">
          <QuadrantTitle importance={0} urgency={1} />
            <div className="task-list">
              {quadrants.quadrant3.map(task => (
                <TaskCard key={task.id} task={task} handleDelete={handleDelete} />
              ))}
            </div>
          </div>

          {/* 第四象限：不重要不紧急 */}
          <div className="quadrant quadrant4">
          <QuadrantTitle importance={0} urgency={0} />
            <div className="task-list">
              {quadrants.quadrant4.map(task => (
                <TaskCard key={task.id} task={task} handleDelete={handleDelete} />
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// 任务卡片组件
function TaskCard({ task, handleDelete }) {
  return (
    <div className="task-card" style={{ padding: '8px' }}>
      <div className="task-content">
        {/* 标题和删除按钮 */}
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h3>{task.title}</h3>
          {/* 删除按钮 */}
          <button
            onClick={() => handleDelete(task.id)}
            style={{
              width: '19px', // 正方形宽度
              height: '19px', // 正方形高度
              backgroundColor: 'white', // 背景留白
              border: '0.5px solid black', // 黑色边框
              color: 'black', // 黑色打叉
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              cursor: 'pointer',
              boxShadow: '0.5px 0.5px 1px rgba(0, 0, 0, 0.2)', // 阴影效果
              fontSize: '10px', // 控制图标大小
              lineHeight: '1', // 确保图标居中
              padding: 0, // 去除默认内边距
              marginTop: '-11px',
              marginRight: '-4px'
            }}
          >
            ✕
          </button>
        </div>
        <p>{task.description}</p>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <p className="deadline">Deadline: {task.due_date}</p>
          {/* 完成按钮 */}
          <button
           // onClick={() => handleComplete(task.id)}
            style={{
              width: '18px', // 圆形宽度
              height: '18px', // 圆形高度
              borderRadius: '50%', // 圆形
              backgroundColor: 'white', // 背景留白
              border: '0.5px solid black', // 黑色边框
              color: 'black', // 黑色打钩
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              cursor: 'pointer',
              boxShadow: '1px 1px 2px rgba(0, 0, 0, 0.2)', // 阴影效果
              fontSize: '12px', // 控制图标大小
              lineHeight: '1', // 确保图标居中
              padding: 0, // 去除默认内边距
              marginTop: '-7px',
              marginRight: '-3px'
            }}
          >
            ✓
          </button>
        </div>
      </div>
    </div>
  );
}

// 任务象限标题组件
function QuadrantTitle({ importance, urgency }) {
  return (
    <div style={{ display: 'flex', gap: '8px', alignItems: 'center', marginBottom: '16px' }}>
      {/* 重要性方块 */}
      <div
        style={{
          width: '40px', // 长方形宽度
          height: '20px', // 长方形高度
          backgroundColor: importance === 1 ? 'red' : 'green',
          color: 'white',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '12px',
          fontFamily: 'Arial, sans-serif',
          border: '1px solid #ccc',
          boxShadow: '2px 2px 4px rgba(0, 0, 0, 0.2)', // 阴影效果
          borderRadius: '4px', // 圆角
        }}
      >
        {importance === 1 ? '重要' : '不重要'}
      </div>
      {/* 紧急性方块 */}
      <div
        style={{
          width: '40px', // 长方形宽度
          height: '20px', // 长方形高度
          backgroundColor: urgency === 1 ? 'red' : 'green',
          color: 'white',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '12px',
          fontFamily: 'Arial, sans-serif',
          border: '1px solid #ccc',
          boxShadow: '2px 2px 4px rgba(0, 0, 0, 0.2)', // 阴影效果
          borderRadius: '4px', // 圆角
        }}
      >
        {urgency === 1 ? '紧急' : '不紧急'}
      </div>
    </div>
  );
}