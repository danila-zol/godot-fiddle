<script setup>
import { sections } from './globals';
</script>

<template>
    <div id="navbar">
        <a v-for="section in sections" :key="section.name" :href="section.link">
            {{ section.name }}
        </a>
    </div>
</template>

<style lang="css" scoped>
#navbar {
    display: flex;
    color: white;
    background-color: black;
}

#navbar * {
    color: white;
    font-style: normal;
    text-decoration: none;
    flex-grow: 1;
    text-align: center;
}
</style>