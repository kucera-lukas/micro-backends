import ROUTES from "./routes";

import Home from "../pages/home";

import { Routes, Route } from "react-router-dom";

const Router = (): JSX.Element => {
  return (
    <Routes>
      <Route
        path={ROUTES.ROOT}
        element={<Home />}
      />
    </Routes>
  );
};

export default Router;
