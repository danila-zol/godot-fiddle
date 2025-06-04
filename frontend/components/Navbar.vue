<script setup>
import { NuxtLink } from '#components';

let props = defineProps({
    isSsr: { default: true }
})
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
        <component :is="props.isSsr ? NuxtLink : 'a'" v-for="section in sections" :key="section.name" :to="section.link"
            :href="section.link">
            {{ section.name }}
        </component>
    </div>
</template>

<style lang="scss" scoped>
@use '~/assets/scss/colors';

#navbar {
    display: flex;
    color: white;
    background-color: colors.$navbar-color;
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
    color: colors.$light-highlight-color
}
</style>