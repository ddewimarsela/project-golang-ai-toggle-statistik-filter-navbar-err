import React, { useState } from 'react';
import axios from 'axios';
import './ChatPage.css';

const ChatPage = () => {
  const [query, setQuery] = useState("");
  const [loading, setLoading] = useState(false);
  const [chatHistory, setChatHistory] = useState([]);

  const handleChat = async () => {
    if (!query.trim()) return;

    setLoading(true);
    try {
      const res = await axios.post("http://localhost:8080/chat", 
        { query }, 
        { withCredentials: true }
      );
      
      const newMessage = {
        type: 'user',
        content: query,
        timestamp: new Date().toLocaleTimeString()
      };
      
      const newResponse = {
        type: 'bot',
        content: res.data.answer,
        timestamp: new Date().toLocaleTimeString()
      };

      setChatHistory(prev => [...prev, newMessage, newResponse]);
      setQuery("");
    } catch (error) {
      console.error("Error querying chat:", error);
      setChatHistory(prev => [...prev, {
        type: 'bot',
        content: "Error occurred while processing your question",
        timestamp: new Date().toLocaleTimeString()
      }]);
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleChat();
    }
  };

  return (
    <div className="chat-page">
      <div className="card chat-card">
        <h2 className="card-title">Ask Questions About Your Energy Data</h2>
        <p className="chat-description">
          Tanyakan tentang konsumsi energi Anda dan dapatkan insight yang bermanfaat
        </p>

        <div className="chat-container">
          <div className="chat-history">
            {chatHistory.length === 0 ? (
              <div className="chat-empty-state">
                <span className="empty-icon">💡</span>
                <p>Mulai bertanya tentang data energi Anda</p>
                <div className="suggestion-chips">
                  <button 
                    className="suggestion-chip"
                    onClick={() => setQuery("Berapa total konsumsi energi?")}
                  >
                    Berapa total konsumsi energi?
                  </button>
                  <button 
                    className="suggestion-chip"
                    onClick={() => setQuery("Kapan penggunaan energi tertinggi?")}
                  >
                    Kapan penggunaan energi tertinggi?
                  </button>
                  <button 
                    className="suggestion-chip"
                    onClick={() => setQuery("Bagaimana tren penggunaan energi?")}
                  >
                    Bagaimana tren penggunaan energi?
                  </button>
                </div>
              </div>
            ) : (
              chatHistory.map((message, index) => (
                <div 
                  key={index} 
                  className={`chat-message ${message.type}-message`}
                >
                  <div className="message-content">
                    <span className="message-icon">
                      {message.type === 'user' ? '👤' : '🤖'}
                    </span>
                    <div className="message-text">
                      {message.content}
                    </div>
                  </div>
                  <span className="message-timestamp">{message.timestamp}</span>
                </div>
              ))
            )}
          </div>

          <div className="chat-input-container">
            <textarea
              className="chat-input"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="Ketik pertanyaan Anda di sini..."
              rows="1"
            />
            <button 
              onClick={handleChat} 
              className={`button primary-button send-button ${loading ? 'loading' : ''}`}
              disabled={!query.trim() || loading}
            >
              {loading ? (
                <span className="spinner"></span>
              ) : (
                <span className="button-icon">Ask</span>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ChatPage;
