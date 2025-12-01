import type {Route} from "./+types/home";
import {redirect, useNavigate} from "react-router";
import CreatePostForm from "~/components/CreatePostForm";
import Post, {type FeedPost} from "~/components/Post";
import {API_URL} from "~/config";
import gohper from './../../public/gohper.svg'
import {getCurrentUser} from "~/services/auth-service";
import {signOut} from "firebase/auth";
import {auth} from "~/firebase-config";

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
    let currentUser = await getCurrentUser();
    if (!currentUser) {
        return redirect("login")
    }
    let authToken = await currentUser.getIdToken();
    console.log("authToken: ", authToken)
    const res = await fetch(`${API_URL}/user/feed`, {
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
        signOut(auth)
        navigate("/login")
    }

    const handleClickPost = (id: number) => () => {
        navigate(`/post-details/${id}`);
    }

    return (
        <div className="container mx-auto px-4 text-center">
            <nav className='nav'>
                <div className='logo-container'>
                    <img src={gohper} className="logo"/>
                    <h1>GopherSocial</h1>
                </div>

                <button className="btn btn-primary" onClick={handleLogout}>Logout</button>
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
