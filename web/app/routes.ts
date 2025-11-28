import {type RouteConfig, index, route} from "@react-router/dev/routes";

export default [
    index("routes/home.tsx"),
    route("login", "./routes/login.tsx"),
    route("register", "./routes/register.tsx"),
    route("confirm/:token", "./routes/confirm.tsx"),
    route("post-details/:pid", "./routes/post-details.tsx"),
] satisfies RouteConfig;
