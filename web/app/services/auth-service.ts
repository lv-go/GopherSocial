import {onAuthStateChanged, type User} from 'firebase/auth';
import {auth} from "~/firebase-config";

export const getCurrentUser = () => {
    return new Promise<User | null>((resolve, reject) => {
        const unsubscribe = onAuthStateChanged(
            auth,
            user => {
                unsubscribe(); // Stop observing after the first emission
                resolve(user);
            },
            error => {
                unsubscribe();
                reject(error);
            }
        );
    });
};