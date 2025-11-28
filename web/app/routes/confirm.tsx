import React from 'react';
import {useNavigate, useParams} from "react-router";
import {API_URL} from "~/config";

export default function Confirm() {
    const { token = '' } = useParams()
    const redirect = useNavigate()

    const handleConfirm = async () => {
        const response = await fetch(`${API_URL}/auth/confirm/${token}`, {
            method: "PUT"
        })

        if (response.ok) {
            redirect("/")
        } else {
            // handle error
            alert("Failed to confirm token")
        }
    }

    return (
        <div>
            <h1>Confirmation</h1>
            <button onClick={handleConfirm}>Click to confirm</button>
        </div>
    );
}
