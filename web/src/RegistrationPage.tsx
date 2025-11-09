import { useState } from "react";
import { API_URL } from "./App.tsx";
import { useNavigate } from "react-router-dom";
import { useCookies } from "react-cookie";

export default function RegistrationPage() {
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const redirect = useNavigate();
  const [_, setCookie] = useCookies(['at']);

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    try {
      const data = await fetch(`${API_URL}/authentication/user`, {
        method: "POST",
        body: JSON.stringify({ username, email, password }),
      })

      const out = await data.json()
      setCookie("at", out.data)

      if (!out.error) {
        redirect("/")
      }
    } catch (error) {
      console.log('error: ', error)
    }
  };

  return (
    <div>
      <h1>Register</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="username">Username</label>
          <input
            type="text"
            id="username"
            name="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div>
          <label htmlFor="email">Email</label>
          <input
            type="text"
            id="email"
            name="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>
        <div>
          <label htmlFor="password">Password</label>
          <input
            type="password"
            id="password"
            name="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <button type="submit">Register</button>
      </form>
    </div>
  );
}
