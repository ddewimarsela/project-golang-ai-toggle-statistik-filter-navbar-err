import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import './Navbar.css';
import ThemeToggle from '../ThemeToggle/ThemeToggle';

const Navbar = ({ theme, toggleTheme }) => {
  const location = useLocation();

  return (
    <nav className="navbar">
      <div className="navbar-container">
        <div className="navbar-brand">
          <Link to="/" className="brand-link">
            <h1>Energy Data Analysis</h1>
          </Link>
        </div>
        <div className="navbar-menu">
          <div className="navbar-start">
            <Link 
              to="/" 
              className={`nav-item ${location.pathname === '/' ? 'active' : ''}`}
            >
              Upload Data
            </Link>
            <Link 
              to="/chat" 
              className={`nav-item ${location.pathname === '/chat' ? 'active' : ''}`}
            >
              Ask Question
            </Link>
            <Link 
              to="/statistics" 
              className={`nav-item ${location.pathname === '/statistics' ? 'active' : ''}`}
            >
              Statistics
            </Link>
            <Link 
              to="/filter" 
              className={`nav-item ${location.pathname === '/filter' ? 'active' : ''}`}
            >
              Filter Data
            </Link>
          </div>
          <div className="navbar-end">
            <ThemeToggle theme={theme} toggleTheme={toggleTheme} />
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
