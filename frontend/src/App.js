import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar/Navbar';
import UploadPage from './pages/UploadPage';
import ChatPage from './pages/ChatPage';
import StatisticsPage from './pages/StatisticsPage';
import FilterPage from './pages/FilterPage';
import { useTheme } from './hooks/useTheme';
import './App.css';

function App() {
  const [theme, setTheme] = useTheme();

  const toggleTheme = () => {
    setTheme(theme === 'light' ? 'dark' : 'light');
  };

  React.useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
  }, [theme]);

  return (
    <Router>
      <div className={`app ${theme}`}>
        <Navbar theme={theme} toggleTheme={toggleTheme} />
        <div className="container">
          <header className="header">
            <p className="app-subtitle">Upload your data and get smart energy insights</p>
          </header>

          <main className="main-content">
            <Routes>
              <Route path="/" element={<UploadPage />} />
              <Route path="/chat" element={<ChatPage />} />
              <Route path="/statistics" element={<StatisticsPage />} />
              <Route path="/filter" element={<FilterPage />} />
            </Routes>
          </main>

        </div>
      </div>
    </Router>
  );
}

export default App;