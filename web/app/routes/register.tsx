import React, {useState} from 'react';
import {Form, useNavigate} from "react-router";
import {API_URL} from "~/config";

export function meta() {
    return [
        {title: "GopherSocial - Register"},
        {name: "description", content: "Register page for GopherSocial."},
    ];
}

export default function Register() {
    const [username, setUsername] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const navigate = useNavigate();

    async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
        event.preventDefault();

        try {
            const resp = await fetch(`${API_URL}/auth/register`, {
                method: "POST",
                body: JSON.stringify({ username, email, password }),
            })
            if (!resp.ok) {
                console.log('error: ', resp.statusText)
                return
            }

            navigate("/login")
        } catch (error) {
            console.log('error: ', error)
        }
    };

    return (
        <div>
            <h1>Register</h1>
            <Form onSubmit={handleSubmit}>
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
            </Form>
        </div>
    );
}
