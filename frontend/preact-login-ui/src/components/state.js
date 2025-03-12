export const tasks = []; // 全局任务列表

// 自定义事件
export const notifyTasksUpdated = new Event('tasksUpdated');

// 触发事件
export const triggerTasksUpdated = () => {
  window.dispatchEvent(notifyTasksUpdated); // 触发事件
};

export const fetchTasks = async (userID) => {
  try {
    const token = localStorage.getItem('jwt_token');
    const response = await fetch(`/api/tasks/allTask/${userID}`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
    const data = await response.json();
    console.log("data===:",data);

    const { tasks: taskList } = data;

    // 检查 tasks 是否为数组
    if (Array.isArray(taskList)) {
      tasks.length = 0; // 清空数组
      tasks.push(...taskList); // 更新任务列表
    } else {
      console.error('Invalid tasks data:', taskList);
    }

  } catch (error) {
    console.error('Failed to fetch tasks:', error);
  }
};