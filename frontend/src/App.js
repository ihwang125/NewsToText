import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './services/AuthContext';
import Navbar from './components/Navbar';
import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import AlertForm from './pages/AlertForm';
import AlertHistory from './pages/AlertHistory';
import ProtectedRoute from './components/ProtectedRoute';

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="App">
          <Navbar />
          <div className="container">
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route
                path="/"
                element={
                  <ProtectedRoute>
                    <Dashboard />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/alerts/new"
                element={
                  <ProtectedRoute>
                    <AlertForm />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/alerts/edit/:id"
                element={
                  <ProtectedRoute>
                    <AlertForm />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/alerts/history"
                element={
                  <ProtectedRoute>
                    <AlertHistory />
                  </ProtectedRoute>
                }
              />
            </Routes>
          </div>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;