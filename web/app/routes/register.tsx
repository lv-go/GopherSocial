import React, {useState} from 'react';
import {Form, Link, useNavigate} from "react-router";
import {API_URL} from "~/config";
import {createUserWithEmailAndPassword} from "firebase/auth";
import {auth} from "~/firebase-config";

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
            const userCredentials = await createUserWithEmailAndPassword(auth, email, password)
            console.log("userCredentials: ", userCredentials)

            navigate("/login")
        } catch (error) {
            console.log('error: ', error)
        }
    };

    return (<div className="flex items-center justify-center min-h-screen bg-base-200">
        <div className="card w-full max-w-sm shadow-2xl bg-base-100">
            <Form onSubmit={handleSubmit} className="card-body gap-5">
                <h2 className="card-title">Register</h2>
                <label className="validator">
                    <span className="label">Email</span>
                    <input type="text" className="input" required id="email" name="email" value={email}
                           onChange={e => setEmail(e.target.value)}
                    />
                </label>
                <div className="form-field">
                    <label htmlFor="password" className="label">Password</label>
                    <input type="password" className="input" required id="password" name="password" value={password}
                           onChange={e => setPassword(e.target.value)}
                    />
                </div>
                <div className="flex justify-between">
                    <button type="submit" className="btn btn-primary">Register</button>
                    <Link to="/login" className="btn btn-link">Login</Link>
                </div>
            </Form>
        </div>
    </div>);
}
