import React from 'react';
import { Routes } from 'react-router'
import { BrowserRouter, Route } from 'react-router-dom';
import Home from './pages/home/Home';
import './App.css'
import EventsWrapper from './pages/Events/EventsWrapper';
import EventWrapper from './pages/Event/EventWrapper';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home/>} />
        <Route path="/events" element={<EventsWrapper />} />
        <Route path="/event/:id" element={<EventWrapper />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
