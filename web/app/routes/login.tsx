import React, {type FormEvent, useState} from 'react';
import {API_URL} from "~/config";
import {Form, Link, useNavigate} from "react-router";

export default function Login() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const navigate = useNavigate();

    const handleLogin = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        try {
            const resp = await fetch(`${API_URL}/auth/login`, {
                method: "POST",
                body: JSON.stringify({email, password}),
            })

            if (!resp.ok) {
                return
            }

            const out = await resp.json()
            sessionStorage.setItem("auth-token", out.token)

            navigate("/")
        } catch (error) {
            console.log('error: ', error)
        }
    }


    return (
        <Form className="login-form" onSubmit={handleLogin}>
            <h1>Login to GopherSocial</h1>
            <div className="form-field">
                <label>Email</label>
                <input type="email" name="email" placeholder="email..." required value={email}
                       onChange={(v) => setEmail(v.target.value)}/>
            </div>

            <div className="form-field">
                <label>Password</label>
                <input type="password" name="password" placeholder="password..." required value={password}
                       onChange={(v) => setPassword(v.target.value)}/>
            </div>

            <div className="form-footer">
                <button type="submit" className="btn btn-primary">Login</button>
                <Link to="/register">Register</Link>
            </div>
        </Form>
    );
}
