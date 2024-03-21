import React from 'react';
import {
    BrowserRouter as Router, Route, Routes
} from "react-router-dom";
import AuthPage from './Auth/Auth';
import TrackerPage from './Task/Task';

function App() {
    return (
        <Router>
            <Routes>
                <Route path="/" element={<AuthPage />} />
                <Route path="/auth" element={<AuthPage />} />
                <Route path="/tracker" element={<TrackerPage />} />
            </Routes>
        </Router>
    );
}

export default App;
