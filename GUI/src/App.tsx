import logo from "./logo.svg"
import { Counter } from "./features/counter/Counter"
import "./App.css"

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <Counter />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <span>
          <span>Learn </span>
          <a
            className="App-link"
            href="https://reactjs.org/"
            target="_blank"
            rel="noopener noreferrer"
          >
            React
          </a>
          <span>, </span>
          <a
            className="App-link"
            href="https://redux.js.org/"
            target="_blank"
            rel="noopener noreferrer"
          >
            Redux
          </a>
          <span>, </span>
          <a
            className="App-link"
            href="https://redux-toolkit.js.org/"
            target="_blank"
            rel="noopener noreferrer"
          >
            Redux Toolkit
          </a>
          ,<span> and </span>
          <a
            className="App-link"
            href="https://react-redux.js.org/"
            target="_blank"
            rel="noopener noreferrer"
          >
            React Redux
          </a>
        </span>
      </header>
    </div>
  )
}

// Install dependencies

// Import dependencies
import React, { useState, useEffect } from 'react';
import { Line } from 'react-chartjs-2';
import io from 'socket.io-client';

// Create component
const ChartPage = () => {
  const [data, setData] = useState([]);

  useEffect(() => {
    // Connect to WebSocket
    const socket = io('http://localhost:3000');

    // Listen for data
    socket.on('data', (newData) => {
      setData((prevData) => [...prevData, newData]);
    });

    // Disconnect from WebSocket on unmount
    return () => {
      socket.disconnect();
    };
  }, []);

  // Format data for chart
  const chartData = {
    labels: data.map((d) => d.timestamp),
    datasets: [
      {
        label: 'Data',
        data: data.map((d) => d.value),
        fill: false,
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.1,
      },
    ],
  };

  return (
    <div>
      <h1>Real-time Chart</h1>
      <Line data={chartData} />
    </div>
  );
};

// export default ChartPage;

export default App
