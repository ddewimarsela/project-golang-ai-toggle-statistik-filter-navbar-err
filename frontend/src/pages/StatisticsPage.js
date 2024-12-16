import React, { useState } from 'react';
import axios from 'axios';
import { StatisticsChart } from '../components/StatisticsChart/StatisticsChart';
import './StatisticsPage.css';

const StatisticsPage = () => {
  const [statistics, setStatistics] = useState(null);
  const [loading, setLoading] = useState(false);

  const fetchStatistics = async () => {
    setLoading(true);
    try {
      const res = await axios.get("http://localhost:8080/statistics", { 
        withCredentials: true 
      });
      setStatistics(res.data.data);
    } catch (error) {
      console.error("Error fetching statistics:", error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="statistics-page">
      <div className="card statistics-card">
        <h2 className="card-title">Analisis Statistik Konsumsi Energi</h2>
        
        <button 
          onClick={fetchStatistics} 
          className="button primary-button fetch-button"
          disabled={loading}
        >
          <span className="button-icon">📊</span>
          {loading ? 'Memuat Data...' : 'Lihat Statistik'}
        </button>

        {statistics && (
          <div className="statistics-content">
            <div className="stats-grid">
              <div className="stat-box total">
                <h3>Total Konsumsi</h3>
                <div className="stat-value">{statistics.total_konsumsi.toFixed(2)}</div>
                <div className="stat-unit">kWh</div>
              </div>
              
              <div className="stat-box average">
                <h3>Rata-rata</h3>
                <div className="stat-value">{statistics.rata_rata.toFixed(2)}</div>
                <div className="stat-unit">kWh/hari</div>
              </div>

              <div className="stat-box min">
                <h3>Minimum</h3>
                <div className="stat-value">{statistics.minimum.toFixed(2)}</div>
                <div className="stat-unit">kWh</div>
              </div>

              <div className="stat-box max">
                <h3>Maximum</h3>
                <div className="stat-value">{statistics.maximum.toFixed(2)}</div>
                <div className="stat-unit">kWh</div>
              </div>

              <div className="stat-box median">
                <h3>Median</h3>
                <div className="stat-value">{statistics.median.toFixed(2)}</div>
                <div className="stat-unit">kWh</div>
              </div>
            </div>

            <div className="charts-section">
              <h3 className="section-title">Visualisasi Data</h3>
              <StatisticsChart statistics={statistics} />
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default StatisticsPage;
