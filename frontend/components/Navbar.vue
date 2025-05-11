<script setup>
const runtimeConfig = useRuntimeConfig()
const sections = [
    {
        name: "Главная",
        link: "/"
    },
    {
        name: "Поиск",
        link: "/search"
    },
    {
        name: "Доки",
        link: runtimeConfig.public.docsUrl
    },
    {
        name: "Форум",
        link: runtimeConfig.public.forumUrl
    },
    {
        name: "Аккаунт",
        link: "/account"
    },
    {
        name: "Создать",
        link: "/new"
    },
]
</script>

<template>
    <div id="navbar">
        <NuxtLink v-for="section in sections" :key="section.name" :to="section.link">
            {{ section.name }}
        </NuxtLink>
    </div>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/colors';

#navbar {
    display: flex;
    color: white;
    background-color: colors.$darkerer-bg-color;
    padding: 10px 0;
    width: 100%;
}

#navbar a {
    color: white;
    font-style: normal;
    font-size: 24px;
    text-decoration: none;
    flex-grow: 1;
    text-align: center;
}

#navbar a:hover {
    color: cyan
}
</style>