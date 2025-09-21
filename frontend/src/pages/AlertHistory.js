import React, { useState, useEffect } from 'react';
import { alertsAPI } from '../services/api';

const AlertHistory = () => {
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchHistory();
  }, []);

  const fetchHistory = async () => {
    try {
      const response = await alertsAPI.getHistory();
      setHistory(response.data);
    } catch (error) {
      setError('Failed to fetch alert history');
      console.error('Error fetching history:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div>
      <h1>Alert History</h1>

      {error && (
        <div className="alert alert-danger">{error}</div>
      )}

      {history.length === 0 ? (
        <div className="card">
          <h3>No alert history yet</h3>
          <p>Your alert notifications will appear here once they start sending.</p>
        </div>
      ) : (
        <div>
          {history.map(item => (
            <div key={item.id} className="card">
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                <div style={{ flex: 1 }}>
                  <h3>{item.news_title}</h3>
                  <p><strong>Source:</strong> {item.news_source}</p>
                  <p><strong>Sent:</strong> {new Date(item.sent_at).toLocaleString()}</p>

                  <div style={{ marginTop: '10px' }}>
                    <a
                      href={item.news_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="btn btn-secondary"
                    >
                      Read Article
                    </a>
                  </div>
                </div>

                <div style={{ marginLeft: '20px' }}>
                  <span
                    style={{
                      padding: '5px 10px',
                      borderRadius: '4px',
                      color: 'white',
                      backgroundColor: item.success ? 'green' : 'red',
                      fontSize: '14px',
                    }}
                  >
                    {item.success ? 'Sent' : 'Failed'}
                  </span>
                  {!item.success && item.error_msg && (
                    <p style={{ color: 'red', fontSize: '12px', marginTop: '5px' }}>
                      {item.error_msg}
                    </p>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default AlertHistory;