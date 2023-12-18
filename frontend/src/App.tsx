// App.tsx
import React, { useState, useEffect } from 'react';
import './App.css';

interface Stock {
  id: string;
  symbol: string;
  openPrice: number;
  currentPrice: number;
  refreshInterval: number;
}

const App: React.FC = () => {
  const [stocks, setStocks] = useState<Stock[]>([]);

  useEffect(() => {
    const fetchStocks = async () => {
      try {
        const response = await fetch('/api/stocks');
        if (response.ok) {
          const data = await response.json();
          setStocks(data);
        } else {
          console.error('Failed to fetch stocks:', response.statusText);
        }
      } catch (error) {
        console.error('Error fetching stocks:', error);
      }
    };

    fetchStocks();
  }, []);

  const updateStockPrice = (updatedStock: Stock) => {
    setStocks((prevStocks) =>
      prevStocks.map((stock) =>
        stock.id === updatedStock.id ? updatedStock : stock
      )
    );
  };

  useEffect(() => {
    const socket = new WebSocket('ws://localhost:8080/api/ws');

    socket.onmessage = (event) => {
      const updatedStock: Stock = JSON.parse(event.data);
      updateStockPrice(updatedStock);
    };

    return () => {
      socket.close();
    };
  }, []);

  return (
    <div className="App">
      <h1>Stocks</h1>
      <table>
        <thead>
          <tr>
            <th>Symbol</th>
            <th>Open Price</th>
            <th>Current Price</th>
            <th>Refresh Interval</th>
          </tr>
        </thead>
        <tbody>
          {stocks.map((stock) => (
            <tr key={stock.id}>
              <td>{stock.symbol}</td>
              <td>{stock.openPrice.toFixed(2)}</td>
              <td>{stock.currentPrice.toFixed(2)}</td>
              <td>{stock.refreshInterval} seconds</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default App;
