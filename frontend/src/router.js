import {
  BrowserRouter as Router,
  Navigate,
  Route,
  Routes,
  Outlet,
} from "react-router-dom";
import RegistrationPage from "./pages/RegistrationPage/RegistrationPage";
import LoginPage from "./pages/LoginPage/LoginPage";
import DashboardPage from "./pages/DashboardPage/DashboardPage";
import VerificationPage from "./pages/VerificationPage/VerificationPage";
import Header from "./components/layouts/Header/Header";

const AppRouter = () => {
  return (
    <Router>
      <Routes>
        <Route
          path='/'
          element={
            <>
              <Header />
              <Outlet />
            </>
          }>
          <Route
            path='/'
            element={<Navigate to={"/dashboard"} />}
          />

          <Route
            path='register'
            element={<RegistrationPage />}
          />
          <Route
            path='login'
            element={<LoginPage />}
          />
          <Route
            path='dashboard'
            element={<DashboardPage />}
          />
          <Route
            path='verification'
            element={<VerificationPage />}
          />
        </Route>
      </Routes>
    </Router>
  );
};

export default AppRouter;
