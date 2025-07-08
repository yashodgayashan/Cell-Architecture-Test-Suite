import React, { useState } from 'react';
import './App.css';

function App() {
  const [apiResponse, setApiResponse] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);

  // The URL for the frontend-svc. This should be configured via an environment
  // variable or an external config file in a real deployment.
  const API_URL_CELL_A = process.env.REACT_APP_API_URL_CELL_A || '/api/cell-a';
  const API_URL_CELL_B = process.env.REACT_APP_API_URL_CELL_B || '/api/cell-b';


  const callApi = async (apiUrl) => {
    setIsLoading(true);
    setError(null);
    setApiResponse(null);

    try {
      const response = await fetch(apiUrl);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setApiResponse(data);
    } catch (e) {
      setError(e.message);
      console.error("API call failed:", e);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Scale-to-Zero Test Client</h1>
        <p>Click a button to send a request to the corresponding cell's frontend service.</p>
        <div className="button-container">
          <button onClick={() => callApi(API_URL_CELL_A)} disabled={isLoading}>
            {isLoading ? 'Loading...' : 'Call Cell A'}
          </button>
          <button onClick={() => callApi(API_URL_CELL_B)} disabled={isLoading}>
            {isLoading ? 'Loading...' : 'Call Cell B'}
          </button>
        </div>

        {apiResponse && (
          <div className="response-container">
            <h3>API Response:</h3>
            <pre>{JSON.stringify(apiResponse, null, 2)}</pre>
          </div>
        )}

        {error && (
          <div className="error-container">
            <h3>Error:</h3>
            <p>{error}</p>
          </div>
        )}
      </header>
    </div>
  );
}

export default App;