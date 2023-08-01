import { Route, BrowserRouter as Router, Routes } from "react-router-dom"
import Page from "./MainPage"

import "./App.css"

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<Page />} />
        </Routes>
      </div>
    </Router>
  )
}

export default App
