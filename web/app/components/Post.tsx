import React from 'react';
import './Post.css';

interface PostComment {
    id: number
    post_id: number
    user_id: number
    content: string
    created_at: string
    user?: {
        id: number
        username: string
    }
}

export interface FeedPost {
    id: number
    userId: number
    comments_count: number
    content: string
    createdAt: string
    tags: string[]
    title?: string
    comments?: PostComment[]
}

interface PostProps {
    post: FeedPost
    onClick: () => void
}

export default function Post({post, onClick}: PostProps) {
    return (
        <div key={post.id} className="post" onClick={onClick}>
            <h2>{post.title}</h2>
            <p>{post.content}</p>

            <div>
                <p>Categories: {post.tags?.join(', ')}</p>
                <p>{new Date(post.createdAt).toDateString()}</p>
            </div>

            <div className="post-bottom">
                Click to see {post.comments_count} comments
            </div>
        </div>
    );
}
