import { ref } from "vue"

export const sections = ref(
    [
        {
            name: "games",
            link: "/games"
        },
        {
            name: "new game",
            link: "/new"
        },
        {
            name: "docs",
            link: import.meta.env.VITE_DOCS_URL
        },
        {
            name: "forum",
            link: import.meta.env.VITE_FORUM_URL
        },
        {
            name: "about",
            link: "/about"
        }
    ]
)

export const blogs = ref(Array(9).fill(
    { title: "Cool Blog", text: "I love bugs", thumb: "https://raw.githubusercontent.com/godotengine/godot/refs/heads/master/icon.svg" }
))