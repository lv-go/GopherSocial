import type {Route} from "./+types/post-details";
import React from 'react';
import {API_URL} from "~/config";
import type {FeedPost} from "~/components/Post";
import {Link, redirect, useNavigate} from "react-router";
import './post-details.css';

export function meta({}: Route.MetaArgs) {
    return [
        {title: "GopherSocial - Post Details"},
        {name: "description", content: "Post details page for GopherSocial."},
    ];
}

export async function clientLoader({params}: Route.ClientLoaderArgs): Promise<FeedPost | Response> {
    let authToken = sessionStorage.getItem("auth-token");
    if (!authToken) {
        return redirect("login")
    }

    const res = await fetch(`${API_URL}/posts/${params.pid}`, {
        method: "GET",
            headers: {
            "Content-Type": "application/json",
                Authorization: `Bearer ${authToken}`
        },
    });

    return await res.json();
}

export default function PostDetails({loaderData: post}: Route.ComponentProps) {
    const navigate = useNavigate();

    return (
        <div>
            <h1>{post.title}</h1>
            <p>{post.content}</p>

            <div className="comments">
                {post.comments?.map(comment => (
                    <div key={comment.id} className="comment">
                        <p>{comment.user?.username}: </p>
                        <p>{comment.content}</p>
                        <p className="comment-date">at {new Date(post.createdAt).toDateString()}</p>
                    </div>
                ))}
            </div>

            <Link to="/">Back to Home</Link>
        </div>
    );
}
