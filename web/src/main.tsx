import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { ConfirmationPage } from './ConfirmationPage.tsx'
import { LoginPage } from "./LoginPage.tsx";
import RegistrationPage from "./RegistrationPage.tsx";

const router = createBrowserRouter([{
  path: "/",
  element: <App/>
}, {
  path: "/confirm/:token",
  element: <ConfirmationPage/>
}, {
  path: "/login",
  element: <LoginPage/>
}, {
  path: "/register",
  element: <RegistrationPage/>
}])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router}/>
  </StrictMode>,
)
