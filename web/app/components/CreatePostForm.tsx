import {type FormEvent, useState} from 'react';
import './CreatePostForm.css';
import {API_URL} from "~/config";
import {Form, useRevalidator} from "react-router";
import {getCurrentUser} from "~/services/auth-service";

export default function CreatePostForm() {
    const [title, setTitle] = useState('')
    const [content, setContent] = useState('')
    const {revalidate} = useRevalidator()

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault()
        const currentUser = await getCurrentUser()
        const at = await currentUser?.getIdToken()
        await fetch(`${API_URL}/posts`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${at}`
            },
            body: JSON.stringify({
                title,
                content
            })
        })

        setTitle('')
        setContent('')
        await revalidate()
    }

    return (
        <Form className="gopher-form">
            <label>
                <input placeholder="Title..." value={title} type="text" onChange={(e) => setTitle(e.target.value)}/>
            </label>
            <label>
                <textarea placeholder="What's in your mind..." value={content}
                          onChange={(e) => setContent(e.target.value)}/>
            </label>
            <button onClick={handleSubmit}>Share</button>
        </Form>
    )
}
