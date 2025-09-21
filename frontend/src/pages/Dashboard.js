import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { alertsAPI } from '../services/api';
import { useAuth } from '../services/AuthContext';

const Dashboard = () => {
  const [alerts, setAlerts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { user } = useAuth();

  useEffect(() => {
    fetchAlerts();
  }, []);

  const fetchAlerts = async () => {
    try {
      const response = await alertsAPI.getAlerts();
      setAlerts(response.data);
    } catch (error) {
      setError('Failed to fetch alerts');
      console.error('Error fetching alerts:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this alert?')) {
      try {
        await alertsAPI.deleteAlert(id);
        setAlerts(alerts.filter(alert => alert.id !== id));
      } catch (error) {
        setError('Failed to delete alert');
        console.error('Error deleting alert:', error);
      }
    }
  };

  const handleToggleActive = async (alert) => {
    try {
      const response = await alertsAPI.updateAlert(alert.id, {
        active: !alert.active,
      });
      setAlerts(alerts.map(a => a.id === alert.id ? response.data : a));
    } catch (error) {
      setError('Failed to update alert');
      console.error('Error updating alert:', error);
    }
  };

  const handleTestAlert = async (alertId) => {
    try {
      await alertsAPI.testAlert(alertId);
      alert('Test alert sent successfully!');
    } catch (error) {
      setError('Failed to send test alert');
      console.error('Error sending test alert:', error);
    }
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h1>Welcome, {user?.email}</h1>
        <Link to="/alerts/new" className="btn">
          Create New Alert
        </Link>
      </div>

      {error && (
        <div className="alert alert-danger">{error}</div>
      )}

      {alerts.length === 0 ? (
        <div className="card">
          <h3>No alerts yet</h3>
          <p>Create your first news alert to get started!</p>
          <Link to="/alerts/new" className="btn">
            Create Alert
          </Link>
        </div>
      ) : (
        <div>
          <h2>Your News Alerts</h2>
          {alerts.map(alert => (
            <div key={alert.id} className="card">
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                <div style={{ flex: 1 }}>
                  <h3>{alert.topic}</h3>
                  <p><strong>Keywords:</strong> {alert.keywords.join(', ')}</p>
                  <p><strong>Frequency:</strong> {alert.frequency}</p>
                  <p><strong>Status:</strong>
                    <span style={{
                      color: alert.active ? 'green' : 'red',
                      fontWeight: 'bold',
                      marginLeft: '5px'
                    }}>
                      {alert.active ? 'Active' : 'Inactive'}
                    </span>
                  </p>
                  {alert.last_checked && (
                    <p><strong>Last Checked:</strong> {new Date(alert.last_checked).toLocaleString()}</p>
                  )}
                </div>

                <div style={{ display: 'flex', gap: '10px', flexDirection: 'column' }}>
                  <button
                    onClick={() => handleToggleActive(alert)}
                    className={`btn ${alert.active ? 'btn-secondary' : 'btn'}`}
                  >
                    {alert.active ? 'Deactivate' : 'Activate'}
                  </button>

                  <button
                    onClick={() => handleTestAlert(alert.id)}
                    className="btn btn-secondary"
                  >
                    Test
                  </button>

                  <Link
                    to={`/alerts/edit/${alert.id}`}
                    className="btn btn-secondary"
                    style={{ textDecoration: 'none', textAlign: 'center' }}
                  >
                    Edit
                  </Link>

                  <button
                    onClick={() => handleDelete(alert.id)}
                    className="btn btn-danger"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Dashboard;