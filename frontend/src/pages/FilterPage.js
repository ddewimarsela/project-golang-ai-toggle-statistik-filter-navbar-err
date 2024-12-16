import React, { useState } from 'react';
import axios from 'axios';
import './FilterPage.css';

const FilterPage = () => {
  const [filters, setFilters] = useState({
    room: '',
    appliance: '',
    date: ''
  });
  const [filteredData, setFilteredData] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleFilter = async () => {
    setLoading(true);
    try {
      const queryParams = new URLSearchParams(filters).toString();
      const response = await axios.get(`http://localhost:8080/filter?${queryParams}`, {
        withCredentials: true
      });
      
      if (response.data.status === 'success') {
        setFilteredData(response.data.data);
      }
    } catch (error) {
      console.error('Error filtering data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    setFilters({
      room: '',
      appliance: '',
      date: ''
    });
    setFilteredData(null);
  };

  return (
    <div className="filter-page">
      <div className="card filter-card">
        <h2 className="card-title">Filter Data Energi</h2>
        <p className="filter-description">
          Filter dan analisis data konsumsi energi berdasarkan ruangan, perangkat, atau tanggal
        </p>

        <div className="filters-container">
          <div className="filter-group">
            <label className="filter-label">
              <span className="label-icon">🏠</span>
              Ruangan
            </label>
            <input
              type="text"
              placeholder="Contoh: Living Room, Kitchen"
              value={filters.room}
              onChange={(e) => setFilters({...filters, room: e.target.value})}
              className="filter-input"
            />
          </div>

          <div className="filter-group">
            <label className="filter-label">
              <span className="label-icon">🔌</span>
              Perangkat
            </label>
            <input
              type="text"
              placeholder="Contoh: AC, TV, Lamp"
              value={filters.appliance}
              onChange={(e) => setFilters({...filters, appliance: e.target.value})}
              className="filter-input"
            />
          </div>

          <div className="filter-group">
            <label className="filter-label">
              <span className="label-icon">📅</span>
              Tanggal
            </label>
            <input
              type="date"
              value={filters.date}
              onChange={(e) => setFilters({...filters, date: e.target.value})}
              className="filter-input"
            />
          </div>

          <div className="filter-actions">
            <button 
              onClick={handleFilter} 
              className={`button primary-button ${loading ? 'loading' : ''}`}
              disabled={loading}
            >
              {loading ? (
                <>
                  <span className="spinner"></span>
                  <span>Filtering...</span>
                </>
              ) : (
                <>
                  <span className="button-icon">🔍</span>
                  Filter Data
                </>
              )}
            </button>
            <button 
              onClick={handleReset} 
              className="button secondary-button"
              disabled={loading}
            >
              <span className="button-icon">🔄</span>
              Reset
            </button>
          </div>
        </div>

        {filteredData && (
          <div className="results-container">
            <h3 className="results-title">
              <span className="results-icon">📊</span>
              Hasil Filter
            </h3>
            <div className="table-container">
              <table className="results-table">
                <thead>
                  <tr>
                    {Object.keys(filteredData).map(key => (
                      <th key={key}>{key}</th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {filteredData[Object.keys(filteredData)[0]].map((_, rowIndex) => (
                    <tr key={rowIndex}>
                      {Object.keys(filteredData).map(key => (
                        <td key={key}>{filteredData[key][rowIndex]}</td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default FilterPage;
