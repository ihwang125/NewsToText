import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { alertsAPI } from '../services/api';

const AlertForm = () => {
  const [formData, setFormData] = useState({
    topic: '',
    keywords: '',
    frequency: 'daily',
    active: true,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [isEdit, setIsEdit] = useState(false);

  const navigate = useNavigate();
  const { id } = useParams();

  useEffect(() => {
    if (id) {
      setIsEdit(true);
      fetchAlert();
    }
  }, [id]);

  const fetchAlert = async () => {
    try {
      setLoading(true);
      const response = await alertsAPI.getAlerts();
      const alert = response.data.find(a => a.id === parseInt(id));

      if (alert) {
        setFormData({
          topic: alert.topic,
          keywords: alert.keywords.join(', '),
          frequency: alert.frequency,
          active: alert.active,
        });
      } else {
        setError('Alert not found');
      }
    } catch (error) {
      setError('Failed to fetch alert');
      console.error('Error fetching alert:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData({
      ...formData,
      [name]: type === 'checkbox' ? checked : value,
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const keywords = formData.keywords
      .split(',')
      .map(k => k.trim())
      .filter(k => k.length > 0);

    if (keywords.length === 0) {
      setError('Please enter at least one keyword');
      setLoading(false);
      return;
    }

    const alertData = {
      topic: formData.topic,
      keywords,
      frequency: formData.frequency,
      ...(isEdit && { active: formData.active }),
    };

    try {
      if (isEdit) {
        await alertsAPI.updateAlert(id, alertData);
      } else {
        await alertsAPI.createAlert(alertData);
      }
      navigate('/');
    } catch (error) {
      setError(error.response?.data?.error || `Failed to ${isEdit ? 'update' : 'create'} alert`);
      console.error(`Error ${isEdit ? 'updating' : 'creating'} alert:`, error);
    } finally {
      setLoading(false);
    }
  };

  if (loading && isEdit) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="card" style={{ maxWidth: '600px', margin: '0 auto' }}>
      <h2>{isEdit ? 'Edit Alert' : 'Create New Alert'}</h2>

      {error && (
        <div className="alert alert-danger">{error}</div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="topic">Topic</label>
          <input
            type="text"
            id="topic"
            name="topic"
            value={formData.topic}
            onChange={handleChange}
            placeholder="e.g., Technology, Stocks, Sports"
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="keywords">Keywords (comma-separated)</label>
          <textarea
            id="keywords"
            name="keywords"
            value={formData.keywords}
            onChange={handleChange}
            placeholder="e.g., artificial intelligence, machine learning, AI"
            rows="3"
            required
          />
          <small style={{ color: '#666' }}>
            Enter keywords separated by commas. Articles matching any of these keywords will trigger alerts.
          </small>
        </div>

        <div className="form-group">
          <label htmlFor="frequency">Frequency</label>
          <select
            id="frequency"
            name="frequency"
            value={formData.frequency}
            onChange={handleChange}
            required
          >
            <option value="realtime">Real-time (every 5 minutes)</option>
            <option value="hourly">Hourly</option>
            <option value="daily">Daily</option>
          </select>
        </div>

        {isEdit && (
          <div className="form-group">
            <label>
              <input
                type="checkbox"
                name="active"
                checked={formData.active}
                onChange={handleChange}
                style={{ marginRight: '8px' }}
              />
              Active
            </label>
          </div>
        )}

        <div style={{ display: 'flex', gap: '10px' }}>
          <button
            type="submit"
            className="btn"
            disabled={loading}
          >
            {loading ? (isEdit ? 'Updating...' : 'Creating...') : (isEdit ? 'Update Alert' : 'Create Alert')}
          </button>

          <button
            type="button"
            className="btn btn-secondary"
            onClick={() => navigate('/')}
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
};

export default AlertForm;