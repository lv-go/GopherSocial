import {getApp, getApps, initializeApp} from "firebase/app";
import {connectAuthEmulator, getAuth} from "firebase/auth";

const firebaseConfig = {
    apiKey: import.meta.env.VITE_FIREBASE_API_KEY || 'test-api-key',
    authDomain: import.meta.env.VITE_FIREBASE_AUTH_DOMAIN || 'test-auth-domain',
    projectId: import.meta.env.VITE_FIREBASE_PROJECT_ID || 'demo-no-project',
    appId: import.meta.env.VITE_FIREBASE_APP_ID || 'test-app-id',
}

const app = !getApps().length
    ? initializeApp(firebaseConfig)
    : getApp()


export const auth = getAuth(app)

// Connect to the Auth emulator in a development environment
if (process.env.NODE_ENV !== 'production') {
    // Ensure the URL matches the one provided by your CLI when you run 'firebase emulators:start'
    connectAuthEmulator(auth, "http://localhost:9099");
}
