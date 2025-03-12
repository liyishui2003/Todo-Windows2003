export default function EditProfileModal({ onClose }) {
    return (
      <div className="modal-overlay">
        <div className="modal">
          <h2>编辑个人信息</h2>
          <form>
            <label>
              用户名：
              <input type="text" name="username" />
            </label>
            <label>
              邮箱：
              <input type="email" name="email" />
            </label>
            <button type="submit">保存</button>
            <button type="button" onClick={onClose}>取消</button>
          </form>
        </div>
      </div>
    );
  }