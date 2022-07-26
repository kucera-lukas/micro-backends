import App from "./App";

import React from "react";
import ReactDOM from "react-dom/client";

import "./index.css";

ReactDOM.createRoot(document.querySelector(`#root`) as HTMLElement).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
