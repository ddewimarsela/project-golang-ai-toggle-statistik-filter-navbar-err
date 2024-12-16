import React, { useState } from 'react';
import axios from 'axios';
import './UploadPage.css';

const UploadPage = () => {
  const [file, setFile] = useState(null);
  const [response, setResponse] = useState("");
  const [loading, setLoading] = useState(false);
  const [dragActive, setDragActive] = useState(false);
  const [fileName, setFileName] = useState("");

  const handleDrag = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  };

  const handleDrop = (e) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      const file = e.dataTransfer.files[0];
      if (file.type === "text/csv") {
        setFile(file);
        setFileName(file.name);
      } else {
        setResponse("Please upload a CSV file");
      }
    }
  };

  const handleFileChange = (e) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      if (file.type === "text/csv") {
        setFile(file);
        setFileName(file.name);
        setResponse("");
      } else {
        setResponse("Please upload a CSV file");
      }
    }
  };

  const handleUpload = async () => {
    if (!file) {
      setResponse("Please select a file first");
      return;
    }

    setLoading(true);
    const formData = new FormData();
    formData.append("file", file);

    try {
      const res = await axios.post('http://localhost:8080/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        withCredentials: true
      });
      setResponse(res.data.answer);
    } catch (error) {
      console.error('Error uploading file:', error);
      setResponse("Error uploading file");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="upload-page">
      <div className="card upload-card">
        <h2 className="card-title">Upload Data Energi</h2>
        

        <div 
          className={`upload-zone ${dragActive ? 'drag-active' : ''}`}
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={handleDrop}
        >
          <input
            type="file"
            onChange={handleFileChange}
            className="file-input"
            accept=".csv"
            id="file-input"
          />
          <label htmlFor="file-input" className="file-label">
            <div className="upload-icon">📁</div>
            <div className="upload-text">
              {fileName ? (
                <span className="file-name">{fileName}</span>
              ) : (
                <>
                  <span className="drag-text">Drag and drop your CSV file here</span>
                  <span className="or-text">or</span>
                  <span className="browse-text">Browse Files</span>
                </>
              )}
            </div>
          </label>
        </div>

        <button 
          onClick={handleUpload} 
          className={`button primary-button upload-button ${loading ? 'loading' : ''}`}
          disabled={!file || loading}
        >
          {loading ? (
            <>
              <span className="spinner"></span>
              <span>Uploading...</span>
            </>
          ) : (
            <>
              <span className="button-icon">📊</span>
              <span>Analyze Data</span>
            </>
          )}
        </button>

        {response && (
          <div className={`response-section ${response.includes('Error') ? 'error' : 'success'}`}>
            <h3 className="response-title">
              {response.includes('Error') ? '❌ Error' : '✅ Analysis Results'}
            </h3>
            <div className="response-content">
              {response}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default UploadPage;
