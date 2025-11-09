import { useState } from "react"
import { API_URL } from "./App"
import { useNavigate } from "react-router-dom"
import { useCookies } from "react-cookie"

export const LoginPage: React.FC = () => {
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const redirect = useNavigate();
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [_, setCookie] = useCookies(['at']);


  const handleLogin = async () => {
    try {
      const data = await fetch(`${API_URL}/authentication/token`, {
        method: "POST",
        body: JSON.stringify({ email, password }),
      })

      const out = await data.json()
      setCookie("at", out.data)

      if (!out.error) {
        redirect("/")
      }
    } catch (error) {
      console.log('error: ', error)
    }
  }


  return (
    <div>
      <h1>Login to GopherSocial</h1>
      <input placeholder="email..." value={email} onChange={(v) => setEmail(v.target.value)} />

      <input type="password" placeholder="password..." value={password} onChange={(v) => setPassword(v.target.value)} />

      <button onClick={handleLogin}>Login</button>
    </div>
  )
}