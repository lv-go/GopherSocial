import React, {type FormEvent, useState} from 'react';
import {API_URL} from "~/config";
import {Form, Link, useNavigate} from "react-router";
import {signInWithEmailAndPassword} from "firebase/auth";
import {auth} from "~/firebase-config";

export default function Login() {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const navigate = useNavigate();

    const handleLogin = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault()
        try {
            const userCredentials = await signInWithEmailAndPassword(auth, email, password)
            console.log("userCredentials: ", userCredentials)
            // sessionStorage.setItem("auth-token", userCredentials.user.accessToken)

            navigate("/")
        } catch (error) {
            console.log('error: ', error)
        }
    }


    return (<div className="flex items-center justify-center min-h-screen bg-base-200">
        <div className="card w-full max-w-sm shadow-2xl bg-base-100">
            <Form className="card-body gap-5" onSubmit={handleLogin}>
                <h1 className="card-title">Login to GopherSocial</h1>
                <label className="validator">
                    <span className="label">Email</span>
                    <input type="email" className="input" name="email" placeholder="email..." required value={email}
                           onChange={(v) => setEmail(v.target.value)}/>
                </label>

                <label className="validator">
                    <span className="label">Password</span>
                    <input type="password" className="input" name="password" placeholder="password..." required value={password}
                           onChange={(v) => setPassword(v.target.value)}/>
                </label>

                <div className="card-actions justify-between">
                    <button type="submit" className="btn btn-primary">Login</button>
                    <Link to="/register" className="btn btn-link">Register</Link>
                </div>
            </Form>
        </div>
    </div>);
}
