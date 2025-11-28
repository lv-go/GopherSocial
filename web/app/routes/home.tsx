import type {Route} from "./+types/home";
import {redirect, useNavigate} from "react-router";
import CreatePostForm from "~/components/CreatePostForm";
import Post, {type FeedPost} from "~/components/Post";
import {API_URL} from "~/config";
import gohper from './../../public/gohper.svg'

export function meta({}: Route.MetaArgs) {
    return [
        {title: "GopherSocial - Home"},
        {name: "description", content: "Welcome to React Router!"},
    ];
}

interface Page<T> {
    items: T[];
    total: number;
    number: number;
    size: number;
}

export async function clientLoader(): Promise<Page<FeedPost> | Response> {
    let authToken = sessionStorage.getItem("auth-token");
    if (!authToken) {
        return redirect("login")
    }
    const res = await fetch(`${API_URL}/users/feed`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${authToken}`
        },
    });
    return res.json();
}

export default function Home({loaderData: posts}: Route.ComponentProps) {
    const navigate = useNavigate();

    const handleLogout = () => {
        console.log("logging out...")
        sessionStorage.removeItem("auth-token")
        navigate("/")
    }

    const handleClickPost = (id: number) => () => {
        navigate(`/post-details/${id}`);
    }

    return (
        <div id="root">
            <nav className='nav'>
                <div className='logo-container'>
                    <img src={gohper} className="logo"/>
                    <h1>GopherSocial</h1>
                </div>

                <button onClick={handleLogout}>Logout</button>
            </nav>

            <p>This is a social media platform for gophers.</p>

            <CreatePostForm />

            <div className='posts'>
                {posts.items.map(post => (
                    <Post key={post.id} post={post} onClick={handleClickPost(post.id)}/>
                ))}

                {posts.items.length === 0 && <p>No posts yet, start following someone or post something</p>}
            </div>

        </div>
    )
}
