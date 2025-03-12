import { Router, Route } from 'preact-router';
import Home from './pages/Home';
import Login from './pages/Login';

export function App() {
  return (
    <Router>
      <Route path="/home" component={Home} />
      <Route path="/login" component={Login} />
    </Router>
  );
}