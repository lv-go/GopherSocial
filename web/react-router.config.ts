import type { Config } from "@react-router/dev/config";

export default {
  // Config options...
  // Server-side render by default, to enable SPA mode set this to `false`
    ssr: false,
    async prerender() {
        return ["/", "/login", "/register", "/confirm/:token", "/post-details/:pid"];
    }
} satisfies Config;
